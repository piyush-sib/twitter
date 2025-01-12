package main

import (
	"context"

	"go.uber.org/dig"
	"twitter/internal"
	"twitter/internal/models"
)

// GetContainer returns fully initialized container.
func getContainer(ctx context.Context, app *models.Application, flg *flags) (*dig.Container, error) {
	// Create Dependency Injection (DI) container
	container, err := internal.GetContainer(ctx, *app)
	if err != nil {
		return nil, err
	}

	for _, provide := range []struct {
		Name     string
		Resource any
		Options  []dig.ProvideOption
	}{
		{
			Name: "flags",
			Resource: func() flags {
				return *flg
			},
		},
	} {
		provideErr := container.Provide(provide.Resource, provide.Options...)
		if provideErr != nil {
			return nil, provideErr
		}
	}
	return container, nil
}
