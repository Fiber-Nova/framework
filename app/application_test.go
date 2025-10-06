package app

import (
	"testing"
)

func TestNewApplication(t *testing.T) {
	app := New(nil)
	
	if app == nil {
		t.Fatal("Expected app to be initialized")
	}
	
	if app.Fiber() == nil {
		t.Error("Expected Fiber instance to be initialized")
	}
	
	if app.Config() == nil {
		t.Error("Expected Config to be initialized")
	}
	
	if app.Container() == nil {
		t.Error("Expected Container to be initialized")
	}
}

func TestApplicationWithConfig(t *testing.T) {
	config := &Config{
		AppName:  "TestApp",
		AppEnv:   "test",
		AppDebug: false,
		AppPort:  "8080",
	}
	
	app := New(config)
	
	if app.Config().AppName != "TestApp" {
		t.Errorf("Expected AppName to be 'TestApp', got '%s'", app.Config().AppName)
	}
	
	if app.Config().AppPort != "8080" {
		t.Errorf("Expected AppPort to be '8080', got '%s'", app.Config().AppPort)
	}
}

func TestServiceContainer(t *testing.T) {
	container := &ServiceContainer{
		services: make(map[string]interface{}),
	}
	
	// Test Bind
	container.Bind("test", "value")
	
	// Test Get
	value, exists := container.Get("test")
	if !exists {
		t.Error("Expected service to exist")
	}
	
	if value != "value" {
		t.Errorf("Expected 'value', got '%v'", value)
	}
	
	// Test non-existent service
	_, exists = container.Get("nonexistent")
	if exists {
		t.Error("Expected service to not exist")
	}
}
