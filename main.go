//
// SmartGoRest - A REST API server in Go using the Ent framework.
//
// Copyright (c) Alessio Saltarin 2026
// Licensed under ISC License - See LICENSE file for details.
//

package main

import (
	"context"
	"log"

	"github.com/guildenstern70/smartgorest/internal/db"
	"github.com/guildenstern70/smartgorest/internal/service"
	"github.com/joho/godotenv"

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
	var dbService = service.NewDbService()
	var client = dbService.GetClient()

	// Populate database with sample data
	dbInitializer := db.NewInitializer(context.Background(), client)
	if err := dbInitializer.Run(); err != nil {
		log.Fatalf("failed initializing database: %v", err)
	}

	dbService.CloseConnection()

}

func header() {
	log.Println("====================================")
	log.Println("  Smart Go Rest - Version: ", version)
	log.Println("====================================")
}
