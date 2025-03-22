package rapidgo_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rwiteshbera/rapidgo"
)

// Setup test server
func setupTestServer() *httptest.Server {
	router := rapidgo.New()

	router.Get("/hello", func(c *rapidgo.Context) {
		c.JSON(http.StatusOK, "Hello, World!")
	})

	router.Get("/user/:id", func(c *rapidgo.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, map[string]string{"id": id})
	})

	router.Get("/user/id", func(c *rapidgo.Context) {
		c.JSON(http.StatusOK, "/user/id")
	})

	router.Get("/search", func(c *rapidgo.Context) {
		query := c.Query("q")
		c.JSON(http.StatusOK, map[string]string{"query": query})
	})

	router.Post("/submit", func(c *rapidgo.Context) {
		var data map[string]string
		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
			return
		}
		c.JSON(http.StatusOK, data)
	})

	router.Put("/update", func(c *rapidgo.Context) {
		c.JSON(http.StatusOK, map[string]string{"status": "updated"})
	})

	router.Delete("/delete", func(c *rapidgo.Context) {
		c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
	})

	router.Patch("/patch", func(c *rapidgo.Context) {
		c.JSON(http.StatusOK, map[string]string{"status": "patched"})
	})

	router.Options("/options", func(c *rapidgo.Context) {
		c.SetHeader("Allow", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.SendStatus(http.StatusNoContent)
	})

	router.Head("/head", func(c *rapidgo.Context) {
		c.SendStatus(http.StatusOK)
	})

	return httptest.NewServer(router.Router)
}

func TestRoutes(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	t.Run("Test GET /hello", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/hello")
		if err != nil {
			t.Fatalf("Failed to make GET request: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("Test GET /user/:id", func(t *testing.T) {

	})

	t.Run("Test GET /user/id", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/user/id")
		if err != nil {
			t.Fatalf("Failed to make GET request: %v", err)
		}
		defer resp.Body.Close()

		// Check status code
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}

		// Check response body
		var result string
		json.NewDecoder(resp.Body).Decode(&result)
		expectedResponse := "/user/id"
		if result != expectedResponse {
			t.Errorf("Expected response %q, got %q", expectedResponse, result)
		}
	})

	t.Run("Test GET /search?q=golang", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/search?q=golang")
		if err != nil {
			t.Fatalf("Failed to make GET request: %v", err)
		}
		defer resp.Body.Close()
		var result map[string]string
		json.NewDecoder(resp.Body).Decode(&result)
		if result["query"] != "golang" {
			t.Errorf("Expected query 'golang', got %s", result["query"])
		}
	})

	t.Run("Test POST /submit", func(t *testing.T) {
		data := `{"name": "Test"}`
		resp, err := http.Post(server.URL+"/submit", "application/json", strings.NewReader(data))
		if err != nil {
			t.Fatalf("Failed to make POST request: %v", err)
		}
		defer resp.Body.Close()
		var result map[string]string
		json.NewDecoder(resp.Body).Decode(&result)
		if result["name"] != "Test" {
			t.Errorf("Expected name Test, got %s", result["name"])
		}
	})

	methods := []struct {
		method string
		path   string
	}{
		{"PUT", "/update"},
		{"DELETE", "/delete"},
		{"PATCH", "/patch"},
	}

	for _, test := range methods {
		t.Run("Test "+test.method+" "+test.path, func(t *testing.T) {
			req, err := http.NewRequest(test.method, server.URL+test.path, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected 200, got %d", resp.StatusCode)
			}
		})
	}
}
