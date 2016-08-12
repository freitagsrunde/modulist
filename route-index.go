package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *App) Index(c *gin.Context) {

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Title":  "Willkommen bei MODULIST",
		"Errors": &struct{}{},
	})
}

func (app *App) Login(c *gin.Context) {}

func (app *App) Logout(c *gin.Context) {}
