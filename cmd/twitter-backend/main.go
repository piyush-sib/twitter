package main

import (
	"context"

	"twitter/internal/closer"
	"twitter/internal/models"
	"twitter/internal/twitter-backend/service"
)

func main() {
	// Parse application flags
	flg, err := getFlags()
	if err != nil {
		panic(err)
	}

	// Define app information
	app := &models.Application{
		Env:  flg.Environment,
		Name: flg.AppName,
	}

	ctx := context.Background()
	// Initialize DI container
	container, err := getContainer(ctx, app, flg)
	if err != nil {
		panic(err)
	}

	// Close resources
	defer func() {
		errCloser := container.Invoke(func(closer *closer.Closer) {
			closer.Close()
		})
		if errCloser != nil {
			panic(errCloser)
		}
	}()
	// Run HTTP server
	err = service.RunHTTPServer(container)
	if err != nil {
		panic(err)
	}

}
