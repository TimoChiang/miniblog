package validator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"miniblog/domain/repository"
)

type Validator interface {
	Exec(s interface{}) map[string][]string
}

type MyValidator struct {
	*validator.Validate //combination package by anonymous
	Repo repository.ArticleRepository
}

func NewValidator(articleRepository repository.ArticleRepository) Validator {
	return &MyValidator{validator.New(), articleRepository}
}

func(v *MyValidator) Exec(s interface{}) map[string][]string{
	// custom validation
	_ = v.RegisterValidation("uniqueInDB", func(fl validator.FieldLevel) bool {
		if fl.FieldName() == "Slug" {
			return !v.Repo.SlugExists(fl.Field().String())
		}
		return false
	})
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
			case "uniqueInDB":
				errorMessage = fmt.Sprintf("%s's value: %s already in use.", err.Field(), err.Value())
			}
			errorMessages[err.StructField()] = append(errorMessages[err.StructField()], errorMessage)
		}
	}

	if len(errorMessages) > 0 {
		return errorMessages
	}
	return nil
}
