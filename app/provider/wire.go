//go:build wireinject
// +build wireinject

package provider

import "github.com/google/wire"

func InitializeApplication() (*Application, error) {
	wire.Build(
		logProviderSet,
		provideDatabase,
		provideJson,
		provideRedis,
		provideValidator,
		provideFormDecoder,
		provideModifier,
		provideCasbin,
		provideApplication,
	)
	return &Application{}, nil
}
