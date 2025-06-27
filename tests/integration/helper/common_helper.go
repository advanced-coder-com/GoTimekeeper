package integration_test_helper

import (
	"bytes"
	"encoding/json"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"testing"
)

type TestingContext struct {
	Email     string
	Password  string
	AuthToken string
	ProjectID []uint64
	TaskID    []uint64
}

func InitConfig(env string) {
	viper.SetConfigFile(env)
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No .env.test file found or error reading it: %v", err)
	}
}

var ErrorMessage struct {
	ErrorMessage string `json:"error"`
}

func DoPost(t *testing.T, client *http.Client, url string, body any) *http.Response {
	buf, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(buf))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("HTTP request failed: %v", err)
	}
	return resp
}

func DoPostAuth(t *testing.T, client *http.Client, url string, body any, token string) *http.Response {
	buf, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(buf))
	if err != nil {
		t.Fatalf("failed to create auth request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("auth HTTP request failed: %v", err)
	}
	return resp
}

func DoPutchAuth(t *testing.T, client *http.Client, url string, body any, token string) *http.Response {
	buf, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(buf))
	if err != nil {
		t.Fatalf("failed to create auth request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("auth HTTP request failed: %v", err)
	}
	return resp
}

func DoDeleteAuth(t *testing.T, client *http.Client, url string, body any, token string) *http.Response {
	buf, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewReader(buf))
	if err != nil {
		t.Fatalf("failed to create auth request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("auth HTTP request failed: %v", err)
	}
	return resp
}
func DoGet(t *testing.T, client *http.Client, url string) *http.Response {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("failed to create GET request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}
	return resp
}

func DoGetAuth(t *testing.T, client *http.Client, url string, token string) *http.Response {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("failed to create GET request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("authenticated GET request failed: %v", err)
	}
	return resp
}

func DecodeJSON(t *testing.T, r io.Reader, v any) {
	if err := json.NewDecoder(r).Decode(v); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
}
