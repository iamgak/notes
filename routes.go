package main

import "github.com/gin-gonic/gin"

func (app *Application) InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/", app.listTodos)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// r.GET("/", app.getNotesListing)

	// Apply middleware to routes
	// r.POST("/secure", authMiddleware(), func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"Secure": true,
	// 	})
	// 	// This route is now protected and requires authentication
	// })
	return r
}
