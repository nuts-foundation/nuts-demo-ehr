package api

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/google/uuid"
	ssi "github.com/nuts-foundation/go-did"
	"github.com/nuts-foundation/go-did/vc"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

	jwt2 "github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/lestrrat-go/jwx/jwt/openid"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/customers"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts/client/auth"
)

const MaxSessionAge = time.Hour
const CustomerID = "cid"
const SessionID = "sid"
const Elevated = "elv"

type Auth struct {
	// sessions maps session identifiers to VPs
	sessions map[string]Session
	// mux is used to secure access to internal state of this struct to prevent racy behaviour
	mux        sync.RWMutex
	customers  customers.Repository
	password   string
	sessionKey *ecdsa.PrivateKey
}

type Session struct {
	Presentation auth.VerifiablePresentation
	CustomerID   int
	StartTime    time.Time
	UserContext  bool
}

type JWTCustomClaims struct {
	CustomerID int    `json:"cis"`
	SessionID  string `json:"sid"`
	jwt2.StandardClaims
}

func NewAuth(key *ecdsa.PrivateKey, customers customers.Repository, passwd string) *Auth {
	result := &Auth{
		sessionKey: key,
		customers:  customers,
		sessions:   map[string]Session{},
		password:   passwd,
	}
	// In theory, we should manage this goroutine to have it properly shut down when the server is shut down,
	// but it's a demo app and everything is in-memory, so no chance of data corruption.
	go func(result *Auth) {
		// Clean up expired sessions to avoid memory leak when running over a long time
		ticker := time.NewTicker(time.Minute)
		for {
			<-ticker.C
			result.mux.Lock()
			for token, session := range result.sessions {
				if session.StartTime.Add(MaxSessionAge).Before(time.Now()) {
					delete(result.sessions, token)
				}
			}
			result.mux.Unlock()
		}
	}(result)
	return result
}

// CreateCustomerJWT creates a JWT that only stores the customer ID. This is required for the IRMA flow.
func (auth *Auth) CreateCustomerJWT(customerId int) ([]byte, error) {
	t := openid.New()
	t.Set(jwt.IssuedAtKey, time.Now())
	t.Set(jwt.ExpirationKey, time.Now().Add(MaxSessionAge))
	t.Set(CustomerID, customerId)

	return jwt.Sign(t, jwa.ES256, auth.sessionKey)
}

func (auth *Auth) GetSession(id string) *Session {
	auth.mux.RLock()
	defer auth.mux.RUnlock()

	for t, session := range auth.sessions {
		if t == id {
			return &session
		}
	}

	return nil
}

func (auth *Auth) GetSessions() map[string]Session {
	auth.mux.RLock()
	defer auth.mux.RUnlock()

	sessions := map[string]Session{}

	for token, session := range auth.sessions {
		sessions[token] = session
	}

	return sessions
}

func (auth *Auth) GetCustomerIDFromHeader(ctx echo.Context) (int, error) {
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
func (auth *Auth) CreateSessionJWT(organizationName, userName string, customerId int, session string, elevated bool) ([]byte, error) {
	t := openid.New()
	t.Set(jwt.SubjectKey, organizationName)
	t.Set("usi", userName)
	t.Set(jwt.IssuedAtKey, time.Now())
	t.Set(jwt.ExpirationKey, time.Now().Add(MaxSessionAge))
	t.Set(CustomerID, customerId)
	t.Set(SessionID, session)
	t.Set(Elevated, elevated)

	return jwt.Sign(t, jwa.ES256, auth.sessionKey)
}

// StoreVP stores the given VP under a new identifier or existing identifier
func (auth *Auth) StoreVP(customerID int, VP auth.VerifiablePresentation) string {
	return auth.createSession(customerID, VP, true)
}

// JWTHandler is like the echo JWT middleware. It checks the JWT and required claims
func (auth *Auth) JWTHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		protectedPaths := []string{
			"/web/private",
			"/web/auth/selfsigned/session", // self-signed auth means require authenticated session for employee info
		}
		for _, path := range protectedPaths {
			if strings.HasPrefix(ctx.Request().RequestURI, path) {
				token, err := auth.extractJWTFromHeader(ctx)
				if err != nil {
					ctx.Echo().Logger.Error(err)
					return echo.NewHTTPError(http.StatusUnauthorized, err)
				}
				if sessionID, ok := token.Get(SessionID); !ok {
					return echo.NewHTTPError(http.StatusUnauthorized, "could not get sessionID from token")
				} else {
					ctx.Set(SessionID, fmt.Sprintf("%s", sessionID))
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

func createPasswordVP(customerDID string) auth.VerifiablePresentation {
	const identifier = "t.tester@example.com"
	id, _ := uuid.NewUUID()
	return auth.VerifiablePresentation{
		VerifiableCredential: []vc.VerifiableCredential{
			{
				Type:   []ssi.URI{ssi.MustParseURI("EmployeeCredential")},
				Issuer: ssi.MustParseURI(customerDID),
				CredentialSubject: []interface{}{
					map[string]interface{}{
						"employer": customerDID,
						"employee": map[string]interface{}{
							"identifier": identifier,
							"roleName":   "Verpleegkundige niveau 2",
							"initials":   "T",
							"familyName": "Tester",
						},
					},
				},
			},
		},
		Proof: []interface{}{map[string]interface{}{
			"challenge": id.String(), // VP needs to be unique, otherwise sessions might be shared/reused
			"identity":  fmt.Sprintf("%s/%s", customerDID, identifier),
			"created":   time.Now().Format(time.RFC3339),
		}},
	}
}

func (auth *Auth) AuthenticatePassword(customerID int, password string) (string, error) {
	customer, err := auth.customers.FindByID(customerID)
	if err != nil {
		return "", errors.New("invalid customer ID")
	}
	if auth.password != password {
		return "", errors.New("authentication failed")
	}
	token := auth.createSession(customerID, createPasswordVP(*customer.Did), false)
	return token, nil
}

func (auth *Auth) createSession(customerID int, presentation auth.VerifiablePresentation, userContext bool) string {
	auth.mux.Lock()
	defer auth.mux.Unlock()

	for k, v := range auth.sessions {
		if reflect.DeepEqual(v.Presentation, presentation) {
			return k
		}
	}

	tokenBytes := make([]byte, 64)
	_, _ = rand.Read(tokenBytes)

	token := hex.EncodeToString(tokenBytes)
	auth.sessions[token] = Session{
		Presentation: presentation,
		StartTime:    time.Now(),
		CustomerID:   customerID,
		UserContext:  userContext,
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
