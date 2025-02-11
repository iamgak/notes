package main

import (
	"github.com/gin-gonic/gin"
)

func (app *Application) InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// read API
	r.GET("/", secureHeaders(), app.LoginMiddleware(), app.listTodos)
	// r.GET("/", secureHeaders(), app.listTodos)
	r.GET("/user/login", secureHeaders(), func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Page: login",
		})
	})

	// write API
	r.POST("/", secureHeaders(), app.LoginMiddleware(), app.createTodo)
	r.POST("/update", secureHeaders(), app.LoginMiddleware(), app.updateTodo)

	// r.GET("/", app.getNotesListing)

	r.GET("/google_login", secureHeaders(), app.GoogleLogin)
	r.GET("/callback", secureHeaders(), app.GoogleCallback)
	return r
}
