//
// SmartGoRest - A REST API server in Go using the Ent framework.
//
// Copyright (c) Alessio Saltarin 2026
// Licensed under ISC License - See LICENSE file for details.
//

package service

import (
	"testing"

	"github.com/guildenstern70/smartgorest/ent"
	"github.com/guildenstern70/smartgorest/ent/enttest"
	_ "github.com/mattn/go-sqlite3"
)

func newTestPersonService(t *testing.T) *PersonService {
	t.Helper()

	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	t.Cleanup(func() {
		_ = client.Close()
	})

	_, err := client.Person.Create().
		SetFirstName("John").
		SetLastName("Doe").
		SetAge(30).
		SetEmail("john.doe@example.com").
		Save(t.Context())
	if err != nil {
		t.Fatalf("seed person John Doe: %v", err)
	}

	_, err = client.Person.Create().
		SetFirstName("Jane").
		SetLastName("Smith").
		SetAge(25).
		SetEmail("jane.smith@example.com").
		Save(t.Context())
	if err != nil {
		t.Fatalf("seed person Jane Smith: %v", err)
	}

	_, err = client.Person.Create().
		SetFirstName("Luca").
		SetLastName("Bianchi").
		SetAge(41).
		SetEmail("luca.bianchi@example.com").
		Save(t.Context())
	if err != nil {
		t.Fatalf("seed person Luca Bianchi: %v", err)
	}

	return NewPersonService(client)
}

func TestGetAllPersons(t *testing.T) {
	service := newTestPersonService(t)

	persons, err := service.GetAllPersons()
	if err != nil {
		t.Fatalf("GetAllPersons returned error: %v", err)
	}
	if len(persons) != 3 {
		t.Fatalf("expected 3 persons, got %d", len(persons))
	}
}

func TestGetFirstNPersons(t *testing.T) {
	service := newTestPersonService(t)

	persons, err := service.GetFirstNPersons(2)
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
	service := newTestPersonService(t)

	_, err := service.GetFirstNPersons(-1)
	if err == nil {
		t.Fatal("expected error for negative n")
	}
}

func TestGetPersonByID(t *testing.T) {
	service := newTestPersonService(t)

	persons, err := service.GetAllPersons()
	if err != nil {
		t.Fatalf("GetAllPersons returned error: %v", err)
	}

	personByID, err := service.GetPersonByID(persons[0].ID)
	if err != nil {
		t.Fatalf("GetPersonByID returned error: %v", err)
	}
	if personByID.ID != persons[0].ID {
		t.Fatalf("expected ID %d, got %d", persons[0].ID, personByID.ID)
	}
}

func TestGetPersonByIDNotFound(t *testing.T) {
	service := newTestPersonService(t)

	_, err := service.GetPersonByID(99999)
	if err == nil {
		t.Fatal("expected not found error")
	}
	if !ent.IsNotFound(err) {
		t.Fatalf("expected ent not found error, got: %v", err)
	}
}

func TestGetPersonByIDRejectsNonPositiveID(t *testing.T) {
	service := newTestPersonService(t)

	_, err := service.GetPersonByID(0)
	if err == nil {
		t.Fatal("expected error for non-positive id")
	}
}

func TestGetPersonByNameAndSurname(t *testing.T) {
	service := newTestPersonService(t)

	personByName, err := service.GetPersonByNameAndSurname("Jane", "Smith")
	if err != nil {
		t.Fatalf("GetPersonByNameAndSurname returned error: %v", err)
	}
	if personByName.FirstName != "Jane" || personByName.LastName != "Smith" {
		t.Fatalf("unexpected person returned: %s %s", personByName.FirstName, personByName.LastName)
	}
}

func TestGetPersonByNameAndSurnameNotFound(t *testing.T) {
	service := newTestPersonService(t)

	_, err := service.GetPersonByNameAndSurname("No", "Body")
	if err == nil {
		t.Fatal("expected not found error")
	}
	if !ent.IsNotFound(err) {
		t.Fatalf("expected ent not found error, got: %v", err)
	}
}

func TestGetPersonByNameAndSurnameRejectsEmptyInput(t *testing.T) {
	service := newTestPersonService(t)

	_, err := service.GetPersonByNameAndSurname("", "Smith")
	if err == nil {
		t.Fatal("expected error when name or surname are empty")
	}
}

func TestPersonServiceNilReceiver(t *testing.T) {
	var service *PersonService

	_, err := service.GetAllPersons()
	if err == nil {
		t.Fatal("expected error for nil service receiver")
	}
}
