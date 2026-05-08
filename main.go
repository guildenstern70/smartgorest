//
// SmartGoRest - A REST API server in Go using the Ent framework.
//
// Copyright (c) Alessio Saltarin 2026
// Licensed under ISC License - See LICENSE file for details.
//

package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/guildenstern70/smartgorest/ent"
	"github.com/guildenstern70/smartgorest/internal/controller/rest"
	"github.com/guildenstern70/smartgorest/internal/db"
	"github.com/guildenstern70/smartgorest/internal/service"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	_ "github.com/lib/pq"
)

var (
	version = "0.1"
)

func main() {

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	header()
	dbService := service.NewDbService()
	defer dbService.CloseConnection()
	client := dbService.GetClient()

	// Populate database with sample data
	dbInitializer := db.NewInitializer(context.Background(), client)
	if err := dbInitializer.Run(); err != nil {
		log.Fatalf("failed initializing database: %v", err)
	}

	// Start the Echo server
	if err := startEcho(client); err != nil {
		log.Fatalf("echo server terminated with error: %v", err)
	}

}

func startEcho(dbClient *ent.Client) error {

	log.Println("Starting Echo server on http://localhost:1323")

	e := echo.New()
	e.Use(middleware.RequestLogger())

	// Services
	personService := service.NewPersonService(dbClient)

	// Controller
	personController := rest.NewPersonController(personService)
	persons := e.Group("/persons")
	persons.GET("", personController.GetPersons)
	persons.GET("/:id", personController.GetPerson)

	server := &http.Server{
		Addr:    ":1323",
		Handler: e,
	}

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.ListenAndServe()
	}()

	signalContext, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	case <-signalContext.Done():
		log.Println("Termination signal received, shutting down Echo server...")
		shutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownContext); err != nil {
			return err
		}

		err := <-serverErr
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}
}

func header() {
	log.Println("====================================")
	log.Println("  Smart Go Rest - Version: ", version)
	log.Println("====================================")
}
