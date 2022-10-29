package provider

import "github.com/go-playground/form/v4"

func provideFormDecoder() *form.Decoder {
	return form.NewDecoder()
}
