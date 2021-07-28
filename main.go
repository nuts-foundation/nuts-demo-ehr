package main

import (
	"context"
	"crypto/elliptic"
	"crypto/sha1"
	"embed"
	"encoding/hex"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/nuts-foundation/nuts-demo-ehr/api"
	"github.com/nuts-foundation/nuts-demo-ehr/client"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/customers"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/dossier"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/patients"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/registry"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/transfer"
	"github.com/nuts-foundation/nuts-demo-ehr/sql"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log2 "github.com/labstack/gommon/log"
)

const assetPath = "web/dist"

//go:embed web/dist/*
var embeddedFiles embed.FS

const apiTimeout = 10 * time.Second

func getFileSystem(useFS bool) http.FileSystem {
	if useFS {
		logrus.Info("using live mode")
		return http.FS(os.DirFS(assetPath))
	}

	logrus.Info("using embed mode")
	fsys, err := fs.Sub(embeddedFiles, assetPath)
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}

func main() {
	// config stuff
	config := loadConfig()
	config.Print(log.Writer())

	// init node API client
	nodeClient := client.HTTPClient{NutsNodeAddress: config.NutsNodeAddress}

	var passwd string
	if config.Credentials.Empty() {
		passwd = generateAuthenticationPassword(config)
		logrus.Infof("Authentication credentials not configured, so they were generated (password=%s)", passwd)
	} else {
		passwd = config.Credentials.Password
	}

	// Initialize services
	customerRepository := customers.NewJsonFileRepository(config.CustomersFile)
	sqlDB := sqlx.MustConnect("sqlite3", config.DBConnectionString)
	sqlDB.SetMaxOpenConns(1)
	//patientRepository := patients.NewSQLitePatientRepository(patients.Factory{}, sqlDB)
	patientRepository := patients.NewFHIRPatientRepository(patients.Factory{}, config.FHIRServerAddress)
	if config.LoadTestPatients {
		allCustomers, err := customerRepository.All()
		if err != nil {
			log.Fatal(err)
		}
		for _, customer := range allCustomers {
			registerPatients(patientRepository, sqlDB, customer.Id)
		}
	}
	auth := api.NewAuth(config.sessionKey, nodeClient, customerRepository, passwd)

	// Initialize wrapper
	apiWrapper := api.Wrapper{
		Auth:                 auth,
		Client:               nodeClient,
		CustomerRepository:   customerRepository,
		PatientRepository:    patientRepository,
		DossierRepository:    dossier.NewSQLiteDossierRepository(dossier.Factory{}, sqlDB),
		TransferRepository:   transfer.NewSQLiteTransferRepository(transfer.Factory{}, sqlDB),
		OrganizationRegistry: registry.NewOrganizationRegistry(&nodeClient),
		FHIRGateway:          &fhir.StubGateway{},
	}
	e := echo.New()
	e.HideBanner = true
	// Register Echo logger middleware but do not log calls to the status endpoint,
	// since that gets called by the Docker healthcheck very, very often which leads to lots of clutter in the log.
	e.GET("/status", func(c echo.Context) error {
		c.Response().WriteHeader(http.StatusNoContent)
		return nil
	})
	loggerConfig := middleware.DefaultLoggerConfig
	loggerConfig.Skipper = func(ctx echo.Context) bool {
		return ctx.Request().RequestURI == "/status"
	}
	e.Use(middleware.LoggerWithConfig(loggerConfig))
	// JWT checking for correct claims
	e.Use(auth.JWTHandler)
	e.Use(sql.Transactional(sqlDB))
	e.Logger.SetLevel(log2.DEBUG)
	e.HTTPErrorHandler = func(err error, ctx echo.Context) {
		if !ctx.Response().Committed {
			ctx.Response().Write([]byte(err.Error()))
			ctx.Echo().Logger.Error(err)
		}
	}
	e.HTTPErrorHandler = httpErrorHandler

	api.RegisterHandlersWithBaseURL(e, apiWrapper, "/web")

	// Setup asset serving:
	// Check if we use live mode from the file system or using embedded files
	useFS := len(os.Args) > 1 && os.Args[1] == "live"
	assetHandler := http.FileServer(getFileSystem(useFS))
	e.GET("/*", echo.WrapHandler(assetHandler))

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.HTTPPort)))
}

func registerPatients(repository patients.Repository, db *sqlx.DB, customerID string) {
	pdate := func(value time.Time) *openapi_types.Date {
		val := openapi_types.Date{value}
		return &val
	}
	pstring := func(value string) *string {
		return &value
	}
	props := []domain.PatientProperties{
		{
			Ssn:       pstring("1234567890"),
			Dob:       pdate(time.Date(1980, 10, 10, 0, 0, 0, 0, time.UTC)),
			FirstName: "Henk",
			Surname:   "de Vries",
			Gender:    domain.PatientPropertiesGenderMale,
			Zipcode:   "6825AX",
		},
		{
			Ssn:       pstring("1234567891"),
			Dob:       pdate(time.Date(1939, 1, 5, 0, 0, 0, 0, time.UTC)),
			FirstName: "Grepelsteeltje",
			Surname:   "Grouw",
			Gender:    domain.PatientPropertiesGenderFemale,
			Zipcode:   "9999AA",
		},
		{
			Ssn:       pstring("1234567892"),
			Dob:       pdate(time.Date(1972, 1, 10, 0, 0, 0, 0, time.UTC)),
			FirstName: "Dibbes",
			Surname:   "Bouwman",
			Gender:    domain.PatientPropertiesGenderMale,
			Zipcode:   "1234ZZ",
		},
		{
			Ssn:       pstring("1234567893"),
			Dob:       pdate(time.Date(2001, 2, 27, 0, 0, 0, 0, time.UTC)),
			FirstName: "Anne",
			Surname:   "von Oben",
			Gender:    domain.PatientPropertiesGenderOther,
			Zipcode:   "7777AX",
		},
	}
	if err := sql.ExecuteTransactional(db, func(ctx context.Context) error {
		for _, prop := range props {
			if _, err := repository.NewPatient(ctx, customerID, prop); err != nil {
				return fmt.Errorf("unable to register test patient: %w", err)
			}
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}
}

func generateAuthenticationPassword(config Config) string {
	pkHashBytes := sha1.Sum(elliptic.Marshal(config.sessionKey.Curve, config.sessionKey.X, config.sessionKey.Y))
	return hex.EncodeToString(pkHashBytes[:])
}

// httpErrorHandler includes the err.Err() string in a { "error": "msg" } json hash
func httpErrorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		msg  interface{}
	)
	type Map map[string]interface{}

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
		if he.Internal != nil {
			err = fmt.Errorf("%v, %v", err, he.Internal)
		}
	} else {
		msg = err.Error()
	}

	if _, ok := msg.(string); ok {
		msg = Map{"error": msg}
	}

	// Send response
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, msg)
		}
		if err != nil {
			c.Logger().Error(err)
		}
	}
}
