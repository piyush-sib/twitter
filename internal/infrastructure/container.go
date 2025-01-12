package infrastructure

import (
	"go.uber.org/dig"
	"log"
	"twitter/internal/infrastructure/mysql"
)

func GetContainer(container *dig.Container) (_ *dig.Container, err error) {
	log.Println("Initialising Infrastructure Dependency Injection Container")
	// dependency injection container for mongodb.
	container, err = mysql.GetContainer(container)
	if err != nil {
		return nil, err
	}

	// Provide resources
	for _, provide := range []struct {
		Name     string
		Resource any
		Options  []dig.ProvideOption
	}{} {
		provideErr := container.Provide(provide.Resource, provide.Options...)
		if provideErr != nil {
			return nil, err
		}
	}
	return container, nil
}
