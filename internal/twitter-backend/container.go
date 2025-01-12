package mapper

import (
	"go.uber.org/dig"
	"twitter/internal/twitter-backend/middlewares"
	"twitter/internal/twitter-backend/repository"
	"twitter/internal/twitter-backend/service"
	"twitter/internal/twitter-backend/utilities"
)

func GetContainer(container *dig.Container) (_ *dig.Container, err error) {
	// dependency injection container for repository.
	container, err = repository.GetContainer(container)
	if err != nil {
		return nil, err
	}

	// dependency injection container for utilities.
	container, err = utilities.GetContainer(container)
	if err != nil {
		return nil, err
	}
	// dependency injection container for middlewares.
	container, err = middleware.GetContainer(container)
	if err != nil {
		return nil, err
	}

	// dependency injection container for service.
	container, err = service.GetContainer(container)
	if err != nil {
		return nil, err
	}

	return container, nil
}
