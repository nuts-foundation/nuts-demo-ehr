package api

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	jwt2 "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/lestrrat-go/jwx/jwt/openid"
	"github.com/nuts-foundation/nuts-demo-ehr/client"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/customers"
)

const MaxSessionAge = time.Hour
const CustomerID = "cid"
const SessionID = "sid"

type Auth struct {
	// sessions maps session identifiers to base64 encoded VPs
	sessions     map[string]Session
	// mux is used to secure access to internal state of this struct to prevent racy behaviour
	mux          sync.Mutex
	nodeClient   client.HTTPClient
	customers    customers.Repository
	password     string
	sessionKey   *ecdsa.PrivateKey

}

type Session struct {
	credential interface{}
	customerID string
	startTime  time.Time
}

type JWTCustomClaims struct {
	CustomerID string `json:"cis"`
	SessionID  string `json:"sid"`
	jwt2.StandardClaims
}

func NewAuth(key *ecdsa.PrivateKey, nodeClient client.HTTPClient, customers customers.Repository, passwd string) *Auth {
	return &Auth{
		sessionKey:   key,
		nodeClient: nodeClient,
		customers:  customers,
		sessions:   map[string]Session{},
		password:   passwd,
	}
}

// CreateCustomerJWT creates a JWT that only stores the customer ID. This is required for the IRMA flow.
func (auth *Auth) CreateCustomerJWT(customerId string) ([]byte, error) {
	t := openid.New()
	t.Set(jwt.IssuedAtKey, time.Now())
	t.Set(jwt.ExpirationKey, time.Now().Add(MaxSessionAge))
	t.Set(CustomerID, customerId)

	signed, err := jwt.Sign(t, jwa.ES256, auth.sessionKey)
	if err != nil {
		return nil, err
	}
	return signed, nil
}

// CreateSessionJWT creates a JWT with customer ID and session ID
func (auth *Auth) CreateSessionJWT(customerId string, session string) ([]byte, error) {
	t := openid.New()
	t.Set(jwt.IssuedAtKey, time.Now())
	// session is valid for 20 minutes
	t.Set(jwt.ExpirationKey, time.Now().Add(20*time.Minute))
	t.Set(CustomerID, customerId)
	t.Set(SessionID, session)

	signed, err := jwt.Sign(t, jwa.ES256, auth.sessionKey)
	if err != nil {
		return nil, err
	}
	return signed, nil
}

// StoreVP stores the given VP under a new identifier or existing identifier
func (auth *Auth) StoreVP(customerID string, VP string) string {
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
				bearerToken := ctx.Request().Header.Get(echo.HeaderAuthorization)
				if bearerToken == "" {
					return ctx.NoContent(http.StatusUnauthorized)
				}
				token, err := auth.ValidateJWT([]byte(bearerToken[7:]))
				if err != nil {
					ctx.Echo().Logger.Error(err)
					return ctx.NoContent(http.StatusUnauthorized)
				}
				if _, ok := token.Get(SessionID); !ok {
					return ctx.NoContent(http.StatusUnauthorized)
				}

				if customerId, ok := token.Get(CustomerID); !ok {
					return ctx.NoContent(http.StatusUnauthorized)
				} else {
					ctx.Set(CustomerID, customerId)
				}
			}
		}
		return next(ctx)
	}
}

func (auth *Auth) AuthenticatePassword(customerID string, password string) (string, error) {
	_, err := auth.customers.FindByID(customerID)
	if err != nil {
		return "", errors.New("invalid customer ID")
	}
	if auth.password != password {
		return "", errors.New("authentication failed")
	}
	token := auth.createSession(customerID, customerID+password)
	return token, nil
}

func (auth *Auth) createSession(customerID string, credential interface{}) string {
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
