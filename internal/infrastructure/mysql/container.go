package mysql

import (
	"go.uber.org/dig"
	"log"
)

func GetContainer(container *dig.Container) (*dig.Container, error) {
	log.Println("Initialising MySQL Dependency Injection Container")
	for _, provide := range []struct {
		Name     string
		Resource any
		Options  []dig.ProvideOption
	}{
		{
			Name:     "mysql client handler",
			Resource: newMySQLClientHandler,
		},
		{
			Name:     "mysql seeder handler",
			Resource: NewMySQLSeederHandler,
		},
	} {
		provideErr := container.Provide(provide.Resource, provide.Options...)
		if provideErr != nil {
			return nil, provideErr
		}
	}
	return container, nil
}
