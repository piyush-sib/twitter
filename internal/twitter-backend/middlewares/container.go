package middleware

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
			Name:     "auth middlewares",
			Resource: NewAuthMiddlewares,
		},
	} {
		provideErr := container.Provide(provide.Resource, provide.Options...)
		if provideErr != nil {
			return nil, provideErr
		}
	}
	return container, nil
}
