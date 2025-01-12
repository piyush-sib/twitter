package main

import (
	"context"

	"twitter/internal/closer"
	"twitter/internal/models"
	"twitter/internal/twitter-backend/service"
)

func main() {
	// Parse application flags
	flg := getFlags()

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

//
//func gracefulShutdown(closables []Closable) {
//	// Catch OS signals for termination
//	sigChan := make(chan os.Signal, 1)
//	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
//
//	<-sigChan // Wait for a signal
//	for _, c := range closables {
//		if err := c.Close(); err != nil {
//			// Log the error, ideally you should use your logger here
//			log.Printf("Failed to close resource: %v", err)
//		}
//	}
//
//	log.Println("Application shutdown gracefully.")
//}
