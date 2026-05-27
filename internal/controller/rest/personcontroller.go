//
// SmartGoRest - A REST API server in Go using the Ent framework.
//
// Copyright (c) Alessio Saltarin 2026
// Licensed under ISC License - See LICENSE file for details.
//

package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/guildenstern70/smartgorest/internal/service"
	"github.com/labstack/echo/v5"
)

// CreatePersonRequest represents the request body for creating a person.
type CreatePersonRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email"`
	Age       int    `json:"age" binding:"required"`
}

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

// GetPersons godoc
// @Summary      Get all persons
// @Description  Get a list of all persons in the system.
// @Tags         persons
// @Accept       json
// @Produce      json
// @Success      200  {array}   dto.PersonDTO
// @Failure      500  {object}  map[string]string
// @Router       /persons [get]
func (pc *PersonController) GetPersons(c *echo.Context) error {

	persons, err := pc.personService.GetAllPersons()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve persons"})
	}
	return c.JSON(http.StatusOK, persons)
}

// GetPerson godoc
// @Summary      Get a person by ID
// @Description  Get details of a specific person by their ID.
// @Tags         persons
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Person ID"
// @Success      200  {object}  dto.PersonDTO
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /persons/{id} [get]
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

// CreatePerson godoc
// @Summary      Create a new person
// @Description  Create a new person with the provided details.
// @Tags         persons
// @Accept       json
// @Produce      json
// @Param        person  body      CreatePersonRequest  true  "Person data"
// @Success      201     {object}  dto.PersonDTO
// @Failure      400     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /persons [post]
func (pc *PersonController) CreatePerson(c *echo.Context) error {
	var req CreatePersonRequest

	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	person, err := pc.personService.CreatePerson(req.FirstName, req.LastName, req.Email, req.Age)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create person: " + err.Error()})
	}

	return c.JSON(http.StatusCreated, person)
}

// DeletePerson godoc
// @Summary      Delete a person
// @Description  Delete a person by their ID.
// @Tags         persons
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Person ID"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /persons/{id} [delete]
func (pc *PersonController) DeletePerson(c *echo.Context) error {
	id := c.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}

	err = pc.personService.DeletePerson(intId)
	if err != nil {
		// Check if it's a NotFoundError
		if err.Error() == "person not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Person not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete person: " + err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
