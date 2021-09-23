package api

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	jwt2 "github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/lestrrat-go/jwx/jwt/openid"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/customers"
)

const MaxSessionAge = time.Hour
const CustomerID = "cid"
const SessionID = "sid"
const Elevated = "elv"

type Auth struct {
	// sessions maps session identifiers to base64 encoded VPs
	sessions map[string]Session
	// mux is used to secure access to internal state of this struct to prevent racy behaviour
	mux        sync.Mutex
	customers  customers.Repository
	password   string
	sessionKey *ecdsa.PrivateKey
}

type Session struct {
	credential interface{}
	customerID int
	startTime  time.Time
}

type JWTCustomClaims struct {
	CustomerID int    `json:"cis"`
	SessionID  string `json:"sid"`
	jwt2.StandardClaims
}

func NewAuth(key *ecdsa.PrivateKey, customers customers.Repository, passwd string) *Auth {
	return &Auth{
		sessionKey: key,
		customers:  customers,
		sessions:   map[string]Session{},
		password:   passwd,
	}
}

// CreateCustomerJWT creates a JWT that only stores the customer ID. This is required for the IRMA flow.
func (auth *Auth) CreateCustomerJWT(customerId int) ([]byte, error) {
	t := openid.New()
	t.Set(jwt.IssuedAtKey, time.Now())
	t.Set(jwt.ExpirationKey, time.Now().Add(MaxSessionAge))
	t.Set(CustomerID, customerId)

	return jwt.Sign(t, jwa.ES256, auth.sessionKey)
}

func (auth Auth) GetCustomerIDFromHeader(ctx echo.Context) (int, error) {
	token, err := auth.extractJWTFromHeader(ctx)
	if err != nil {
		ctx.Echo().Logger.Error(err)
		return 0, echo.NewHTTPError(http.StatusUnauthorized, err)
	}
	rawID, ok := token.Get(CustomerID)
	if !ok {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "missing customerID in token")
	}
	return int(rawID.(float64)), nil
}

// CreateSessionJWT creates a JWT with customer ID and session ID
func (auth *Auth) CreateSessionJWT(subject string, customerId int, session string, elevated bool) ([]byte, error) {
	t := openid.New()
	t.Set(jwt.SubjectKey, subject)
	t.Set(jwt.IssuedAtKey, time.Now())
	t.Set(jwt.ExpirationKey, time.Now().Add(MaxSessionAge))
	t.Set(CustomerID, customerId)
	t.Set(SessionID, session)
	t.Set(Elevated, elevated)

	return jwt.Sign(t, jwa.ES256, auth.sessionKey)
}

// StoreVP stores the given VP under a new identifier or existing identifier
func (auth *Auth) StoreVP(customerID int, VP string) string {
	return auth.createSession(customerID, VP)
}

// JWTHandler is like the echo JWT middleware. It checks the JWT and required claims
func (auth *Auth) JWTHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		protectedPaths := []string{
			"/web/private",
		}
		for _, path := range protectedPaths {
			if strings.HasPrefix(ctx.Request().RequestURI, path) {
				token, err := auth.extractJWTFromHeader(ctx)
				if err != nil {
					ctx.Echo().Logger.Error(err)
					return echo.NewHTTPError(http.StatusUnauthorized, err)
				}
				if _, ok := token.Get(SessionID); !ok {
					return echo.NewHTTPError(http.StatusUnauthorized, "could not get sessionID from token")
				}

				customerId, ok := customerIDFromToken(token)
				if !ok {
					return echo.NewHTTPError(http.StatusUnauthorized, "could not get customerID from token")
				}
				ctx.Set(CustomerID, customerId)
			}
		}
		return next(ctx)
	}
}

func (auth *Auth) AuthenticatePassword(customerID int, password string) (string, error) {
	_, err := auth.customers.FindByID(customerID)
	if err != nil {
		return "", errors.New("invalid customer ID")
	}
	if auth.password != password {
		return "", errors.New("authentication failed")
	}
	token := auth.createSession(customerID, fmt.Sprintf("%d%s", customerID, password))
	return token, nil
}

func (auth *Auth) createSession(customerID int, credential interface{}) string {
	auth.mux.Lock()
	defer auth.mux.Unlock()

	for k, v := range auth.sessions {
		if v.credential == credential {
			return k
		}
	}
	tokenBytes := make([]byte, 64)
	_, _ = rand.Read(tokenBytes)

	token := hex.EncodeToString(tokenBytes)
	auth.sessions[token] = Session{
		credential: credential,
		startTime:  time.Now(),
		customerID: customerID,
	}
	return token
}

func (auth *Auth) ValidateJWT(token []byte) (jwt.Token, error) {
	pubKey := auth.sessionKey.PublicKey
	t, err := jwt.Parse(token, jwt.WithVerify(jwa.ES256, pubKey), jwt.WithValidate(true))
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (auth *Auth) extractJWTFromHeader(ctx echo.Context) (jwt.Token, error) {
	bearerToken := ctx.Request().Header.Get(echo.HeaderAuthorization)
	if bearerToken == "" {
		return nil, errors.New("no bearer token")
	}
	return auth.ValidateJWT([]byte(bearerToken[7:]))
}
