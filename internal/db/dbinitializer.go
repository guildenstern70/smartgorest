package db

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/guildenstern70/smartgorest/ent"
)

type Initializer struct {
	ctx    context.Context
	client *ent.Client
}

func NewInitializer(ctx context.Context, client *ent.Client) *Initializer {
	if ctx == nil {
		ctx = context.Background()
	}

	return &Initializer{
		ctx:    ctx,
		client: client,
	}
}

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
	if err := dbinit.createPerson("John", "Doe", "john.doe@example.com", 30); err != nil {
		return err
	}
	if err := dbinit.createPerson("Jane", "Smith", "jane.smith@example.com", 25); err != nil {
		return err
	}

	firstNames := []string{"Liam", "Olivia", "Noah", "Emma", "Lucas", "Mia", "Ethan", "Sofia", "Ava", "Elijah", "Isabella", "Mason"}
	lastNames := []string{"Brown", "Wilson", "Taylor", "Anderson", "Thomas", "Moore", "Jackson", "Martin", "Lee", "Walker", "Harris", "Clark"}
	domains := []string{"example.com", "mail.com", "test.org", "acme.dev"}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate 12 additional records with random combinations.
	for i := 0; i < 12; i++ {
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

		if err := dbinit.createPerson(firstName, lastName, email, age); err != nil {
			return err
		}
	}
	log.Println("Done.")

	return nil
}

func (dbinit *Initializer) createPerson(
	firstName string,
	lastName string,
	email string,
	age int) error {
	_, err := dbinit.client.Person.
		Create().
		SetAge(age).
		SetFirstName(firstName).
		SetLastName(lastName).
		SetEmail(email).
		Save(dbinit.ctx)
	if err != nil {
		return fmt.Errorf("failed creating user: %w", err)
	}

	return nil
}
