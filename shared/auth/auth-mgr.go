package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	AuthorizationHeader string = "Authorization"
	BearerPrefix        string = "Bearer "

	// How long a JWT is valid for after sending a request
	ValidTime time.Duration = time.Second * 5
)

// Manager for API authorization
type AuthorizationManager struct {
	// The path to the API authorization key file
	keyPath string

	// The API authorization key
	key []byte

	// The JWT signing method
	signingMethod jwt.SigningMethod
}

// Creates a new API authorization manager.
// Note that the key is not loaded until one of the load methods is called or it's lazy loaded via AddAuthHeader.
func NewAuthorizationManager(keyPath string) *AuthorizationManager {
	return &AuthorizationManager{
		keyPath:       keyPath,
		signingMethod: jwt.SigningMethodHS384,
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
	}

	// Create a new token from the secret
	now := time.Now()
	claims := &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(ValidTime)),
	}
	token := jwt.NewWithClaims(m.signingMethod, claims)
	tokenString, err := token.SignedString(m.key)
	if err != nil {
		return fmt.Errorf("error signing API token: %w", err)
	}

	// Add the token to the request header
	request.Header.Add(AuthorizationHeader, BearerPrefix+tokenString)
	return nil
}

// Validates the provided request by checking the authorization header.
// If the key is not loaded, this will attempt to load it.
func (m *AuthorizationManager) ValidateRequest(request *http.Request) error {
	// Lazy load the key
	if m.key == nil {
		err := m.LoadAuthKey()
		if err != nil {
			return fmt.Errorf("error loading API key: %w", err)
		}
	}

	// Make sure the header exists
	header := request.Header.Get(AuthorizationHeader)
	if header == "" {
		return errors.New("missing authorization header")
	}

	// Check the header prefix
	if !strings.HasPrefix(header, BearerPrefix) {
		return errors.New("authorization header is missing the expected prefix")
	}
	tokenString := strings.TrimPrefix(header, BearerPrefix)

	// Parse the token
	token, err := jwt.Parse(string(tokenString), m.keyFunc, jwt.WithValidMethods(
		[]string{
			m.signingMethod.Alg(),
		},
	))
	if err != nil {
		return fmt.Errorf("error parsing JWT token: %w", err)
	}
	if !token.Valid {
		return errors.New("invalid JWT token")
	}
	return nil
}

// Returns the key expected for JWT signatures. Used by JWT's parser.
func (m *AuthorizationManager) keyFunc(token *jwt.Token) (any, error) {
	return m.key, nil
}
