package main

import (
	"context"
	"crypto/elliptic"
	"crypto/sha1"
	"crypto/tls"
	"embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/nuts-foundation/nuts-demo-ehr/domain"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/acl"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/sharedcareplan"
	openapiTypes "github.com/oapi-codegen/runtime/types"

	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nuts-foundation/nuts-demo-ehr/api"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/customers"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/dossier"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/episode"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/fhir"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/notification"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/patients"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/reports"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/transfer"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/transfer/receiver"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/transfer/sender"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
	"github.com/nuts-foundation/nuts-demo-ehr/internal/keyring"
	nutsClient "github.com/nuts-foundation/nuts-demo-ehr/nuts/client"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts/registry"
	"github.com/nuts-foundation/nuts-demo-ehr/sql"

	"github.com/sirupsen/logrus"
)

const assetPath = "web/dist"

//go:embed web/dist/*
var embeddedFiles embed.FS

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

	// Read the authentication key
	var authorizer *nutsClient.Authorizer
	if keyPath := config.NutsNodeKeyPath; keyPath != "" {
		key, err := keyring.Open(keyPath)
		if err != nil {
			logrus.Fatalf("Failed to open nuts-node key: %v", err)
		}

		authorizer = &nutsClient.Authorizer{Key: key, Audience: config.NutsNodeAPIAudience}
	}

	// init node API nutsClient
	nodeClient := nutsClient.HTTPClient{NutsNodeAddress: config.NutsNodeAddress, Authorizer: authorizer}
	vcRegistry := registry.NewVerifiableCredentialRegistry(&nodeClient)
	customerRepository := customers.NewJsonFileRepository(config.CustomersFile)

	server := createServer()

	registerEHR(server, config, customerRepository, vcRegistry, &nodeClient)

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
	server.Logger.SetLevel(1)
	server.HTTPErrorHandler = func(err error, ctx echo.Context) {
		if !ctx.Response().Committed {
			_, _ = ctx.Response().Write([]byte(err.Error()))
			ctx.Echo().Logger.Error(err)
		}
	}
	server.Binder = &fhirBinder{}
	server.HTTPErrorHandler = httpErrorHandler
	return server
}

