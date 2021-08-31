package main

import (
	"context"
	"crypto/elliptic"
	"crypto/sha1"
	"embed"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/nuts-foundation/nuts-demo-ehr/api"
	"github.com/nuts-foundation/nuts-demo-ehr/client"
	"github.com/nuts-foundation/nuts-demo-ehr/client/auth"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	auth_service "github.com/nuts-foundation/nuts-demo-ehr/domain/auth"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/customers"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/dossier"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/inbox"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/patients"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/registry"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/transfer"
	http2 "github.com/nuts-foundation/nuts-demo-ehr/http"
	"github.com/nuts-foundation/nuts-demo-ehr/proxy"
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

	customerRepository := customers.NewJsonFileRepository(config.CustomersFile)

	server := createServer()
	registerEHR(server, config, customerRepository)
	if config.FHIR.Proxy.Enable {
		registerFHIRProxy(server, config, customerRepository)
	}

	// Start server
	server.Logger.Fatal(server.Start(fmt.Sprintf(":%d", config.HTTPPort)))
}

func createServer() *echo.Echo {
	server := echo.New()
	server.HideBanner = true
	// Register Echo logger middleware but do not log calls to the status endpoint,
	// since that gets called by the Docker healthcheck very, very often which leads to lots of clutter in the log.
	server.GET("/status", func(c echo.Context) error {
		c.Response().WriteHeader(http.StatusNoContent)
		return nil
	})
	loggerConfig := middleware.DefaultLoggerConfig
	loggerConfig.Skipper = func(ctx echo.Context) bool {
		return ctx.Request().RequestURI == "/status"
	}
	server.Use(middleware.LoggerWithConfig(loggerConfig))
	server.Logger.SetLevel(log2.DEBUG)
	server.HTTPErrorHandler = func(err error, ctx echo.Context) {
		if !ctx.Response().Committed {
			ctx.Response().Write([]byte(err.Error()))
			ctx.Echo().Logger.Error(err)
		}
	}
	server.HTTPErrorHandler = httpErrorHandler
	return server
}

func registerFHIRProxy(server *echo.Echo, config Config, customerRepository customers.Repository) {
	authService, err := auth_service.NewService(config.NutsNodeAddress)
	if err != nil {
		log.Fatal(err)
	}

	fhirURL, err := url.Parse(config.FHIR.Server.Address)
	if err != nil {
		log.Fatal(err)
	}
	proxyServer := proxy.NewServer(authService, customerRepository, *fhirURL, config.FHIR.Proxy.Path)

	// set security filter
	server.Use(proxyServer.AuthMiddleware())

	server.Any(config.FHIR.Proxy.Path+"/*", func(c echo.Context) error {
		// Logic performed by middleware
		return nil
	}, proxyServer.Handler)
}

func registerEHR(server *echo.Echo, config Config, customerRepository customers.Repository) {
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
	sqlDB := sqlx.MustConnect("sqlite3", config.DBConnectionString)
	sqlDB.SetMaxOpenConns(1)

	authService, err := auth_service.NewService(config.NutsNodeAddress)
	if err != nil {
		log.Fatal(err)
	}

	fhirClientFactory := fhir.NewFactory(fhir.WithURL(config.FHIR.Server.Address))
	patientRepository := patients.NewFHIRPatientRepository(patients.Factory{}, fhirClientFactory)
	transferRepository := transfer.NewSQLiteTransferRepository(sqlDB)
	orgRegistry := registry.NewOrganizationRegistry(&nodeClient)
	vcRegistry := registry.NewVerifiableCredentialRegistry(&nodeClient)
	transferService := transfer.NewTransferService(authService, fhirClientFactory, transferRepository, customerRepository, orgRegistry, vcRegistry)
	tenantInitializer := func(tenant string) error {
		return fhir.InitializeTenant(config.FHIR.Server.Type, config.FHIR.Server.Address, tenant)
	}

	if config.LoadTestPatients {
		allCustomers, err := customerRepository.All()
		if err != nil {
			log.Fatal(err)
		}
		for _, customer := range allCustomers {
			if err := tenantInitializer(customer.Id); err != nil {
				log.Fatal(err)
			}
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
		TransferRepository:   transferRepository,
		OrganizationRegistry: orgRegistry,
		TransferService:      transferService,
		Inbox:                inbox.NewService(customerRepository, inbox.NewRepository(sqlDB), orgRegistry, authService),
		TenantInitializer:    tenantInitializer,
	}

	// JWT checking for correct claims
	server.Use(auth.JWTHandler)
	server.Use(sql.Transactional(sqlDB))

	// for requests that require Nuts AccesToken
	server.Use(authMiddleware(authService))

	api.RegisterHandlersWithBaseURL(server, apiWrapper, "/web")

	// Setup asset serving:
	// Check if we use live mode from the file system or using embedded files
	useFS := len(os.Args) > 1 && os.Args[1] == "live"
	assetHandler := http.FileServer(getFileSystem(useFS))

	server.GET("/*", echo.WrapHandler(assetHandler))
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

func authMiddleware(authService auth_service.Service) echo.MiddlewareFunc {
	config := http2.Config{
		Skipper: func(e echo.Context) bool {
			return e.Request().RequestURI != "/web/external/transfer/notify"
		},
		AccessF: func(request *http.Request, token *auth.TokenIntrospectionResponse) error {
			service := token.Service
			if service == nil {
				return errors.New("access-token doesn't contain 'service' claim")
			}
			if *service != "eOverdracht-receiver" {
				return fmt.Errorf("access-token contains incorrect 'service' claim: %s, must be eOverdracht-receiver", *service)
			}

			return nil
		},
	}
	return http2.SecurityFilter{Auth: authService}.AuthWithConfig(config)
}
