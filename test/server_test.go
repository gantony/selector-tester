package test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

// Make sure the server is running before running this test.

var token string = "eyJhbGciOiJSUzI1NiIsImtpZCI6IktORWk1enZpeGpNQ0tSZndnMndmQUhVZzFmLXFIWG03U3REYS1iWFZiakUifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxNzQ5NzE3OTk4LCJpYXQiOjE3NDk2MzE1OTgsImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwianRpIjoiN2E1OTliZDEtOGRiZC00Mjk1LWE1YWYtOTllYWY4MjE4ZGE4Iiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJkZWZhdWx0Iiwic2VydmljZWFjY291bnQiOnsibmFtZSI6ImphbmUiLCJ1aWQiOiI5NmZjZWQzZS0yNjY1LTQ2MjQtOTE1MC1lMjA4MzIyMTRlZTAifX0sIm5iZiI6MTc0OTYzMTU5OCwic3ViIjoic3lzdGVtOnNlcnZpY2VhY2NvdW50OmRlZmF1bHQ6amFuZSJ9.XGvzR3jum7bnC3qxGmzbHVlDk9Nunl3EPjNG_Zvqf7uZuW4pZyWvnmShPy1YR36C3ILhnouEK4ysXuf50fQYw--oFGpwC_Vg56hHyg4Yomne3PCXLDD-yEB7wEMe9hp42MLfQXAPF4mENznBOaoUciluCrhPhr83PKIy1hASm8bKRHKmzOhHeQqOaiWwwSj2OZ-3nSuxh-0cpRznM_ssTVlhtUrKB8tMa0lS0r2Kt2Maa_UobeQ0rTekgJC_45HOUVEjhALTLidLMKZzMr6qFp060r0LRzy3FmHomeqmgQZCZEOT4r0aaTYrRW_qDmsbECoApEsAe7rx0-Uz3AjG3A"

func TestProxy(t *testing.T) {
	client := &http.Client{}

	payload := `{"page_size":100,"page_num":0,"sort_by":[],"time_range":{"from":"2025-06-11T08:45:59Z","to":"2025-06-11T09:00:59Z"},"selector":""}`
	// Test with valid URL and payload
	proxyReq := fmt.Sprintf(`{"url":"https://localhost:9443/tigera-elasticsearch/flowLogs/search", "payload":%q, "token": %q}`, payload, token)
	logrus.Infof("Payload: %s", proxyReq)
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/proxy", strings.NewReader(proxyReq))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Accept any 2xx status as success, since proxied endpoint may vary
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		t.Errorf("Expected status 2xx, got %s", resp.Status)
	}

	_, _ = io.ReadAll(resp.Body)
	// Optionally, check for non-empty body or specific content if you know what to expect

	// Test with invalid JSON
	req, err = http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/proxy", strings.NewReader(`{"url": "invalid"`))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %s", resp.Status)
	}
}

func TestProxyCORSPreflight(t *testing.T) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodOptions, "http://127.0.0.1:8080/proxy", nil)
	if err != nil {
		t.Fatalf("Failed to create OPTIONS request: %v", err)
	}
	req.Header.Set("Origin", "http://localhost")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type, Authorization")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make OPTIONS request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status 200 or 204, got %s", resp.Status)
	}

	if resp.Header.Get("Access-Control-Allow-Origin") == "" {
		t.Error("Missing Access-Control-Allow-Origin header in preflight response")
	}
	if resp.Header.Get("Access-Control-Allow-Headers") == "" {
		t.Error("Missing Access-Control-Allow-Headers header in preflight response")
	}
	if resp.Header.Get("Access-Control-Allow-Methods") == "" {
		t.Error("Missing Access-Control-Allow-Methods header in preflight response")
	}
}
