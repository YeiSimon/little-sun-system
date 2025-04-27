package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
)

// GoogleOAuthConfig holds the configuration for Google OAuth
type GoogleAuth struct {
	Config *oauth2.Config
}

// NewGoogleAuth creates a new Google auth handler
func NewGoogleAuth(configPath string) (*GoogleAuth, error) {
	// Read client secret file
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %w", err)
	}

	// Configure OAuth2
	config, err := google.ConfigFromJSON(b,
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile")
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %w", err)
	}

	return &GoogleAuth{
		Config: config,
	}, nil
}

// ValidateIDToken validates a Google ID token and returns the claims
func (g *GoogleAuth) ValidateIDToken(ctx context.Context, idToken string) (map[string]interface{}, error) {
	// Validate the ID token
	payload, err := idtoken.Validate(ctx, idToken, g.Config.ClientID)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Additional validation
	// 1. Verify iss is accounts.google.com
	if iss, ok := payload.Claims["iss"].(string); !ok || iss != "https://accounts.google.com" {
		return nil, fmt.Errorf("invalid token issuer: %v", payload.Claims["iss"])
	}

	// 2. Validate aud matches client ID
	if aud, ok := payload.Claims["aud"].(string); !ok || aud != g.Config.ClientID {
		return nil, fmt.Errorf("invalid audience: %v", payload.Claims["aud"])
	}

	// 3. Check email_verified
	if emailVerified, ok := payload.Claims["email_verified"].(bool); !ok || !emailVerified {
		return nil, fmt.Errorf("email not verified")
	}

	// Return the claims
	return payload.Claims, nil
}

// GetClientID extracts the client ID from the config file
func GetClientID(configPath string) (string, error) {
	type WebConfig struct {
		ClientID string `json:"client_id"`
	}

	type Config struct {
		Web WebConfig `json:"web"`
	}

	// Read client secret file
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("unable to read client secret file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(b, &config); err != nil {
		return "", fmt.Errorf("unable to parse client secret file: %w", err)
	}

	return config.Web.ClientID, nil
}