func registerEHR(server *echo.Echo, config Config, customerRepository customers.Repository, vcRegistry registry.VerifiableCredentialRegistry, nodeClient *nutsClient.HTTPClient) {
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

	fhirNotifier := transfer.FireAndForgetNotifier{}
	var tlsClientConfig *tls.Config
	var err error
	if config.TLS.Client.IsConfigured() {
		log.Println("Configuring TLS client certificate for calls to remote Nuts Nodes and FHIR servers.")
		if tlsClientConfig, err = config.TLS.Client.Load(); err != nil {
			log.Fatal(err)
		}
		fhirNotifier.TLSConfig = tlsClientConfig
	}

	fhirClientFactory := fhir.NewFactory(
		fhir.WithURL(config.FHIR.Server.Address),
		fhir.WithMultiTenancyEnabled(config.FHIR.Server.SupportsMultiTenancy()),
		fhir.WithTLS(tlsClientConfig),
	)
	patientRepository := patients.NewFHIRPatientRepository(patients.Factory{}, fhirClientFactory)
	reportRepository := reports.NewFHIRRepository(fhirClientFactory)
	orgRegistry := registry.NewOrganizationRegistry(nodeClient)
	dossierRepository := dossier.NewSQLiteDossierRepository(dossier.Factory{}, sqlDB)
	transferSenderRepo := sender.NewTransferRepository(sqlDB)
	transferReceiverRepo := receiver.NewTransferRepository(sqlDB)
	transferSenderService := sender.NewTransferService(nodeClient, fhirClientFactory, transferSenderRepo, customerRepository, dossierRepository, patientRepository, orgRegistry, vcRegistry, fhirNotifier)
	transferReceiverService := receiver.NewTransferService(nodeClient, fhirClientFactory, transferReceiverRepo, customerRepository, orgRegistry, vcRegistry, fhirNotifier)
	tenantInitializer := func(tenant int) error {
		if !config.FHIR.Server.SupportsMultiTenancy() {
			return nil
		}

		return fhir.InitializeTenant(config.FHIR.Server.Address, strconv.Itoa(tenant))
	}

	// Shared Care Plan
	var scpService *sharedcareplan.Service
	if config.SharedCarePlanning.Enabled() {
		scpRepository, err := sharedcareplan.NewRepository(sqlDB)
		if err != nil {
			log.Fatal(err)
		}
		scpFHIRClient := fhir.NewFactory(fhir.WithURL(config.SharedCarePlanning.CarePlanService.FHIRBaseURL))()
		scpService = &sharedcareplan.Service{DossierRepository: dossierRepository, PatientRepository: patientRepository, Repository: scpRepository, FHIRClient: scpFHIRClient}
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

	aclRepository, err := acl.NewRepository(sqlDB)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize wrapper
	apiWrapper := api.Wrapper{
		APIAuth:                 auth,
		ACL:                     aclRepository,
		NutsClient:              nodeClient,
		CustomerRepository:      customerRepository,
		PatientRepository:       patientRepository,
		ReportRepository:        reportRepository,
		DossierRepository:       dossier.NewSQLiteDossierRepository(dossier.Factory{}, sqlDB),
		TransferSenderRepo:      transferSenderRepo,
		OrganizationRegistry:    orgRegistry,
		TransferSenderService:   transferSenderService,
		TransferReceiverService: transferReceiverService,
		TransferReceiverRepo:    transferReceiverRepo,
		ZorginzageService:       domain.ZorginzageService{NutsClient: nodeClient},
		SharedCarePlanService:   scpService,
		FHIRService:             fhir.Service{ClientFactory: fhirClientFactory},
		EpisodeService:          episode.NewService(fhirClientFactory, nodeClient, orgRegistry, vcRegistry, aclRepository),
		TenantInitializer:       tenantInitializer,
		NotificationHandler:     notification.NewHandler(nodeClient, fhirClientFactory, transferReceiverService, orgRegistry, vcRegistry),
	}

	// JWT checking for correct claims
	server.Use(auth.JWTHandler)
	server.Use(sql.Transactional(sqlDB))

	api.RegisterHandlersWithBaseURL(server, apiWrapper, "/web")

	// Setup asset serving:
	// Check if we use live mode from the file system or using embedded files
	useFS := len(os.Args) > 1 && os.Args[1] == "live"
	assetHandler := http.FileServer(getFileSystem(useFS))

	server.GET("/*", echo.WrapHandler(assetHandler))
}

func registerPatients(repository patients.Repository, db *sqlx.DB, customerID int) {
	pdate := func(value time.Time) *openapiTypes.Date {
		val := openapiTypes.Date{Time: value}
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
			Gender:    types.Male,
			Zipcode:   "6825AX",
		},
		{
			Ssn:       pstring("1234567891"),
			Dob:       pdate(time.Date(1939, 1, 5, 0, 0, 0, 0, time.UTC)),
			FirstName: "Grepelsteeltje",
			Surname:   "Grouw",
			Gender:    types.Female,
			Zipcode:   "9999AA",
		},
		{
			Ssn:       pstring("1234567892"),
			Dob:       pdate(time.Date(1972, 1, 10, 0, 0, 0, 0, time.UTC)),
			FirstName: "Dibbes",
			Surname:   "Bouwman",
			Gender:    types.Male,
			Zipcode:   "1234ZZ",
		},
		{
			Ssn:       pstring("1234567893"),
			Dob:       pdate(time.Date(2001, 2, 27, 0, 0, 0, 0, time.UTC)),
			FirstName: "Anne",
			Surname:   "von Oben",
			Gender:    types.Other,
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
			msg = fmt.Sprintf("%v, %v", err, he.Internal)
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

type fhirBinder struct{}

func (cb *fhirBinder) Bind(i interface{}, c echo.Context) (err error) {
	// You may use default binder
	db := new(echo.DefaultBinder)
	if err = db.Bind(i, c); err != echo.ErrUnsupportedMediaType {
		return
	}

	if strings.Contains(c.Request().Header.Get("Content-Type"), "application/fhir+json") {
		var bytes []byte
		if bytes, err = io.ReadAll(c.Request().Body); err != nil {
			return
		}

		if err = json.Unmarshal(bytes, i); err != nil {
			return
		}
	}

	return
}
