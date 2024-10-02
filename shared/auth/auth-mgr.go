package auth

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rocket-pool/node-manager-core/api/types"
	"github.com/rocket-pool/node-manager-core/log"
)

const (
	AuthorizationHeader string = "Authorization"
	BearerPrefix        string = "Bearer "

	// How long a JWT is valid for after sending a request
	DefaultRequestLifespan time.Duration = time.Second * 5
)

// Manager for API authorization
type AuthorizationManager struct {
	// The path to the API authorization key file
	keyPath string

	// The API authorization key
	key []byte

	// The JWT signing method
	signingMethod jwt.SigningMethod

	// The amount of time a request is valid for
	requestLifespan time.Duration

	// The name to use for the issuer (for source tracing in logs)
	clientName string
}

// Creates a new API authorization manager.
// Note that the key is not loaded until one of the load methods is called or it's lazy loaded via AddAuthHeader.
func NewAuthorizationManager(keyPath string, clientName string, requestLifespan time.Duration) *AuthorizationManager {
	return &AuthorizationManager{
		keyPath:         keyPath,
		signingMethod:   jwt.SigningMethodHS384,
		requestLifespan: requestLifespan,
		clientName:      clientName,
	}
}

// Sets the API authorization key directly - useful for testing.
func (m *AuthorizationManager) SetKey(key []byte) {
	m.key = key
}

// Loads the provided API authorization key from disk.
func (m *AuthorizationManager) LoadAuthKey() error {
	// Read the file
	keyData, err := os.ReadFile(m.keyPath)
	if err != nil {
		return fmt.Errorf("error reading API key [%s] from disk: %w", m.keyPath, err)
	}
	m.key = keyData
	return nil
}

// Adds the API authorization header to the provided request.
// If the key is not loaded, this will attempt to load it.
func (m *AuthorizationManager) AddAuthHeader(request *http.Request) error {
	// Lazy load the key
	if m.key == nil {
		err := m.LoadAuthKey()
		if err != nil {
			return fmt.Errorf("error loading API key: %w", err)
		}
		if len(m.key) == 0 {
			return errors.New("API key is empty")
		}
	}

	// Create a new token from the secret
	now := time.Now()
	claims := &jwt.RegisteredClaims{
		Issuer:    m.clientName,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(m.requestLifespan)),
	}
	token := jwt.NewWithClaims(m.signingMethod, claims)
	tokenString, err := token.SignedString(m.key)
	if err != nil {
		return fmt.Errorf("error signing API token: %w", err)
	}

	// Add the token to the request headers
	request.Header.Add(AuthorizationHeader, BearerPrefix+tokenString)
	return nil
}

// Validates the provided request by checking the authorization header.
// If the key is not loaded, this will attempt to load it.
func (m *AuthorizationManager) ValidateRequest(request *http.Request) (string, error) {
	// Lazy load the key
	if m.key == nil {
		err := m.LoadAuthKey()
		if err != nil {
			return "", fmt.Errorf("error loading API key: %w", err)
		}
		if len(m.key) == 0 {
			return "", errors.New("API key is empty")
		}
	}

	// Make sure the header exists
	header := request.Header.Get(AuthorizationHeader)
	if header == "" {
		return "", errors.New("missing authorization header")
	}

	// Check the header prefix
	if !strings.HasPrefix(header, BearerPrefix) {
		return "", errors.New("authorization header is missing the expected prefix")
	}
	tokenString := strings.TrimPrefix(header, BearerPrefix)

	// Parse the token
	var claims jwt.RegisteredClaims
	token, err := jwt.ParseWithClaims(string(tokenString), &claims, m.keyFunc, jwt.WithValidMethods(
		[]string{
			m.signingMethod.Alg(),
		},
	))

	// Return
	clientName := claims.Issuer
	if err != nil {
		return clientName, fmt.Errorf("error parsing JWT token: %w", err)
	}
	if !token.Valid {
		return clientName, errors.New("invalid JWT token")
	}
	return clientName, nil
}

// Returns a request handler that validates the request before passing it to the next handler.
func (m *AuthorizationManager) GetRequestHandler(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientName, err := m.ValidateRequest(r)
		if err != nil {
			logger.Warn("Request failed authorization",
				log.Err(err),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("remoteAddr", r.RemoteAddr),
				slog.String("clientName", clientName),
			)

			// Create the response
			msg := types.ApiResponse[any]{
				Error: fmt.Sprintf("Authorization failed (%s)", err.Error()),
			}
			bytes, _ := json.Marshal(msg)
			w.Header().Add("Content-Type", "application/json")

			// Write it
			w.WriteHeader(http.StatusUnauthorized)
			_, writeErr := w.Write(bytes)
			if writeErr != nil {
				logger.Error("Error writing auth failure response",
					log.Err(writeErr),
					slog.String("path", r.URL.Path),
					slog.String("method", r.Method),
					slog.String("remoteAddr", r.RemoteAddr),
					slog.String("clientName", clientName),
				)
			}
			return
		}

		// Valid request
		logger.Debug("Request authorized",
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
			slog.String("remoteAddr", r.RemoteAddr),
			slog.String("clientName", clientName),
		)
		next.ServeHTTP(w, r)
	})
}

// Returns the key expected for JWT signatures. Used by JWT's parser.
func (m *AuthorizationManager) keyFunc(token *jwt.Token) (any, error) {
	return m.key, nil
}
