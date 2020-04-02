package validator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

type Validator interface {
	Exec(s interface{}) map[string][]string
}

type MyValidator struct {
	*validator.Validate //combination package by anonymous
}

func NewValidator() Validator {
	return &MyValidator{validator.New()}
}

func(v *MyValidator) Exec(s interface{}) map[string][]string{
	return errorMessage(v.Struct(s))
}

//custom error message
func errorMessage(err error) map[string][]string {
	errorMessages := make(map[string][]string)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var errorMessage string
			switch err.ActualTag() {
			case "required":
				errorMessage = fmt.Sprintf("%s  is a required field.", err.Field())
			case "min":
				errorMessage = fmt.Sprintf("%s  must be at least %s character in length.", err.Field(), err.Param())
			case "max":
				errorMessage = fmt.Sprintf("%s   must be a maximum of %s characters in length.", err.Field(), err.Param())
			}
			errorMessages[err.StructField()] = append(errorMessages[err.StructField()], errorMessage)
		}
	}

	if len(errorMessages) > 0 {
		return errorMessages
	}
	return nil
}
