//
// SmartGoRest - A REST API server in Go using the Ent framework.
//
// Copyright (c) Alessio Saltarin 2026
// Licensed under ISC License - See LICENSE file for details.
//

package web

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func HomePage(c *echo.Context) error {
	viewData := map[string]any{
		"APIPortalURL": "/api",
	}

	return c.Render(http.StatusOK, "index", viewData)
}

func APIPortalPage(c *echo.Context) error {
	viewData := map[string]any{
		"SwaggerURL": "/swagger/index.html",
	}

	return c.Render(http.StatusOK, "api", viewData)
}

func DownloadOpenAPIJSON(c *echo.Context) error {
	// Offer the generated swagger spec as a downloadable OpenAPI JSON artifact.
	c.Response().Header().Set(echo.HeaderContentDisposition, `attachment; filename="smartgorest-openapi.json"`)
	return c.File("docs/swagger.json")
}
