package api

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/types"
	"github.com/nuts-foundation/nuts-demo-ehr/nuts/client/iam"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/lestrrat-go/jwx/jwt/openid"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/customers"
)

const MaxSessionAge = time.Hour
const CustomerID = "cid"
const SessionID = "sid"

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
	Presentation *iam.VerifiablePresentation
	CustomerID   string
	StartTime    time.Time
	UserInfo     UserInfo
}

type UserInfo struct {
	Identifier string
	RoleName   string
	Initials   string
	FamilyName string
}

type JWTCustomClaims struct {
	CustomerID string `json:"cis"`
	SessionID  string `json:"sid"`
	JWTStandardClaims
}

type JWTStandardClaims struct {
	Audience  string `json:"aud,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
	Id        string `json:"jti,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	Issuer    string `json:"iss,omitempty"`
	NotBefore int64  `json:"nbf,omitempty"`
	Subject   string `json:"sub,omitempty"`
}

type CreateAuthorizationRequestParams types.CreateAuthorizationRequestParams

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

// CreateCustomerJWT creates a JWT that only stores the customer ID.
func (auth *Auth) CreateCustomerJWT(customerId string) ([]byte, error) {
	t := openid.New()
	_ = t.Set(jwt.IssuedAtKey, time.Now())
	_ = t.Set(jwt.ExpirationKey, time.Now().Add(MaxSessionAge))
	_ = t.Set(CustomerID, customerId)

	return jwt.Sign(t, jwa.ES256, auth.sessionKey)
}

func (auth *Auth) GetSession(id string) *Session {
	auth.mux.RLock()
	defer auth.mux.RUnlock()

	session, ok := auth.sessions[id]
	if !ok {
		return nil
	}
	return &session
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

func (auth *Auth) GetCustomerIDFromHeader(ctx echo.Context) (string, error) {
	token, err := auth.extractJWTFromHeader(ctx)
	if err != nil {
		ctx.Echo().Logger.Error(err)
		return "", echo.NewHTTPError(http.StatusUnauthorized, err)
	}
	rawID, ok := token.Get(CustomerID)
	if !ok {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "missing customerID in token")
	}
	return rawID.(string), nil
}

// CreateSessionJWT creates a JWT with customer ID and session ID
func (auth *Auth) CreateSessionJWT(organizationName, userName string, customerId string, session string) ([]byte, error) {
	t := openid.New()
	t.Set(jwt.SubjectKey, organizationName)
	t.Set("usi", userName)
	t.Set(jwt.IssuedAtKey, time.Now())
	t.Set(jwt.ExpirationKey, time.Now().Add(MaxSessionAge))
	t.Set(CustomerID, customerId)
	t.Set(SessionID, session)

	return jwt.Sign(t, jwa.ES256, auth.sessionKey)
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

func (auth *Auth) AuthenticatePassword(customerID string, password string) (string, UserInfo, error) {
	_, err := auth.customers.FindByID(customerID)
	if err != nil {
		return "", UserInfo{}, errors.New("invalid customer ID")
	}
	if auth.password != password {
		return "", UserInfo{}, errors.New("authentication failed")
	}
	userInfo := UserInfo{
		Identifier: "t.tester@example.com",
		RoleName:   "Verpleegkundige niveau 2",
		Initials:   "T",
		FamilyName: "Tester",
	}
	token := auth.createSession(customerID, userInfo)
	return token, userInfo, nil
}

func (auth *Auth) createSession(customerID string, userInfo UserInfo) string {
	auth.mux.Lock()
	defer auth.mux.Unlock()

	tokenBytes := make([]byte, 64)
	_, _ = rand.Read(tokenBytes)

	token := hex.EncodeToString(tokenBytes)
	auth.sessions[token] = Session{
		CustomerID: customerID,
		StartTime:  time.Now(),
		UserInfo:   userInfo,
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
