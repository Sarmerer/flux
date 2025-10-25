package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/flow/internal/errors"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) error {
	if email == "" {
		return errors.NewValidationError("Email is required").WithField("email")
	}

	email = strings.TrimSpace(email)
	if !emailRegex.MatchString(email) {
		return errors.NewValidationError("Invalid email format").WithField("email")
	}

	return nil
}

type PasswordValidationResult struct {
	Valid   bool
	Errors  []string
	Details map[string]bool
}

func ValidatePassword(password string) error {
	if password == "" {
		return errors.NewValidationError("Password is required").WithField("password")
	}

	var validationErrors []string

	if len(password) < 8 {
		validationErrors = append(validationErrors, "must be at least 8 characters long")
	}

	hasUpper := false
	for _, char := range password {
		if char >= 'A' && char <= 'Z' {
			hasUpper = true
			break
		}
	}
	if !hasUpper {
		validationErrors = append(validationErrors, "must contain at least one uppercase letter")
	}

	hasLower := false
	for _, char := range password {
		if char >= 'a' && char <= 'z' {
			hasLower = true
			break
		}
	}
	if !hasLower {
		validationErrors = append(validationErrors, "must contain at least one lowercase letter")
	}

	hasNumber := false
	for _, char := range password {
		if char >= '0' && char <= '9' {
			hasNumber = true
			break
		}
	}
	if !hasNumber {
		validationErrors = append(validationErrors, "must contain at least one number")
	}

	if len(validationErrors) > 0 {
		message := fmt.Sprintf("Password %s", strings.Join(validationErrors, ", "))
		return errors.NewValidationError(message).WithField("password")
	}

	return nil
}

func ValidateRequired(value string, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return errors.NewValidationError(fmt.Sprintf("%s is required", fieldName)).WithField(fieldName)
	}
	return nil
}

func ValidateMinLength(value string, minLength int, fieldName string) error {
	if len(value) < minLength {
		return errors.NewValidationError(
			fmt.Sprintf("%s must be at least %d characters long", fieldName, minLength),
		).WithField(fieldName)
	}
	return nil
}

func ValidateMaxLength(value string, maxLength int, fieldName string) error {
	if len(value) > maxLength {
		return errors.NewValidationError(
			fmt.Sprintf("%s must not exceed %d characters", fieldName, maxLength),
		).WithField(fieldName)
	}
	return nil
}
