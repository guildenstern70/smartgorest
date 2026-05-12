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
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/guildenstern70/smartgorest/docs"
	"github.com/guildenstern70/smartgorest/ent"
	"github.com/guildenstern70/smartgorest/internal/controller/rest"
	"github.com/guildenstern70/smartgorest/internal/controller/web"
	"github.com/guildenstern70/smartgorest/internal/db"
	"github.com/guildenstern70/smartgorest/internal/service"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	_ "github.com/lib/pq"
	echoSwagger "github.com/swaggo/echo-swagger/v2"
)

// @title Smart Go Rest API
// @version 1.0
// @description A REST API server in Go using the Ent framework.

// @contact.name Alessio Saltarin
// @contact.email alessiosaltarin@gmail.com

// @BasePath  /api/v1

// @license.name ISC License
// @license.url https://opensource.org/license/isc

var (
	// version is the current application version displayed at startup.
	version = "0.1"
)

// Template HTML Renderer for Web Controller
type Template struct {
	templates *template.Template
}

func (t *Template) Render(_ *echo.Context, w io.Writer, name string, data any) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// main bootstraps configuration, database setup, sample data, and HTTP server startup.
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

// startEcho configures and runs the Echo server with graceful shutdown.
func startEcho(dbClient *ent.Client) error {

	log.Println("Starting Echo server on http://localhost:1323")

	e := echo.New()
	e.Use(middleware.RequestLogger())

	// Services
	personService := service.NewPersonService(dbClient)

	// Web Controller
	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Renderer = t
	e.GET("/", web.HomePage)

	// Open API v3.0
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// REST Controller
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

// header prints the startup banner.
func header() {
	log.Println("====================================")
	log.Println("  Smart Go Rest - Version: ", version)
	log.Println("====================================")
}
