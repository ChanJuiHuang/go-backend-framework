package provider

import (
	"github.com/go-playground/mold/v4"
	"github.com/go-playground/mold/v4/modifiers"
)

func provideModifier() *mold.Transformer {
	return modifiers.New()
}
