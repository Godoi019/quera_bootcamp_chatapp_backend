package fiber

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofrs/uuid/v5"
)

var validate = validator.New()

func init() {
	// Register UUID v7 validator (strict)
	_ = RegisterValidation("uuidv7", func(fl validator.FieldLevel) bool {
		f := fl.Field().Interface().(uuid.UUID)

		return !f.IsNil() && f.Version() == 7
	})

	// Register general UUID validator (any valid UUID)
	_ = RegisterValidation("uuid", func(fl validator.FieldLevel) bool {
		f := fl.Field().Interface().(uuid.UUID)

		return !f.IsNil()
	})
}

func RegisterValidation(tag string, fn validator.Func) error {
	return validate.RegisterValidation(tag, fn)
}

// getJSONFieldName extracts the JSON field name from struct tags
func getJSONFieldName(fieldName string, structType reflect.Type) string {
	if field, found := structType.FieldByName(fieldName); found {
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {
			// Handle cases like "service_id,omitempty"
			if commaIndex := strings.Index(jsonTag, ","); commaIndex != -1 {
				return jsonTag[:commaIndex]
			}
			return jsonTag
		}
	}
	// Fallback to snake_case conversion of field name
	return toSnakeCase(fieldName)
}

// toSnakeCase converts CamelCase to snake_case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

func ValidateStruct[T any](body T) []*ErrorResponse {
	var errors []*ErrorResponse

	err := validate.Struct(body)

	if err != nil {
		bodyType := reflect.TypeOf(body)
		if bodyType.Kind() == reflect.Ptr {
			bodyType = bodyType.Elem()
		}

		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse

			// Extract field name from the namespace
			namespace := err.StructNamespace()
			var fieldName string
			if idx := strings.LastIndex(namespace, "."); idx != -1 {
				fieldName = namespace[idx+1:]
			} else {
				fieldName = namespace
			}

			// Get the JSON field name using reflection
			element.Field = getJSONFieldName(fieldName, bodyType)
			element.Tag = err.Tag()

			switch element.Tag {
			case "required":
				element.Value = "This field is required"
			case "uuidv7":
				element.Value = "Must be a valid UUID version 7"
			case "uuid":
				element.Value = "Must be a valid UUID"
			default:
				element.Value = err.Param()
			}
			errors = append(errors, &element)
		}
	}
	return errors
}

func ParseRequestBody[T any](c fiber.Ctx, body *T) *ValidationError {
	errParse := c.Bind().Body(body)
	res := &ValidationError{}
	res.Message = ErrRequiredBodyNotFound

	if errParse != nil {
		res.Message = errParse
		return res
	}

	errValidation := ValidateStruct[T](*body)

	if errValidation != nil {
		res.Errors = errValidation

		return res
	}

	return nil
}

func ParseRequestParams[T any](c fiber.Ctx, params *T) *ValidationError {
	errParse := c.Bind().Query(params)
	res := &ValidationError{}
	res.Message = ErrRequiredUrlQueriesNotFound

	if errParse != nil {
		res.Message = errParse
		return res
	}

	errValidation := ValidateStruct[T](*params)

	if errValidation != nil {
		res.Errors = errValidation

		return res
	}

	return nil
}
