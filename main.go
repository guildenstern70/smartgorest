//
// SmartGoRest - A REST API server in Go using the Ent framework.
//
// Copyright (c) Alessio Saltarin 2026
// Licensed under ISC License - See LICENSE file for details.
//

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/guildenstern70/smartgorest/ent"
	"github.com/guildenstern70/smartgorest/internal/db"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

var (
	version = "0.1"
	dbHost  = "aws-0-eu-west-1.pooler.supabase.com"
	dbPort  = "6543"
	dbUser  = "postgres.nxpjqkbxqrxaauxzygag"
	dbName  = "postgres"
)

func main() {

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	header()
	var client = connectToDatabase()
	defer func(client *ent.Client) {
		err := client.Close()
		if err != nil {
			log.Fatalf("failed closing connection to postgres: %v", err)
		}
	}(client)

	// Populate database with sample data
	dbInitializer := db.NewInitializer(context.Background(), client)
	if err := dbInitializer.Run(); err != nil {
		log.Fatalf("failed initializing database: %v", err)
	}

}

func header() {
	log.Println("====================================")
	log.Println("  Smart Go Rest - Version: ", version)
	log.Println("====================================")
}

func connectToDatabase() *ent.Client {

	log.Println("Connecting to the database...")

	// Get password from environment variable
	dbPass := os.Getenv("DB_PASSWORD")
	if dbPass == "" {
		log.Fatal("DB_PASSWORD environment variable is not set")
	}

	connectionStringTemplate := "host=%s port=%s user=%s dbname=%s password=%s"
	connectionString := fmt.Sprintf(connectionStringTemplate,
		dbHost, dbPort, dbUser, dbName, dbPass)

	client, err := ent.Open("postgres", connectionString)
	if err == nil {
		log.Println("Successfully connected to ", dbHost)
	} else {
		log.Fatalf("Failed opening connection to postgres: %v", err)
	}

	// Run the auto-migration tool.
	log.Println("Checking database schema...")
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	log.Println("Done.")

	return client
}
