//
// SmartGoRest - A REST API server in Go using the Ent framework.
//
// Copyright (c) Alessio Saltarin 2026
// Licensed under ISC License - See LICENSE file for details.
//

package service

import (
	"context"
	"fmt"

	"github.com/guildenstern70/smartgorest/ent"
	"github.com/guildenstern70/smartgorest/ent/person"
)

// PersonService is a struct that provides person-oriented read operations.
type PersonService struct {
	dbClient *ent.Client
}

// NewPersonService is the Constructor for PersonService
func NewPersonService(dbClient *ent.Client) *PersonService {
	return &PersonService{
		dbClient: dbClient,
	}
}

// CountPersons returns the total number of persons in the database
func (ps *PersonService) CountPersons() (int, error) {
	if err := ps.validateInput(); err != nil {
		return 0, err
	}

	return ps.dbClient.Person.Query().Count(context.Background())
}

// GetAllPersons returns all persons in the database, including their phones.
func (ps *PersonService) GetAllPersons() ([]*ent.Person, error) {
	if err := ps.validateInput(); err != nil {
		return nil, err
	}

	return ps.dbClient.Person.Query().WithPhones().All(context.Background())
}

// GetFirstNPersons returns up to n persons ordered by ascending ID, including their phones.
func (ps *PersonService) GetFirstNPersons(n int) ([]*ent.Person, error) {
	if err := ps.validateInput(); err != nil {
		return nil, err
	}
	if n < 0 {
		return nil, fmt.Errorf("n cannot be negative")
	}

	return ps.dbClient.Person.Query().WithPhones().Order(person.ByID()).Limit(n).All(context.Background())
}

// GetPersonByID returns a person by its ID, including their phones.
func (ps *PersonService) GetPersonByID(id int) (*ent.Person, error) {
	if err := ps.validateInput(); err != nil {
		return nil, err
	}
	if id <= 0 {
		return nil, fmt.Errorf("id must be positive")
	}

	return ps.dbClient.Person.Query().WithPhones().Where(person.IDEQ(id)).Only(context.Background())
}

// GetPersonByNameAndSurname returns one person matching first and last name, including their phones.
func (ps *PersonService) GetPersonByNameAndSurname(name string, surname string) (*ent.Person, error) {
	if err := ps.validateInput(); err != nil {
		return nil, err
	}
	if name == "" || surname == "" {
		return nil, fmt.Errorf("name and surname are required")
	}

	return ps.dbClient.Person.Query().
		WithPhones().
		Where(person.FirstNameEQ(name), person.LastNameEQ(surname)).
		Only(context.Background())
}

// CreatePerson creates a new person with the provided details.
func (ps *PersonService) CreatePerson(firstName, lastName, email string, age int) (*ent.Person, error) {
	if err := ps.validateInput(); err != nil {
		return nil, err
	}
	if firstName == "" || lastName == "" {
		return nil, fmt.Errorf("first name and last name are required")
	}
	if age < 0 {
		return nil, fmt.Errorf("age cannot be negative")
	}

	return ps.dbClient.Person.Create().
		SetFirstName(firstName).
		SetLastName(lastName).
		SetEmail(email).
		SetAge(age).
		Save(context.Background())
}

// DeletePerson deletes a person by their ID.
func (ps *PersonService) DeletePerson(id int) error {
	if err := ps.validateInput(); err != nil {
		return err
	}
	if id <= 0 {
		return fmt.Errorf("id must be positive")
	}

	_, err := ps.dbClient.Person.Delete().Where(person.IDEQ(id)).Exec(context.Background())
	return err
}

// validateInput ensures the service and database client are initialized.
func (ps *PersonService) validateInput() error {
	if ps == nil {
		return fmt.Errorf("person service is nil")
	}
	if ps.dbClient == nil {
		return fmt.Errorf("person service db client is nil")
	}
	return nil
}
