package validation

import (
	"fmt"
	"regexp"
	"strings"
)

// Validator provides validation functionality
type Validator struct {
	errors map[string][]string
}

// New creates a new Validator instance
func New() *Validator {
	return &Validator{
		errors: make(map[string][]string),
	}
}

// Required validates that a field is not empty
func (v *Validator) Required(field string, value string) *Validator {
	if strings.TrimSpace(value) == "" {
		v.addError(field, fmt.Sprintf("The %s field is required", field))
	}
	return v
}

// Email validates that a field contains a valid email
func (v *Validator) Email(field string, value string) *Validator {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		v.addError(field, fmt.Sprintf("The %s field must be a valid email address", field))
	}
	return v
}

// Min validates that a string has a minimum length
func (v *Validator) Min(field string, value string, min int) *Validator {
	if len(value) < min {
		v.addError(field, fmt.Sprintf("The %s field must be at least %d characters", field, min))
	}
	return v
}

// Max validates that a string has a maximum length
func (v *Validator) Max(field string, value string, max int) *Validator {
	if len(value) > max {
		v.addError(field, fmt.Sprintf("The %s field must not exceed %d characters", field, max))
	}
	return v
}

// Between validates that a string length is between min and max
func (v *Validator) Between(field string, value string, min int, max int) *Validator {
	length := len(value)
	if length < min || length > max {
		v.addError(field, fmt.Sprintf("The %s field must be between %d and %d characters", field, min, max))
	}
	return v
}

// Fails returns true if validation has failed
func (v *Validator) Fails() bool {
	return len(v.errors) > 0
}

// Passes returns true if validation has passed
func (v *Validator) Passes() bool {
	return !v.Fails()
}

// Errors returns all validation errors
func (v *Validator) Errors() map[string][]string {
	return v.errors
}

// addError adds an error message for a field
func (v *Validator) addError(field string, message string) {
	if _, exists := v.errors[field]; !exists {
		v.errors[field] = []string{}
	}
	v.errors[field] = append(v.errors[field], message)
}
