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

	return c.Render(http.StatusOK, "index", nil)
}
