package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register custom tag name function to use json tags
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// ValidateStruct validates a struct and returns validation errors
func ValidateStruct(s interface{}) map[string]interface{} {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	errors := make(map[string]interface{})
	
	for _, err := range err.(validator.ValidationErrors) {
		fieldName := err.Field()
		
		switch err.Tag() {
		case "required":
			errors[fieldName] = fmt.Sprintf("%s is required", fieldName)
		case "email":
			errors[fieldName] = fmt.Sprintf("%s must be a valid email address", fieldName)
		case "min":
			errors[fieldName] = fmt.Sprintf("%s must be at least %s characters long", fieldName, err.Param())
		case "max":
			errors[fieldName] = fmt.Sprintf("%s must be at most %s characters long", fieldName, err.Param())
		case "url":
			errors[fieldName] = fmt.Sprintf("%s must be a valid URL", fieldName)
		case "oneof":
			errors[fieldName] = fmt.Sprintf("%s must be one of: %s", fieldName, err.Param())
		default:
			errors[fieldName] = fmt.Sprintf("%s is invalid", fieldName)
		}
	}

	return errors
}

// IsValid checks if a struct is valid
func IsValid(s interface{}) bool {
	return validate.Struct(s) == nil
}