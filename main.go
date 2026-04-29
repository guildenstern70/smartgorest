/*
 * Project SmartGoRest
 * Copyright (c) Alessio Saltarin 2026
 * Licensed under ISC Licence - See LICENSE
 */

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/guildenstern70/smartgorest/ent"

	_ "github.com/lib/pq"
)

var (
	version = "0.1"
	dbHost  = "aws-0-eu-west-1.pooler.supabase.com"
	dbPort  = "6543"
	dbUser  = "postgres.nxpjqkbxqrxaauxzygag"
	dbPass  = "Inq8bxXV5DP8sJtN"
	dbName  = "postgres"
)

const VERSION = "0.1"

func main() {
	var connectionStringTemplate = "host=%s port=%s user=%s dbname=%s password=%s"
	var connectionString = fmt.Sprintf(connectionStringTemplate,
		dbHost, dbPort, dbUser, dbName, dbPass)
	client, err := ent.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer func(client *ent.Client) {
		err := client.Close()
		if err != nil {
			log.Fatalf("failed closing connection to postgres: %v", err)
		}
	}(client)

	// Run the auto-migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
}
