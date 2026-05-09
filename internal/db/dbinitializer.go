//
// SmartGoRest - A REST API server in Go using the Ent framework.
//
// Copyright (c) Alessio Saltarin 2026
// Licensed under ISC License - See LICENSE file for details.
//

package db

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/guildenstern70/smartgorest/ent"
	"github.com/guildenstern70/smartgorest/ent/phone"
	"github.com/guildenstern70/smartgorest/internal/service"
)

// Initializer is a struct that seeds the database with initial sample data.
type Initializer struct {
	ctx    context.Context
	client *ent.Client
}

// NewInitializer creates a database initializer with a safe default context.
func NewInitializer(ctx context.Context, client *ent.Client) *Initializer {
	if ctx == nil {
		ctx = context.Background()
	}

	return &Initializer{
		ctx:    ctx,
		client: client,
	}
}

// Run populates the database with sample data when it is empty.
func (dbinit *Initializer) Run() error {
	if dbinit == nil {
		return fmt.Errorf("db initializer is nil")
	}
	if dbinit.client == nil {
		return fmt.Errorf("db initializer client is nil")
	}
	if dbinit.ctx == nil {
		dbinit.ctx = context.Background()
	}

	log.Println("Initializing the database...")

	// Check if database already has persons
	personService := service.NewPersonService(dbinit.client)
	count, err := personService.CountPersons()
	if err != nil {
		return fmt.Errorf("failed to count persons: %w", err)
	}

	if count > 0 {
		log.Printf("Database already contains %d person(s). Skipping initialization.", count)
		return nil
	}

	phones := dbinit.addPhones()
	dbinit.addPersons(phones)
	log.Println("Database initialized.")

	return nil
}

// addPhones creates a pool of sample phone records.
func (dbinit *Initializer) addPhones() []*ent.Phone {
	log.Println("Adding phones...")

	prefixes := []string{"+1", "+39", "+44", "+49", "+33", "+34"}
	kinds := []phone.Kind{phone.KindCellular, phone.KindHome, phone.KindWork}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	phones := make([]*ent.Phone, 0, 20)
	for i := 0; i < 20; i++ {
		prefix := prefixes[r.Intn(len(prefixes))]
		number := fmt.Sprintf("%03d%07d", r.Intn(1000), r.Intn(10000000))
		kind := kinds[r.Intn(len(kinds))]

		p, err := dbinit.createPhone(prefix, number, kind)
		if err != nil {
			log.Printf("failed to create phone: %v", err)
			return phones
		}
		phones = append(phones, p)
	}

	log.Println("... done.")
	return phones
}

// addPersons creates sample people and assigns phones from the pool.
func (dbinit *Initializer) addPersons(phones []*ent.Phone) {

	log.Println("Adding persons...")

	firstNames := []string{"Liam", "Olivia", "Noah", "Emma", "Lucas", "Mia", "Ethan", "Sofia", "Ava", "Elijah", "Isabella", "Mason"}
	lastNames := []string{"Brown", "Wilson", "Taylor", "Anderson", "Thomas", "Moore", "Jackson", "Martin", "Lee", "Walker", "Harris", "Clark"}
	domains := []string{"example.com", "mail.com", "test.org", "acme.dev"}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate 15 Persons with random combinations.
	for i := 0; i < 15; i++ {
		firstName := firstNames[r.Intn(len(firstNames))]
		lastName := lastNames[r.Intn(len(lastNames))]
		age := 18 + r.Intn(48)
		email := fmt.Sprintf(
			"%s.%s%d@%s",
			strings.ToLower(firstName),
			strings.ToLower(lastName),
			r.Intn(1000),
			domains[r.Intn(len(domains))],
		)

		// Assign 1 or 2 random phones to this person, consuming them from the pool.
		count := 1 + r.Intn(2)
		assigned := make([]*ent.Phone, 0, count)
		for j := 0; j < count && len(phones) > 0; j++ {
			idx := r.Intn(len(phones))
			assigned = append(assigned, phones[idx])
			phones = append(phones[:idx], phones[idx+1:]...)
		}

		if err := dbinit.createPerson(firstName, lastName, email, age, assigned); err != nil {
			log.Printf("failed to create person: %v", err)
			return
		}
	}
	log.Println("... done.")

}

// createPhone persists a single phone entity.
func (dbinit *Initializer) createPhone(prefix string, number string, kind phone.Kind) (*ent.Phone, error) {
	p, err := dbinit.client.Phone.
		Create().
		SetPrefix(prefix).
		SetNumber(number).
		SetKind(kind).
		Save(dbinit.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating phone: %w", err)
	}

	return p, nil
}

// createPerson persists a single person with optional phone relations.
func (dbinit *Initializer) createPerson(
	firstName string,
	lastName string,
	email string,
	age int,
	phones []*ent.Phone) error {
	_, err := dbinit.client.Person.
		Create().
		SetAge(age).
		SetFirstName(firstName).
		SetLastName(lastName).
		SetEmail(email).
		AddPhones(phones...).
		Save(dbinit.ctx)
	if err != nil {
		return fmt.Errorf("failed creating user: %w", err)
	}

	return nil
}
