package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iamgak/todo/models"
	"golang.org/x/oauth2"
)

// type User struct {
// 	ID        string
// 	Email     string
// 	Name      string
// 	GoogleID  string
// 	Picture   string
// 	OAuth     string
// 	CreatedAt time.Time
// }

type Config struct {
	GoogleLoginConfig oauth2.Config
}

var AppConfig Config

func (app *Application) ListTodos(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	todos, err := app.Model.Todo.ToDoListing(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	msg := []byte("User is active " + app.Username)
	err = app.Model.Todo.Publish(ctx, msg)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, todos)
}

func (app *Application) UpdateTodo(c *gin.Context) {
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
	err = app.Model.Todo.UpdateTodo(ctx, id, &todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (app *Application) SoftDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	err = app.Model.Todo.SoftDelete(ctx, app.UserID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Deleted Successfully")
}

func (app *Application) SetVisibility(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notes_id, err := strconv.Atoi(c.Param("notes_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	err = app.Model.Todo.SetVisibility(ctx, app.UserID, id, notes_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Visibility Set Successfully")
}

func (app *Application) CreateTodo(c *gin.Context) {
	var todo models.ToDo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//validation
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	todo.UserID = app.UserID
	err := app.Model.Todo.CreateTodo(ctx, &todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	msg := []byte("New to-do item added")
	err = app.Model.Todo.Publish(ctx, msg)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, todo)

}

func (app *Application) Prompt(c *gin.Context) {
	var req models.QueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var result models.ResponseBuffer
	result, err := models.Prompt(req.Parameter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
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

	var user models.UserData
	if err := json.NewDecoder(userInfo.Body).Decode(&user); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		fmt.Println("Error decoding user info:", err)
		return
	}

	// fmt.Println(user)
	user.OAuth = "google"
	user.IpAddr = c.ClientIP()
	loginToken, err := app.Model.Users.AuthUser(user)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		fmt.Println("Error authenticating user info:", err)
		return
	}

	c.SetCookie("ldata", loginToken, 3600, "/", "", false, true)
	ctx := context.Background()
	// defer cancel()
	msg := []byte("User is active " + app.Username)
	err = app.Model.Todo.Publish(ctx, msg)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, "login Successfully")

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
