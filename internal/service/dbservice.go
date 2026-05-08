//
// SmartGoRest - A REST API server in Go using the Ent framework.
//
// Copyright (c) Alessio Saltarin 2026
// Licensed under ISC License - See LICENSE file for details.
//

package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/guildenstern70/smartgorest/ent"
	"github.com/joho/godotenv"
)

const (
	connectTimeout = 5 * time.Second
	schemaTimeout  = 10 * time.Second
)

type DBService struct {
	dbClient *ent.Client
}

// NewDbService is the Constructor for DBService
func NewDbService() *DBService {
	dbs := &DBService{}
	dbs.dbClient = dbs.setupClient()
	return dbs
}

func (dbs *DBService) GetClient() *ent.Client {
	if dbs.dbClient == nil {
		log.Fatal("DBService client is nil")
	}
	return dbs.dbClient
}

func (dbs *DBService) CloseConnection() {
	if dbs.dbClient != nil {
		if err := dbs.dbClient.Close(); err != nil {
			log.Printf("Failed closing database client: %v", err)
		} else {
			log.Println("Database client closed successfully")
		}
	}
}

func (dbs *DBService) setupClient() *ent.Client {
	connectionString := dbs.getConnectionString()
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	drv, err := entsql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("failed opening connection to database: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()
	if err := drv.DB().PingContext(ctx); err != nil {
		_ = drv.Close()
		if strings.Contains(host, "pooler.supabase.com") && port == "6453" {
			log.Printf("hint: Supabase pooler commonly uses port 6543; current DB_PORT is %s", port)
		}
		log.Fatalf("failed connecting to %s:%s within %s: %v", host, port, connectTimeout, err)
	}

	log.Printf("Successfully connected to %s:%s", host, port)
	client := ent.NewClient(ent.Driver(drv))

	log.Println("Checking database schema...")
	ctxSchema, cancelSchema := context.WithTimeout(context.Background(), schemaTimeout)
	defer cancelSchema()
	if err := client.Schema.Create(ctxSchema); err != nil {
		if errors.Is(ctxSchema.Err(), context.DeadlineExceeded) {
			log.Fatalf("schema migration timed out after %s (possible lock contention or connectivity issue): %v", schemaTimeout, err)
		}
		log.Fatalf("failed creating schema resources: %v", err)
	}

	log.Println("Database schema is up to date.")

	return client
}

func (dbs *DBService) getConnectionString() string {

	// Load environment variables
	loadedEnv := false
	for _, envPath := range []string{".env", "../../.env"} {
		if err := godotenv.Load(envPath); err == nil {
			loadedEnv = true
			break
		}
	}
	if !loadedEnv {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Get database connection details
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		log.Fatal("DB_HOST environment variable is not set")
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		log.Fatal("DB_PORT environment variable is not set")
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		log.Fatal("DB_USER environment variable is not set")
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("DB_NAME environment variable is not set")
	}
	dbPass := os.Getenv("DB_PASSWORD")
	if dbPass == "" {
		log.Fatal("DB_PASSWORD environment variable is not set")
	}

	// Connect to Postgres once
	connectionStringTemplate := "host=%s port=%s user=%s dbname=%s password=%s sslmode=require connect_timeout=5"
	return fmt.Sprintf(connectionStringTemplate, dbHost, dbPort, dbUser, dbName, dbPass)
}
