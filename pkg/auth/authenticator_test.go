package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewAuthenticator(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		wantErr     bool
		errContains string
	}{
		{
			name: "valid config",
			config: &Config{
				Server: "https://api.example.com",
			},
			wantErr: false,
		},
		{
			name:        "empty config",
			config:      nil,
			wantErr:    true,
			errContains: "server URL is required",
		},
		{
			name: "no server URL",
			config: &Config{
				Server: "",
			},
			wantErr:    true,
			errContains: "server URL is required",
		},
		{
			name: "invalid server URL",
			config: &Config{
				Server: "://invalid-url",
			},
			wantErr:    true,
			errContains: "invalid server URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth, err := NewAuthenticator(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAuthenticator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && tt.errContains != "" {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("NewAuthenticator() error = %v, want error containing %v", err, tt.errContains)
				}
			}
			if !tt.wantErr && auth == nil {
				t.Error("NewAuthenticator() returned nil authenticator")
			}
		})
	}
}

func TestAuthenticator_Authenticate(t *testing.T) {
	tests := []struct {
		name           string
		username       string
		password       string
		serverResponse *AuthResponse
		serverStatus   int
		serverError    error
		wantToken      string
		wantErr        bool
		errContains    string
	}{
		{
			name:     "successful authentication",
			username: "testuser",
			password: "testpass",
			serverResponse: &AuthResponse{
				Token: "valid-token",
			},
			serverStatus: http.StatusOK,
			wantToken:    "valid-token",
		},
		{
			name:         "empty username",
			username:     "",
			password:     "testpass",
			wantErr:      true,
			errContains:  "username and password are required",
		},
		{
			name:         "empty password",
			username:     "testuser",
			password:     "",
			wantErr:      true,
			errContains:  "username and password are required",
		},
		{
			name:         "server error",
			username:     "testuser",
			password:     "testpass",
			serverStatus: http.StatusInternalServerError,
			wantErr:      true,
			errContains:  "authentication failed with status code: 500",
		},
		{
			name:     "empty token response",
			username: "testuser",
			password: "testpass",
			serverResponse: &AuthResponse{
				Token: "",
			},
			serverStatus: http.StatusOK,
			wantErr:      true,
			errContains:  "received empty token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request method
				if r.Method != http.MethodPost {
					t.Errorf("Expected POST request, got %s", r.Method)
				}

				// Verify Content-Type
				if r.Header.Get("Content-Type") != "application/json" {
					t.Errorf("Expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
				}

				// If status code is set, return it
				if tt.serverStatus != 0 {
					w.WriteHeader(tt.serverStatus)
				}

				// If response data is set, return it
				if tt.serverResponse != nil {
					json.NewEncoder(w).Encode(tt.serverResponse)
				}
			}))
			defer server.Close()

			// Create authenticator
			config := &Config{
				Server:   server.URL,
				Timeout: 5 * time.Second,
			}
			auth, err := NewAuthenticator(config)
			if err != nil {
				t.Fatalf("Failed to create authenticator: %v", err)
			}

			// Execute authentication
			token, err := auth.Authenticate(tt.username, tt.password)

			// Verify results
			if (err != nil) != tt.wantErr {
				t.Errorf("Authenticate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && tt.errContains != "" {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Authenticate() error = %v, want error containing %v", err, tt.errContains)
				}
			}
			if !tt.wantErr && token != tt.wantToken {
				t.Errorf("Authenticate() token = %v, want %v", token, tt.wantToken)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	if config == nil {
		t.Error("DefaultConfig() returned nil")
	}
	if config.AuthPath != "/auth" {
		t.Errorf("DefaultConfig().AuthPath = %v, want /auth", config.AuthPath)
	}
	if config.Timeout != 10*time.Second {
		t.Errorf("DefaultConfig().Timeout = %v, want 10s", config.Timeout)
	}
} 