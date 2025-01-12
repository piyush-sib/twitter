package internal

import (
	"context"
	"log"

	"go.uber.org/dig"
	"twitter/internal/closer"
	"twitter/internal/infrastructure"
	"twitter/internal/models"
	"twitter/internal/structuredlogger"
	tb "twitter/internal/twitter-backend"
)

// GetContainer initialises the Dependency Injection Container.
// It is a basic implementation of uber/dig initialisation, usage of functions like `initResources` is
// just here to reduce cyclomatic complexity in case we have a significant amount of dependencies to initialise.
func GetContainer(ctx context.Context, app models.Application) (*dig.Container, error) {
	log.Println("Initialising Dependency Injection Container")

	c, err := getNewDigContainer(ctx, app)
	if err != nil {
		return nil, err
	}

	// initialise twitter container.
	c, err = tb.GetContainer(c)
	if err != nil {
		return nil, err
	}
	// initialise infrastructure container.
	c, err = infrastructure.GetContainer(c)
	if err != nil {
		return nil, err
	}

	// initialise structured-logger container.
	c, err = structuredlogger.GetContainer(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func getNewDigContainer(ctx context.Context, app models.Application) (*dig.Container, error) {
	c := dig.New()

	err := c.Provide(func() context.Context {
		return ctx
	})
	if err != nil {
		return nil, err
	}

	_ = c.Provide(func() models.Application {
		return app
	})

	err = c.Provide(closer.NewCloser)
	if err != nil {
		return nil, err
	}
	return c, err
}
