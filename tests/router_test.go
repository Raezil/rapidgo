package rapidgo_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rwiteshbera/rapidgo"
)

func TestBasicRoutes(t *testing.T) {
	engine := rapidgo.New()

	// Define a simple route
	engine.Get("/hello", func(c *rapidgo.Context) {
		c.JSON(http.StatusOK, "Hello, World!")
	})

	// Define a route with parameters
	engine.Get("/user/:id", func(c *rapidgo.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, map[string]string{"id": id})
	})

	// Create a test server
	server := httptest.NewServer(engine.Router)
	defer server.Close()

	t.Run("Test /hello route", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/hello")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Check status code
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		// Check response body
		body := new(bytes.Buffer)
		body.ReadFrom(resp.Body)
		if !strings.Contains(body.String(), `"Hello, World!"`) {
			t.Errorf("Unexpected response body: %s", body.String())
		}
	})

	t.Run("Test /user/:id route", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/user/123")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		// Check status code
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		// Parse JSON response
		var result map[string]string
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			t.Fatalf("Failed to decode JSON response: %v", err)
		}

		// Check response content
		expected := map[string]string{"id": "123"}
		if result["id"] != expected["id"] {
			t.Errorf("Unexpected JSON response: got %v, want %v", result, expected)
		}
	})
}
