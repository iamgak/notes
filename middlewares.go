package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/iamgak/todo/models"
	"github.com/joho/godotenv"
)

func secureHeaders() gin.HandlerFunc {
	return (func(c *gin.Context) {
		c.Header("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		c.Header("Referrer-Policy", "origin-when-cross-origin")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "deny")
		c.Header("X-XSS-Protection", "0")
		c.Next()
	})
}

func (app *Application) recoverPanic() gin.HandlerFunc {
	// return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	return func(c *gin.Context) {

		// Create a deferred function (which will always be run in the event of a panic
		// as Go unwinds the stack).
		defer func() {
			// Use the builtin recover function to check if there has been a panic or
			// not.
			if err := recover(); err != nil {
				// If there was a panic, set a "Connection: close" header on the
				// response. This acts as a trigger to make Go's HTTP server
				// automatically close the current connection after a response has been
				// sent.
				c.Header("Connection", "close")
				// The value returned by recover() has the type interface{}, so we use
				// fmt.Errorf() to normalize it into an error and call our
				// serverErrorResponse() helper. In turn, this will log the error using
				// our custom Logger type at the ERROR level and send the client a 500
				// Internal Server Error response.
				// app.ErrorMessage(w, 404, fmt.Errorf("%s", err))

				c.JSON(500, gin.H{"error": "Internal Server Error"})
				// log some message
				c.Abort()
			}
		}()

		c.Next()
		// })
	}
}

func (app *Application) LoginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("ldata")

		if err != nil || cookie == "" || len(cookie) < 40 {
			// deleteCookie(c)
			c.Redirect(http.StatusFound, "/user/login/")
			c.Abort()
			return
		}

		err = godotenv.Load()
		if err != nil {
			fmt.Println("Error fetching env:", err)
			// c.Redirect(http.StatusInternalServerError, "/user/login/")
			// c.AbortWithStatusJSON(http.StatusInternalServerError,
			c.JSON(200, gin.H{
				"message": "Internal Server Error ",
			})
			c.Abort()
			// )
			return
		}

		// userID, err := app.Model.Users.ValidToken(cookie)
		// if err != nil {
		// 	// deleteCookie(c)
		// 	c.AbortWithStatus(http.StatusInternalServerError)
		// 	fmt.Println("Error fetching user info:", err)
		// 	return
		// }

		// if userID <= 0 {
		// 	// deleteCookie(c)
		// 	c.Redirect(http.StatusMovedPermanently, "/user/login/")
		// 	c.Abort()
		// 	return
		// }

		// Parse the token
		SIGNING_KEY := os.Getenv("SIGNING_KEY")
		if SIGNING_KEY == "" {
			c.JSON(200, gin.H{
				"message": "Internal Server Error ",
			})

			fmt.Println("Error fetching info from env: ", SIGNING_KEY)
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(cookie, &models.MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Verify the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("[error] Unexpected signing method: %v", token.Header["alg"])
			}

			// Return the secret key
			return []byte(SIGNING_KEY), nil
		})

		if err != nil {
			deleteCookie(c)
			c.AbortWithStatus(http.StatusInternalServerError)
			fmt.Println("Error fetching info from token:", err)
			return
		}
		// Check if the token is valid
		if claims, ok := token.Claims.(*models.MyCustomClaims); ok && token.Valid {
			// Set the username in the request context
			c.Header("Username", claims.Username)

			app.UserID = int(claims.UserID)
			app.Username = claims.Username
			app.isAuthenticated = true
			c.Next()
		} else {
			// Return an errordeleteCookie(c)
			c.Redirect(http.StatusMovedPermanently, "/user/login/")
			c.Abort()
			return
		}
	}
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check authentication logic here
		if true {
			c.Next()
		} else {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
		}
	}
}

func deleteCookie(c *gin.Context) {
	c.SetCookie("ldata", "", -1, "/", "", false, true)
}
