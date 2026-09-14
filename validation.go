package standardgh

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

func validateStruct(s any) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var errs []string
	for _, e := range err.(validator.ValidationErrors) {
		field := e.Field()
		tag := e.Tag()
		param := e.Param()

		switch tag {
		case "required":
			errs = append(errs, fmt.Sprintf("%s is required", field))
		case "min":
			errs = append(errs, fmt.Sprintf("%s must be at least %s", field, param))
		case "max":
			errs = append(errs, fmt.Sprintf("%s must be at most %s", field, param))
		case "len":
			errs = append(errs, fmt.Sprintf("%s must be length %s", field, param))
		case "eq":
			errs = append(errs, fmt.Sprintf("%s must be equal to %s", field, param))
		case "ne":
			errs = append(errs, fmt.Sprintf("%s must not be equal to %s", field, param))
		case "gt":
			errs = append(errs, fmt.Sprintf("%s must be greater than %s", field, param))
		case "gte":
			errs = append(errs, fmt.Sprintf("%s must be greater than or equal to %s", field, param))
		case "lt":
			errs = append(errs, fmt.Sprintf("%s must be less than %s", field, param))
		case "lte":
			errs = append(errs, fmt.Sprintf("%s must be less than or equal to %s", field, param))
		case "oneof":
			errs = append(errs, fmt.Sprintf("%s must be one of [%s]", field, param))
		case "email":
			errs = append(errs, fmt.Sprintf("%s must be a valid email", field))
		case "url":
			errs = append(errs, fmt.Sprintf("%s must be a valid URL", field))
		case "uuid":
			errs = append(errs, fmt.Sprintf("%s must be a valid UUID", field))
		case "alpha":
			errs = append(errs, fmt.Sprintf("%s must contain only letters", field))
		case "alphanum":
			errs = append(errs, fmt.Sprintf("%s must contain only alphanumeric characters", field))
		case "numeric":
			errs = append(errs, fmt.Sprintf("%s must be a valid number", field))
		case "boolean":
			errs = append(errs, fmt.Sprintf("%s must be a valid boolean", field))
		case "json":
			errs = append(errs, fmt.Sprintf("%s must be a valid JSON", field))
		case "datetime":
			errs = append(errs, fmt.Sprintf("%s must be a valid datetime using format %s", field, param))
		default:
			errs = append(errs, fmt.Sprintf("%s failed on '%s' validation", field, tag))
		}
	}

	return fmt.Errorf("validation failed: %s", strings.Join(errs, ", "))
}
