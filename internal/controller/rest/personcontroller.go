//
// SmartGoRest - A REST API server in Go using the Ent framework.
//
// Copyright (c) Alessio Saltarin 2026
// Licensed under ISC License - See LICENSE file for details.
//

package rest

import (
	"net/http"
	"strconv"

	"github.com/guildenstern70/smartgorest/internal/service"
	"github.com/labstack/echo/v5"
)

// PersonController is a struct that exposes person endpoints over HTTP.
type PersonController struct {
	personService *service.PersonService
}

// NewPersonController creates a controller backed by a person service.
func NewPersonController(personService *service.PersonService) *PersonController {
	return &PersonController{
		personService: personService,
	}
}

// GetPersons handles requests that fetch all persons.
func (pc *PersonController) GetPersons(c *echo.Context) error {

	persons, err := pc.personService.GetAllPersons()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve persons"})
	}
	return c.JSON(http.StatusOK, persons)
}

// GetPerson handles requests that fetch one person by ID.
func (pc *PersonController) GetPerson(c *echo.Context) error {

	var err error

	id := c.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid ID format"})
	}
	person, err := pc.personService.GetPersonByID(intId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve person"})
	}
	return c.JSON(http.StatusOK, person)
}
