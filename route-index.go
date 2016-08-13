package main

import (
	"log"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// Structs

type LoginPayload struct {
	Mail     string `form:"login-mail" conform:"trim,email" validate:"required,email"`
	Password string `form:"login-password" validate:"required"`
}

// Functions

// Index renders the page first visible when
// navigating to MODULIST start page.
func (app *App) Index(c *gin.Context) {

	// Check if user is already logged in.
	// TODO: Do this.

	c.HTML(http.StatusOK, "index.html", gin.H{
		"PageTitle": "Willkommen bei MODULIST",
		"MainTitle": "Willkommen bei MODULIST",
	})
}

// Login provides all necessary functionality
// in order to log in a user.
func (app *App) Login(c *gin.Context) {

	// Check if user is already logged in.
	// TODO: Do this.

	var Payload LoginPayload

	err := c.BindWith(&Payload, binding.FormPost)
	if err != nil {
		log.Fatal("WHAT")
	}

	// Check sent content for validity.
	ErrorDesc := app.ConformAndValidate(Payload)
	if ErrorDesc != nil {

		// If payload did not pass, report errors to user.
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"PageTitle": "Willkommen bei MODULIST",
			"MainTitle": "Willkommen bei MODULIST",
			"Errors":    ErrorDesc,
		})
	}
}

func (app *App) Logout(c *gin.Context) {}
