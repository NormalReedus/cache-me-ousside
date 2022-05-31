package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/magnus-bb/cache-me-ousside/cache"
)

// newConfigValidator returns a new instance of config validator
// with custom validation rules that are used for the Config.
func newConfigValidator() *validator.Validate {
	validate := validator.New()

	// Checks if a string is a path containing a valid, existing directory
	validate.RegisterValidation("filepath", func(fl validator.FieldLevel) bool {
		path := fl.Field().String()

		directory := filepath.Dir(path)

		if _, err := os.Stat(directory); os.IsNotExist(err) {
			return false
		}

		return true
	})

	// Checks if a string is a valid route indentifier
	validate.RegisterValidation("route", func(fl validator.FieldLevel) bool {
		route := fl.Field().String()

		match, err := regexp.MatchString(`^/[\w\-\._~:/?#[\]@!\$&'\(\)\*\+,;=.]+$|^\*$`, route)

		return match && err == nil
	})

	return validate
}

// formatValidationError aggregates validator.ValidationErrors into a human-readable error
// that can be used to display a single validation error to the user.
func formatValidationError(validationErrors validator.ValidationErrors) error {
	var errorMessages []string

	for _, err := range validationErrors {
		fmt.Println("Namespace:", err.Namespace())
		fmt.Println("StructNamespace:", err.StructNamespace())
		fmt.Println("Field:", err.Field())
		fmt.Println("StructField:", err.StructField())
		fmt.Println("Tag:", err.Tag())
		fmt.Println("ActualTag:", err.ActualTag())
		fmt.Println("Value:", err.Value())
		fmt.Println("Param:", err.Param())
		fmt.Println("Kind:", err.Kind())
		fmt.Println("Type:", err.Type())
		fmt.Println("Error:", err.Error())
		fmt.Println()

		fieldName := err.Field()
		errorMessages = append(errorMessages, validationErrorMap[fieldName](err))
	}

	aggregatedErrors := errors.New(strings.Join(errorMessages, "; "))

	return fmt.Errorf("config validation failed: %w", aggregatedErrors)
}

// requiredErrorMsg returns a string formatted to explain a missing configuration property.
func requiredErrorMsg(prop string) string {
	return fmt.Sprintf("configuration is missing the '%s' property", prop)
}

// validationErrorMap maps field names to functions that return a validation error message
// depending on the validation details given with the validator.FieldError err.
var validationErrorMap = map[string]func(err validator.FieldError) string{
	"Capacity": func(err validator.FieldError) string {
		return fmt.Sprintf("'Capacity' must be a non-zero uint64, it is %d", err.Value())
	},

	"CapacityUnit": func(err validator.FieldError) string {
		units := cache.ValidCapacityUnits

		firstUnits := units[:len(units)-1]

		firstUnitsQuoted := make([]string, len(firstUnits))
		for i, unit := range firstUnits {
			firstUnitsQuoted[i] = fmt.Sprintf("%q", unit)
		}

		firstUnitsString := strings.Join(firstUnitsQuoted, ", ")
		lastUnitString := fmt.Sprintf("%q", units[len(units)-1])

		return fmt.Sprintf("'CapacityUnit' must be omitted or set to either %s, or %s, it is %q", firstUnitsString, lastUnitString, err.Value())
	},

	"Hostname": func(err validator.FieldError) string {
		return fmt.Sprintf("'Hostname' must be omitted or set to a valid rfc1123 hostname, it is %q", err.Value())
	},

	"Port": func(err validator.FieldError) string {
		return fmt.Sprintf("'Port' must be omitted or set to a number between 1 and 65535, it is %d", err.Value())
	},

	"ApiUrl": func(err validator.FieldError) string {
		tag := err.Tag()

		if tag == "required" {
			return requiredErrorMsg("ApiUrl")

		} else if tag == "url" {
			return fmt.Sprintf("'ApiUrl' value %q is not a valid URL", err.Value())
		}

		return "" // should never happen
	},

	//TODO: cache and bust
}
