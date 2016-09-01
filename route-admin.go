package main

import (
	"net/http"

	"github.com/freitagsrunde/modulist/db"
	"github.com/gin-gonic/gin"
)

// Functions

func (app *App) ListUsers(c *gin.Context) {

	// Check if user is authorized.
	User, err := app.Authorize(c.Request, db.PRIVILEGE_ADMIN)
	if err != nil {
		c.Redirect(http.StatusFound, "/modules")

		return
	}

	// Update expiration time of session.
	app.CreateSession(c, *User)

	// Fetch all users registered in database.
	var Users []db.User
	app.DB.Find(&Users)

	c.HTML(http.StatusOK, "admin-users.html", gin.H{
		"PageTitle": "Admin - Nutzerverwaltung",
		"User":      User,
		"Users":     Users,
	})
}

func (app *App) CreateUser(c *gin.Context) {}

func (app *App) DeleteUser(c *gin.Context) {}

func (app *App) SendFeedback(c *gin.Context) {

	// Check if user is authorized.
	User, err := app.Authorize(c.Request, db.PRIVILEGE_ADMIN)
	if err != nil {
		c.Redirect(http.StatusFound, "/modules")

		return
	}

	// Update expiration time of session.
	app.CreateSession(c, *User)

	// TODO: Correct admin send feedback page behaviour.

	c.HTML(http.StatusOK, "admin-send-feedback.html", gin.H{
		"PageTitle":     "Admin - Feedback versenden",
		"User":          User,
		"FeedbackMails": struct{}{},
		"Error":         "",
		"MailHeader":    "A",
		"MailFooter":    "B",
	})
}

func (app *App) UpdateMailTemplate(c *gin.Context) {}
