//
// SmartGoRest - A REST API server in Go using the Ent framework.
//
// Copyright (c) Alessio Saltarin 2026
// Licensed under ISC License - See LICENSE file for details.
//

package dto

import "github.com/guildenstern70/smartgorest/ent"

// PhoneDTO is the Data Transfer Object for the Phone entity.
type PhoneDTO struct {
	ID     int    `json:"id,omitempty"`
	Prefix string `json:"prefix,omitempty"`
	Number string `json:"number,omitempty"`
	Kind   string `json:"kind,omitempty"`
}

// PersonDTO is the Data Transfer Object for the Person entity.
// The "phones" field maps the Ent edges (phones relation) into a flat, serialisable list.
type PersonDTO struct {
	ID        int         `json:"id,omitempty"`
	FirstName string      `json:"first_name,omitempty"`
	LastName  string      `json:"last_name,omitempty"`
	Age       int         `json:"age,omitempty"`
	Email     string      `json:"email,omitempty"`
	Phones    []*PhoneDTO `json:"phones,omitempty"`
}

// PhoneDTOFromEnt maps an ent.Phone entity to a PhoneDTO.
func PhoneDTOFromEnt(p *ent.Phone) *PhoneDTO {
	if p == nil {
		return nil
	}
	return &PhoneDTO{
		ID:     p.ID,
		Prefix: p.Prefix,
		Number: p.Number,
		Kind:   string(p.Kind),
	}
}

// PersonDTOFromEnt maps an ent.Person entity (with eager-loaded phones) to a PersonDTO.
func PersonDTOFromEnt(p *ent.Person) *PersonDTO {
	if p == nil {
		return nil
	}

	phones := make([]*PhoneDTO, 0, len(p.Edges.Phones))
	for _, ph := range p.Edges.Phones {
		phones = append(phones, PhoneDTOFromEnt(ph))
	}

	return &PersonDTO{
		ID:        p.ID,
		FirstName: p.FirstName,
		LastName:  p.LastName,
		Age:       p.Age,
		Email:     p.Email,
		Phones:    phones,
	}
}

// PersonDTOsFromEnt maps a slice of ent.Person entities to a slice of PersonDTO pointers.
func PersonDTOsFromEnt(persons []*ent.Person) []*PersonDTO {
	dtos := make([]*PersonDTO, 0, len(persons))
	for _, p := range persons {
		dtos = append(dtos, PersonDTOFromEnt(p))
	}
	return dtos
}
