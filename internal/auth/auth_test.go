package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	password1 := "Password123"
	password2 := "Cucumber1ee4"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrong",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Doesn't match hashes",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "hahahash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestMakeValidateJWT(t *testing.T) {
	interval := 1 * time.Second
	userID := uuid.New()
	//check if it is possible to create a token
	jwt, err := MakeJWT(userID, "secret", time.Second*1)
	if err != nil {
		t.Errorf("expected to create token")
		return
	}
	//check retrieve a id
	id, err := ValidateJWT(jwt, "secret")
	if err != nil {
		t.Errorf("expected to retrieve id")
		return
	}
	//check retrieve an expected id
	if id != userID {
		t.Errorf("expected to retrieve the same id")
		return
	}
	//check with a wrong secret token
	_, err = ValidateJWT(jwt, "notsecret")
	if err == nil {
		t.Errorf("expected to get an error")
		return
	}
	time.Sleep(interval)
	//check retrieving after expiration
	_, err = ValidateJWT(jwt, "secret")
	if err == nil {
		t.Errorf("expected to get an error")
		return
	}
}

func TestGetBearer(t *testing.T) {
	header1 := http.Header{}
	token := "5hz86ha0bbfr"
	header1.Add("Authorization", "Bearer "+token)
	header2 := http.Header{}
	tests := []struct {
		name    string
		header  http.Header
		token   string
		wantErr bool
	}{
		{
			name:    "Header exists",
			header:  header1,
			token:   token,
			wantErr: false,
		},
		{
			name:    "Header doesn't exist",
			header:  header2,
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetBearerToken(tt.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if token != tt.token {
				t.Errorf("GetBearerToken() token = %v, actual = %v", token, tt.token)
			}
		})
	}
}
