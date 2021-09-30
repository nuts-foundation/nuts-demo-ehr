package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/spf13/pflag"
)

const defaultPrefix = "NUTS_"
const defaultDelimiter = "."
const configFileFlag = "configfile"
const defaultConfigFile = "server.config.yaml"
const defaultHTTPPort = 1304
const defaultNutsNodeAddress = "http://localhost:1323"
const defaultCustomerFile = "customers.json"
const defaultLogLevel = "info"

// defaultHAPIFHIRServer configures usage of the HAPI FHIR Server (https://hapifhir.io/)
var defaultHAPIFHIRServer = FHIRServer{
	Address: "http://localhost:8080/fhir",
}

func defaultConfig() Config {
	return Config{
		HTTPPort:        defaultHTTPPort,
		NutsNodeAddress: defaultNutsNodeAddress,
		Verbosity:       defaultLogLevel,
		FHIR: FHIR{
			Server: defaultHAPIFHIRServer,
			Proxy: FHIRProxy{
				Enable: true,
				Path:   "/fhir",
			},
		},
		CustomersFile:      defaultCustomerFile,
		Credentials:        Credentials{Password: "demo"},
		DBConnectionString: "demo-ehr.db?cache=shared",
		LoadTestPatients:   false,
	}
}

type Config struct {
	Credentials     Credentials `koanf:"credentials"`
	Verbosity       string      `koanf:"verbosity"`
	HTTPPort        int         `koanf:"port"`
	NutsNodeAddress string      `koanf:"nutsnodeaddr"`
	FHIR            FHIR        `koanf:"fhir"`
	CustomersFile   string      `koanf:"customersfile"`
	Branding        Branding    `koanf:"branding"`
	// Database connection string, accepts all options for the sqlite3 driver
	// https://github.com/mattn/go-sqlite3#connection-string
	DBConnectionString string `koanf:"dbConnectionString"`
	// Load a set of test patients on startup. Should be disabled for permanent data stores.
	LoadTestPatients bool `koanf:"loadTestPatients"`
	// If set, this key wil be used to sign JWTs. If not set, a new one is generated on each start up.
	// Developer tip: set the sessionPemKey so the session keeps valid after a server reboot.
	SessionPemKey string `koanf:"sessionPemKey"`
	sessionKey    *ecdsa.PrivateKey
}

type FHIR struct {
	Server FHIRServer `koanf:"server"`
	Proxy  FHIRProxy  `koanf:"proxy"`
}

type FHIRServer struct {
	Type    string `koanf:"type"`
	Address string `koanf:"address"`
}

func (server FHIRServer) SupportsMultiTenancy() bool {
	return server.Type == "hapi-multi-tenant"
}

type FHIRProxy struct {
	Enable bool   `koanf:"enable"`
	Path   string `koanf:"path"`
}

type Credentials struct {
	Password string `koanf:"password" json:"-"` // json omit tag to avoid having it printed in server log
}

type Branding struct {
	// Logo defines a path that points to a file on disk to be used as logo, to be displayed in the application.
	Logo string `koanf:"logo"`
}

func (c Credentials) Empty() bool {
	return len(c.Password) == 0
}

func generateSessionKey() (*ecdsa.PrivateKey, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		logrus.Printf("failed to generate private key: %s", err)
		return nil, err
	}
	keyBytes, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, err
	}
	block := pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes}
	pem.Encode(log.Writer(), &block)

	return key, nil
}

func (c Config) Print(writer io.Writer) error {
	if _, err := fmt.Fprintln(writer, "========== CONFIG: =========="); err != nil {
		return err
	}
	var pr Config = c
	if pr.SessionPemKey != "" {
		// Don't print the private key.
		pr.SessionPemKey = "Redacted"
	}
	data, _ := json.MarshalIndent(pr, "", "  ")
	if _, err := fmt.Println(writer, string(data)); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(writer, "========= END CONFIG ========="); err != nil {
		return err
	}
	return nil
}

func loadConfig() Config {
	flagset := loadFlagSet(os.Args[1:])

	var k = koanf.New(".")

	// Prepare koanf for parsing the config file
	configFilePath := resolveConfigFile(flagset)
	// Check if the file exists
	if _, err := os.Stat(configFilePath); err == nil {
		logrus.Infof("Loading config from file: %s", configFilePath)
		if err := k.Load(file.Provider(configFilePath), yaml.Parser()); err != nil {
			logrus.Fatalf("error while loading config from file: %v", err)
		}
	} else {
		logrus.Infof("Using default config because no file was found at: %s", configFilePath)
	}
	// load env flags, can't return error
	_ = k.Load(envProvider(), nil)

	config := defaultConfig()

	// Unmarshal values of the config file into the config struct, potentially replacing default values
	if err := k.Unmarshal("", &config); err != nil {
		log.Fatalf("error while unmarshalling config: %v", err)
	}

	if len(config.SessionPemKey) > 0 {
		log.Print("sessionPemKey set, trying to parse it...")
		block, _ := pem.Decode([]byte(config.SessionPemKey))
		if block == nil || block.Type != "EC PRIVATE KEY" {
			log.Fatalf("unable to parse sessionPemKey as PEM")
		}
		key, err := x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			log.Fatalf("unable to parse sessionPemKey as EC Private key: %v", err)
		}
		config.sessionKey = key
	} else {
		sessionKey, err := generateSessionKey()
		if err != nil {
			log.Fatalf("unable to generate session key: %v", err)
		}
		config.sessionKey = sessionKey
	}

	return config
}

func loadFlagSet(args []string) *pflag.FlagSet {
	f := pflag.NewFlagSet("config", pflag.ContinueOnError)
	f.String(configFileFlag, defaultConfigFile, "Nuts config file")
	f.Usage = func() {
		fmt.Println(f.FlagUsages())
		os.Exit(0)
	}

	err := f.Parse(args)
	if err != nil {
		panic(err)
	}

	return f
}

// resolveConfigFile resolves the path of the config file using the following sources:
// 1. commandline params (using the given flags)
// 2. environment vars,
// 3. default location.
func resolveConfigFile(flagset *pflag.FlagSet) string {
	k := koanf.New(defaultDelimiter)

	// load env flags, can't return error
	_ = k.Load(envProvider(), nil)

	// load cmd flags, without a parser, no error can be returned
	_ = k.Load(posflag.Provider(flagset, defaultDelimiter, k), nil)

	configFile := k.String(configFileFlag)

	return configFile
}

func envProvider() *env.Env {
	return env.Provider(defaultPrefix, defaultDelimiter, func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, defaultPrefix)), "_", defaultDelimiter, -1)
	})
}
