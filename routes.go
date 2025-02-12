package main

import (
	"github.com/gin-gonic/gin"
)

func (app *Application) InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/user/login", secureHeaders(), func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Page: login",
		})
	})

	authorise := r.Group("/")
	authorise.Use(app.LoginMiddleware(), secureHeaders())
	{

		// read API
		authorise.GET("/", app.ListTodos)
		// authorise.GET("/",app.listTodos)

		// write API
		authorise.POST("/", app.CreateTodo)
		authorise.POST("/:id/update/", app.UpdateTodo)
		authorise.POST("/:id/visibilty/:object_id/", app.SetVisibility)
		authorise.POST("/:id/delete/", app.SoftDelete)
	}

	// r.GET("/", app.getNotesListing)

	r.GET("/google_login", secureHeaders(), app.GoogleLogin)
	r.GET("/callback", secureHeaders(), app.GoogleCallback)
	return r
}
