package routing

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestBasicRoutes(t *testing.T) {
	app := fiber.New()
	router := New(app)
	
	router.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("GET")
	})
	
	router.Post("/test", func(c *fiber.Ctx) error {
		return c.SendString("POST")
	})
	
	// Test GET
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	
	body, _ := io.ReadAll(resp.Body)
	if string(body) != "GET" {
		t.Errorf("Expected 'GET', got '%s'", string(body))
	}
	
	// Test POST
	req = httptest.NewRequest("POST", "/test", nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	
	body, _ = io.ReadAll(resp.Body)
	if string(body) != "POST" {
		t.Errorf("Expected 'POST', got '%s'", string(body))
	}
}

func TestRouteGroup(t *testing.T) {
	app := fiber.New()
	router := New(app)
	
	router.Group("/api", func(group *RouterGroup) {
		group.Get("/test", func(c *fiber.Ctx) error {
			return c.SendString("GROUP")
		})
	})
	
	req := httptest.NewRequest("GET", "/api/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	
	body, _ := io.ReadAll(resp.Body)
	if string(body) != "GROUP" {
		t.Errorf("Expected 'GROUP', got '%s'", string(body))
	}
}

func TestAllHTTPMethods(t *testing.T) {
	app := fiber.New()
	router := New(app)
	
	router.Get("/test", func(c *fiber.Ctx) error { return c.SendString("GET") })
	router.Post("/test", func(c *fiber.Ctx) error { return c.SendString("POST") })
	router.Put("/test", func(c *fiber.Ctx) error { return c.SendString("PUT") })
	router.Delete("/test", func(c *fiber.Ctx) error { return c.SendString("DELETE") })
	router.Patch("/test", func(c *fiber.Ctx) error { return c.SendString("PATCH") })
	
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	
	for _, method := range methods {
		req := httptest.NewRequest(method, "/test", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Error testing %s: %v", method, err)
		}
		
		body, _ := io.ReadAll(resp.Body)
		if string(body) != method {
			t.Errorf("Expected '%s', got '%s'", method, string(body))
		}
	}
}
