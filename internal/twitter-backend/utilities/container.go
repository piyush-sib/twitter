package utilities

import (
	"go.uber.org/dig"
)

func GetContainer(container *dig.Container) (*dig.Container, error) {
	for _, provide := range []struct {
		Name     string
		Resource any
		Options  []dig.ProvideOption
	}{
		{
			Name:     "jwt utilities",
			Resource: NewJWTUtils,
		},
		{
			Name:     "password utilities",
			Resource: NewPasswordUtils,
		},
	} {
		provideErr := container.Provide(provide.Resource, provide.Options...)
		if provideErr != nil {
			return nil, provideErr
		}
	}
	return container, nil
}
