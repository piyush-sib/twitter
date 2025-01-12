package service

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
			Name:     "user handler",
			Resource: NewUserHandler,
		},
		{
			Name:     "feed handler",
			Resource: NewFeedHandler,
		},
		{
			Name:     "tweets handler",
			Resource: NewTweetsHandler,
		},
	} {
		provideErr := container.Provide(provide.Resource, provide.Options...)
		if provideErr != nil {
			return nil, provideErr
		}
	}
	return container, nil
}
