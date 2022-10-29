package provider

import "github.com/go-playground/validator/v10"

func provideValidator() *validator.Validate {
	return validator.New()
}
