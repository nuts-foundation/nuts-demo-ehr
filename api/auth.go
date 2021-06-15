package api

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/nuts-foundation/nuts-demo-ehr/domain/customers"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-demo-ehr/client"
)

const MaxSessionAge = time.Hour

type Auth struct {
	// sessions maps session cookies to base64 encoded VPs
	sessions   map[string]Session
	// mux is used to secure access to internal state of this struct to prevent racy behaviour
	mux        sync.Mutex
	nodeClient client.HTTPClient
	customers  customers.Repository
	password   string
}

type Session struct {
	credential interface{}
	startTime  time.Time
	customerID string
}

func NewAuth(nodeClient client.HTTPClient, customers customers.Repository, passwd string) *Auth {
	return &Auth{
		nodeClient: nodeClient,
		customers:  customers,
		sessions:   map[string]Session{},
		password:   passwd,
	}
}

// StoreVP stores the given VP under a new identifier or existing identifier
func (auth *Auth) StoreVP(VP string) string {
	// TODO: Resolve customer ID from vp
	return auth.createSession("", VP)
}

func (auth *Auth) VPHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		protectedPaths := []string{
			"/web/private",
		}
		sessions := auth.getSessions()
		for _, path := range protectedPaths {
			if strings.HasPrefix(ctx.Request().RequestURI, path) {
				// check cookie
				authorized := false
				token := ""
				for _, c := range ctx.Cookies() {
					if c.Name == "session" {
						token = c.Value
						_, ok := sessions[token]
						if !ok {
							return ctx.NoContent(http.StatusForbidden)
						}
						authorized = true
					}
				}
				if !authorized {
					return ctx.NoContent(http.StatusForbidden)
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

// getSessions must be used to for read access to sessions, because it cleans up old, expired sessions.
func (auth *Auth) getSessions() map[string]Session {
	auth.mux.Lock()
	defer auth.mux.Unlock()

	sessions := make(map[string]Session, 0)
	for key, session := range auth.sessions {
		if session.startTime.Add(MaxSessionAge).Before(time.Now()) {
			logrus.Infof("Session %s expired, cleaning up.", key)
			delete(auth.sessions, key)
			continue
		}
		sessions[key] = session
	}
	return sessions
}

// StoreVP stores the given VP under a new identifier or existing identifier
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
