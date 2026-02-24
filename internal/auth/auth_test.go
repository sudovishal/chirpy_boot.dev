package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name        string
		authHeader  string
		wantToken   string
		wantErr     bool
		expectedErr string
	}{
		{
			name:       "valid bearer token",
			authHeader: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			wantToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			wantErr:    false,
		},
		{
			name:       "valid bearer token with special characters",
			authHeader: "Bearer abc123-xyz_456.789",
			wantToken:  "abc123-xyz_456.789",
			wantErr:    false,
		},
		{
			name:        "missing authorization header",
			authHeader:  "",
			wantToken:   "",
			wantErr:     true,
			expectedErr: "missing authorization header",
		},
		{
			name:        "missing Bearer scheme",
			authHeader:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			wantToken:   "",
			wantErr:     true,
			expectedErr: "invalid authorization header format",
		},
		{
			name:        "wrong scheme (Basic instead of Bearer)",
			authHeader:  "Basic dXNlcjpwYXNz",
			wantToken:   "",
			wantErr:     true,
			expectedErr: "invalid authorization header format",
		},
		{
			name:        "Bearer with no token",
			authHeader:  "Bearer",
			wantToken:   "",
			wantErr:     true,
			expectedErr: "invalid authorization header format",
		},
		{
			name:        "Bearer with empty token",
			authHeader:  "Bearer ",
			wantToken:   "",
			wantErr:     true,
			expectedErr: "invalid authorization header format",
		},
		{
			name:        "too many parts",
			authHeader:  "Bearer token extra",
			wantToken:   "",
			wantErr:     true,
			expectedErr: "invalid authorization header format",
		},
		{
			name:        "lowercase bearer",
			authHeader:  "bearer mytoken",
			wantToken:   "",
			wantErr:     true,
			expectedErr: "invalid authorization header format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := http.Header{}
			if tt.authHeader != "" {
				headers.Set("Authorization", tt.authHeader)
			}

			token, err := GetBearerToken(headers)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetBearerToken() expected error but got none")
					return
				}
				if err.Error() != tt.expectedErr {
					t.Errorf("GetBearerToken() error = %v, want %v", err.Error(), tt.expectedErr)
				}
			} else {
				if err != nil {
					t.Errorf("GetBearerToken() unexpected error = %v", err)
					return
				}
				if token != tt.wantToken {
					t.Errorf("GetBearerToken() token = %v, want %v", token, tt.wantToken)
				}
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "Correct password",
			password:      password1,
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Incorrect password",
			password:      "wrongPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Password doesn't match different hash",
			password:      password1,
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash",
			password:      password1,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "secret", time.Hour)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}
