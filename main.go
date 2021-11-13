package main

import (
	"context"
	"crypto/elliptic"
	"crypto/sha1"
	"embed"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/reports"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nuts-foundation/nuts-demo-ehr/domain/episode"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/notification"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/transfer/receiver"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/transfer/sender"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"

	"github.com/nuts-foundation/nuts-demo-ehr/api"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/customers"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/dossier"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/patients"
	httpAuth "github.com/nuts-foundation/nuts-demo-ehr/http/auth"
	"github.com/nuts-foundation/nuts-demo-ehr/http/proxy"
	nutsClient "github.com/nuts-foundation/nuts-demo-ehr/nuts/client"
	nutsAuthClient "github.com/nuts-foundation/nuts-demo-ehr/nuts/client/auth"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts/registry"
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
	logrusLevel, err := logrus.ParseLevel(config.Verbosity)
	if err != nil {
		panic(err)
	}
	logrus.SetLevel(logrusLevel)

	if config.FHIR.Server.Type == "" {
		logrus.Fatal("Invalid FHIR server type, valid options are: 'hapi-multi-tenant', 'hapi' or 'other'")
	}

	// init node API nutsClient
	nodeClient := nutsClient.HTTPClient{NutsNodeAddress: config.NutsNodeAddress}
	vcRegistry := registry.NewVerifiableCredentialRegistry(&nodeClient)
	customerRepository := customers.NewJsonFileRepository(config.CustomersFile)

	server := createServer()

	registerEHR(server, config, customerRepository, vcRegistry)

	if config.FHIR.Proxy.Enable {
		registerFHIRProxy(server, config, customerRepository, vcRegistry)
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

func registerFHIRProxy(server *echo.Echo, config Config, customerRepository customers.Repository, vcRegistry registry.VerifiableCredentialRegistry) {
	authService, err := httpAuth.NewService(config.NutsNodeAddress)
	if err != nil {
		log.Fatal(err)
	}

	fhirURL, err := url.Parse(config.FHIR.Server.Address)
	if err != nil {
		log.Fatal(err)
	}
	proxyServer := proxy.NewServer(authService, customerRepository, vcRegistry, *fhirURL, config.FHIR.Proxy.Path, config.FHIR.Server.SupportsMultiTenancy())

	// set security filter
	server.Use(proxyServer.AuthMiddleware())

	server.Any(config.FHIR.Proxy.Path+"/*", func(c echo.Context) error {
		// For all non-Task updates, logic is performed by middleware

		// For the Task update it does not proxy the call but forwards it to here.
		// The request is altered to /web/internal/{customerID}/Task/{ID}
		server.ServeHTTP(c.Response(), c.Request())
		return nil
	}, proxyServer.Handler)
}

func registerEHR(server *echo.Echo, config Config, customerRepository customers.Repository, vcRegistry registry.VerifiableCredentialRegistry) {
	// init node API nutsClient
	nodeClient := nutsClient.HTTPClient{NutsNodeAddress: config.NutsNodeAddress}

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

	authService, err := httpAuth.NewService(config.NutsNodeAddress)
	if err != nil {
		log.Fatal(err)
	}

	fhirClientFactory := fhir.NewFactory(fhir.WithURL(config.FHIR.Server.Address), fhir.WithMultiTenancyEnabled(config.FHIR.Server.SupportsMultiTenancy()))
	patientRepository := patients.NewFHIRPatientRepository(patients.Factory{}, fhirClientFactory)
	reportRepository := reports.NewFHIRRepository(fhirClientFactory)
	orgRegistry := registry.NewOrganizationRegistry(&nodeClient)
	dossierRepository := dossier.NewSQLiteDossierRepository(dossier.Factory{}, sqlDB)
	transferSenderRepo := sender.NewTransferRepository(sqlDB)
	transferReceiverRepo := receiver.NewTransferRepository(sqlDB)
	transferSenderService := sender.NewTransferService(authService, fhirClientFactory, transferSenderRepo, customerRepository, dossierRepository, patientRepository, orgRegistry, vcRegistry)
	transferReceiverService := receiver.NewTransferService(authService, fhirClientFactory, transferReceiverRepo, customerRepository, orgRegistry, vcRegistry)
	tenantInitializer := func(tenant int) error {
		if !config.FHIR.Server.SupportsMultiTenancy() {
			return nil
		}

		return fhir.InitializeTenant(config.FHIR.Server.Address, strconv.Itoa(tenant))
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
	auth := api.NewAuth(config.sessionKey, customerRepository, passwd)

	// Initialize wrapper
	apiWrapper := api.Wrapper{
		APIAuth:                 auth,
		NutsAuth:                nodeClient,
		CustomerRepository:      customerRepository,
		PatientRepository:       patientRepository,
		ReportRepository:        reportRepository,
		DossierRepository:       dossier.NewSQLiteDossierRepository(dossier.Factory{}, sqlDB),
		TransferSenderRepo:      transferSenderRepo,
		OrganizationRegistry:    orgRegistry,
		TransferSenderService:   transferSenderService,
		TransferReceiverService: transferReceiverService,
		TransferReceiverRepo:    transferReceiverRepo,
		EpisodeService:          episode.NewService(fhirClientFactory, authService, orgRegistry, vcRegistry),
		TenantInitializer:       tenantInitializer,
		NotificationHandler:     notification.NewHandler(authService, fhirClientFactory, transferReceiverService, orgRegistry),
	}

	// JWT checking for correct claims
	server.Use(auth.JWTHandler)
	server.Use(sql.Transactional(sqlDB))

	// for requests that require Nuts AccessToken
	server.Use(authMiddleware(authService))

	api.RegisterHandlersWithBaseURL(server, apiWrapper, "/web")

	// Setup asset serving:
	// Check if we use live mode from the file system or using embedded files
	useFS := len(os.Args) > 1 && os.Args[1] == "live"
	assetHandler := http.FileServer(getFileSystem(useFS))

	server.GET("/*", echo.WrapHandler(assetHandler))
}

func registerPatients(repository patients.Repository, db *sqlx.DB, customerID int) {
	pdate := func(value time.Time) *openapi_types.Date {
		val := openapi_types.Date{value}
		return &val
	}
	pstring := func(value string) *string {
		return &value
	}
	props := []types.PatientProperties{
		{
			Ssn:       pstring("1234567890"),
			Dob:       pdate(time.Date(1980, 10, 10, 0, 0, 0, 0, time.UTC)),
			FirstName: "Henk",
			Surname:   "de Vries",
			Gender:    types.PatientPropertiesGenderMale,
			Zipcode:   "6825AX",
		},
		{
			Ssn:       pstring("1234567891"),
			Dob:       pdate(time.Date(1939, 1, 5, 0, 0, 0, 0, time.UTC)),
			FirstName: "Grepelsteeltje",
			Surname:   "Grouw",
			Gender:    types.PatientPropertiesGenderFemale,
			Zipcode:   "9999AA",
		},
		{
			Ssn:       pstring("1234567892"),
			Dob:       pdate(time.Date(1972, 1, 10, 0, 0, 0, 0, time.UTC)),
			FirstName: "Dibbes",
			Surname:   "Bouwman",
			Gender:    types.PatientPropertiesGenderMale,
			Zipcode:   "1234ZZ",
		},
		{
			Ssn:       pstring("1234567893"),
			Dob:       pdate(time.Date(2001, 2, 27, 0, 0, 0, 0, time.UTC)),
			FirstName: "Anne",
			Surname:   "von Oben",
			Gender:    types.PatientPropertiesGenderOther,
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

func authMiddleware(authService httpAuth.Service) echo.MiddlewareFunc {
	config := httpAuth.Config{
		Skipper: func(e echo.Context) bool {
			return !strings.HasPrefix(e.Request().RequestURI, "/web/external/")
		},
		AccessF: func(ctx echo.Context, request *http.Request, token *nutsAuthClient.TokenIntrospectionResponse) error {
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
	return httpAuth.SecurityFilter{Auth: authService}.AuthWithConfig(config)
}
