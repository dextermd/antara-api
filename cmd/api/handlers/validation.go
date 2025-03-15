package handlers

import (
	"antara-api/common"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"reflect"
	"strings"
)

func (h *Handler) ValidateBodyRequest(c echo.Context, payload any) []*common.ValidationError {
	var errors []*common.ValidationError

	var validate *validator.Validate
	validate = validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(payload)
	if err != nil {
		reflected := reflect.ValueOf(payload)
		for _, e := range err.(validator.ValidationErrors) {
			field, _ := reflected.Type().FieldByName(e.StructField())

			key := field.Tag.Get("json")
			if key == "" {
				key = strings.ToLower(e.StructField())
			}

			condition := e.ActualTag()
			keyToTitleCase := strings.Replace(key, "_", " ", -1)
			param := e.Param()

			errMessage := fmt.Sprintf("%s field must be %s", keyToTitleCase, condition)
			if param != "" {
				errMessage += fmt.Sprintf(" %s", param)
			}

			switch condition {
			case "required":
				errMessage = fmt.Sprintf("%s field is required", keyToTitleCase)
			case "email":
				errMessage = fmt.Sprintf("%s field must be a valid email address", keyToTitleCase)
			case "password":
				errMessage = fmt.Sprintf("%s field must be a valid password", keyToTitleCase)
			case "min":

				errMessage = fmt.Sprintf("%s field must be at least %s characters", keyToTitleCase, strings.ToLower(param))
			case "max":

				errMessage = fmt.Sprintf("%s field must be at most %s characters", keyToTitleCase, strings.ToLower(param))
			case "eqfield":
				errMessage = fmt.Sprintf("%s field must be equal to %s", keyToTitleCase, strings.ToLower(param))
			}

			validationError := &common.ValidationError{
				Error:     errMessage,
				Key:       key,
				Condition: condition,
			}
			errors = append(errors, validationError)
		}
	}

	return errors
}
