package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

// Note: In Go 1.20+, math/rand is automatically seeded.
// No explicit rand.Seed() needed.

// BaseURL returns the API Gateway base URL
func BaseURL() string {
	return "http://localhost:8080"
}

// makeRequest performs an HTTP request with optional JSON body and auth token
func makeRequest(method, path string, body interface{}, token string) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, BaseURL()+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	return http.DefaultClient.Do(req)
}

// parseResponse reads and unmarshals the response body.
// Note: Caller is responsible for closing resp.Body.
func parseResponse[T any](resp *http.Response) (T, error) {
	var result T
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("failed to read response: %w", err)
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return result, fmt.Errorf("failed to unmarshal response: %w (body: %s)", err, string(body))
	}
	return result, nil
}

// waitForService waits for a service to be healthy
func waitForService(url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	backoff := 500 * time.Millisecond

	for time.Now().Before(deadline) {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return nil
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(backoff)
		if backoff < 5*time.Second {
			backoff *= 2
		}
	}
	return fmt.Errorf("service not ready after %v", timeout)
}

// generateTestEmail creates a unique email for test isolation
func generateTestEmail() string {
	return fmt.Sprintf("test_%d_%d@example.com", time.Now().UnixNano(), rand.Intn(10000))
}

// generateTestName creates a unique name for test isolation
func generateTestName(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano()%100000)
}
