package service

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/dig"
	"twitter/internal/infrastructure/mysql"
	"twitter/internal/twitter-backend/middlewares"
)

type newHTTPServerParams struct {
	dig.In

	AppName     string `name:"appname"`
	Environment string `name:"environment"`

	Context       context.Context
	MySQLSeeder   *mysql.MySQLSeederHandler
	UserHandler   *UserHandler
	FeedsHandler  *FeedHandler
	TweetsHandler *TweetsHandler
	AuthHandler   *middleware.AuthMiddlewares
}

func RunHTTPServer(container *dig.Container) error {
	err := container.Invoke(func(params newHTTPServerParams) error {
		return runHTTP(params)
	})
	if err != nil {
		return err
	}
	return nil
}

func runHTTP(params newHTTPServerParams) error {
	// Create a new router
	router := mux.NewRouter()

	// Register routes
	RegisterRoutes(router, params.UserHandler, params.AuthHandler, params.FeedsHandler, params.TweetsHandler)

	// Create an HTTP server
	server := &http.Server{
		Addr:    ":8080", // Change to desired port
		Handler: router,
	}

	// Channel to listen for OS signals for shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Start the server in a goroutine
	go func() {
		log.Printf("Starting server on %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", server.Addr, err)
		}
	}()

	params.MySQLSeeder.MigrateSchema()

	// Wait for a signal to shut down
	<-stop
	log.Println("Shutting down server...")

	// Create a context for server shutdown
	ctx, cancel := context.WithTimeout(params.Context, 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
	return nil
}
