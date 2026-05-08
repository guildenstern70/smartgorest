//
// SmartGoRest - A REST API server in Go using the Ent framework.
//
// Copyright (c) Alessio Saltarin 2026
// Licensed under ISC License - See LICENSE file for details.
//

package service

import (
	"os"
	"testing"

	"github.com/guildenstern70/smartgorest/ent"
	_ "github.com/lib/pq"
)

var personService *PersonService

// TestMain sets up the test database once for all tests
func TestMain(m *testing.M) {

	// Setup DB
	var dbService = NewDbService()
	var dbClient = dbService.GetClient()
	personService = NewPersonService(dbClient)

	// Run tests
	code := m.Run()

	// Cleanup
	_ = dbClient.Close()
	os.Exit(code)
}

func TestGetAllPersons(t *testing.T) {

	persons, err := personService.GetAllPersons()
	if err != nil {
		t.Fatalf("GetAllPersons returned error: %v", err)
	}
	if len(persons) == 0 {
		t.Fatalf("expected more than 1 person, got %d", len(persons))
	}
}

func TestGetFirstNPersons(t *testing.T) {

	persons, err := personService.GetFirstNPersons(2)
	if err != nil {
		t.Fatalf("GetFirstNPersons returned error: %v", err)
	}
	if len(persons) != 2 {
		t.Fatalf("expected 2 persons, got %d", len(persons))
	}
	if persons[0].ID >= persons[1].ID {
		t.Fatalf("expected persons ordered by ID ascending, got ids %d and %d", persons[0].ID, persons[1].ID)
	}
}

func TestGetFirstNPersonsRejectsNegativeN(t *testing.T) {

	_, err := personService.GetFirstNPersons(-1)
	if err == nil {
		t.Fatal("expected error for negative n")
	}
}

func TestGetPersonByID(t *testing.T) {

	persons, err := personService.GetFirstNPersons(1)
	if err != nil {
		t.Fatalf("GetAllPersons returned error: %v", err)
	}

	personByID, err := personService.GetPersonByID(persons[0].ID)
	if err != nil {
		t.Fatalf("GetPersonByID returned error: %v", err)
	}
	if personByID.ID != persons[0].ID {
		t.Fatalf("expected ID %d, got %d", persons[0].ID, personByID.ID)
	}
}

func TestGetPersonByIDNotFound(t *testing.T) {

	_, err := personService.GetPersonByID(99999)
	if err == nil {
		t.Fatal("expected not found error")
	}
	if !ent.IsNotFound(err) {
		t.Fatalf("expected ent not found error, got: %v", err)
	}
}

func TestGetPersonByIDRejectsNonPositiveID(t *testing.T) {

	_, err := personService.GetPersonByID(0)
	if err == nil {
		t.Fatal("expected error for non-positive id")
	}
}

func TestGetPersonByNameAndSurname(t *testing.T) {

	persons, err := personService.GetFirstNPersons(1)

	personByName, err := personService.GetPersonByNameAndSurname(persons[0].FirstName, persons[0].LastName)
	if err != nil {
		t.Fatalf("GetPersonByNameAndSurname returned error: %v", err)
	}
	if personByName.FirstName != persons[0].FirstName || personByName.LastName != persons[0].LastName {
		t.Fatalf("unexpected person returned: %s %s", personByName.FirstName, personByName.LastName)
	}
}

func TestGetPersonByNameAndSurnameNotFound(t *testing.T) {

	_, err := personService.GetPersonByNameAndSurname("No", "Body")
	if err == nil {
		t.Fatal("expected not found error")
	}
	if !ent.IsNotFound(err) {
		t.Fatalf("expected ent not found error, got: %v", err)
	}
}

func TestGetPersonByNameAndSurnameRejectsEmptyInput(t *testing.T) {

	_, err := personService.GetPersonByNameAndSurname("", "Smith")
	if err == nil {
		t.Fatal("expected error when name or surname are empty")
	}
}

func TestCountPersons(t *testing.T) {

	count, err := personService.CountPersons()
	if err != nil {
		t.Fatalf("CountPersons returned error: %v", err)
	}
	if count == 0 {
		t.Fatalf("expected 3 persons, got %d", count)
	}
}
