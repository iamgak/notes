package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iamgak/todo/models"
	"golang.org/x/oauth2"
)

type Config struct {
	GoogleLoginConfig oauth2.Config
}

var AppConfig Config

func (app *Application) listTodos(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	todos, err := app.Model.ToDoListing(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todos)
}

func (app *Application) updateTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var todo models.ToDo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	err = app.Model.UpdateTodo(ctx, id, &todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todo)

}

func (app *Application) createTodo(c *gin.Context) {
	var todo models.ToDo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	err := app.Model.CreateTodo(ctx, &todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, todo)

}

func (app *Application) GoogleCallback(c *gin.Context) {
	// state := c.Query("state")
	code := c.Query("code")
	googleOauthConfig, err := InitGoogleOAuth()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	token, err := googleOauthConfig.Exchange(c.Request.Context(), code)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	client := googleOauthConfig.Client(c.Request.Context(), token)
	userInfo, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		fmt.Println("Error fetching user info:", err)
		return
	}
	defer userInfo.Body.Close()

	// Read the response body and print the user information
	body, err := ioutil.ReadAll(userInfo.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		fmt.Println("Error reading user info response:", err)
		return
	}
	c.String(http.StatusOK, "User information: %v", string(body))

}

func (app *Application) GoogleLogin(c *gin.Context) {
	googleOauthConfig, err := InitGoogleOAuth()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	url := googleOauthConfig.AuthCodeURL("state-token")
	c.Redirect(http.StatusTemporaryRedirect, url)
}
