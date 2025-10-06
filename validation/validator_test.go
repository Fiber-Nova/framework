package validation

import (
	"testing"
)

func TestRequired(t *testing.T) {
	v := New()
	v.Required("name", "")
	
	if !v.Fails() {
		t.Error("Expected validation to fail for empty required field")
	}
	
	v2 := New()
	v2.Required("name", "John")
	
	if v2.Fails() {
		t.Error("Expected validation to pass for non-empty required field")
	}
}

func TestEmail(t *testing.T) {
	v := New()
	v.Email("email", "invalid")
	
	if !v.Fails() {
		t.Error("Expected validation to fail for invalid email")
	}
	
	v2 := New()
	v2.Email("email", "test@example.com")
	
	if v2.Fails() {
		t.Error("Expected validation to pass for valid email")
	}
}

func TestMin(t *testing.T) {
	v := New()
	v.Min("password", "12", 3)
	
	if !v.Fails() {
		t.Error("Expected validation to fail for string shorter than minimum")
	}
	
	v2 := New()
	v2.Min("password", "123", 3)
	
	if v2.Fails() {
		t.Error("Expected validation to pass for string meeting minimum length")
	}
}

func TestMax(t *testing.T) {
	v := New()
	v.Max("name", "12345", 3)
	
	if !v.Fails() {
		t.Error("Expected validation to fail for string exceeding maximum")
	}
	
	v2 := New()
	v2.Max("name", "123", 3)
	
	if v2.Fails() {
		t.Error("Expected validation to pass for string within maximum length")
	}
}

func TestBetween(t *testing.T) {
	v := New()
	v.Between("username", "ab", 3, 10)
	
	if !v.Fails() {
		t.Error("Expected validation to fail for string shorter than minimum")
	}
	
	v2 := New()
	v2.Between("username", "12345678901", 3, 10)
	
	if !v2.Fails() {
		t.Error("Expected validation to fail for string longer than maximum")
	}
	
	v3 := New()
	v3.Between("username", "john", 3, 10)
	
	if v3.Fails() {
		t.Error("Expected validation to pass for string within range")
	}
}

func TestMultipleValidations(t *testing.T) {
	v := New()
	v.Required("email", "invalid").Email("email", "invalid")
	
	errors := v.Errors()
	if len(errors["email"]) != 1 {
		t.Errorf("Expected 1 error for email field, got %d", len(errors["email"]))
	}
	
	v2 := New()
	v2.Required("email", "").Email("email", "")
	
	errors2 := v2.Errors()
	if len(errors2["email"]) != 2 {
		t.Errorf("Expected 2 errors for email field, got %d", len(errors2["email"]))
	}
}
