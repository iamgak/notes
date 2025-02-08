package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func InitRedis(name, passw string) *redis.Client {

	redis_port := 6379
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", name, redis_port),
		Password: passw, // no password set
		DB:       0,     // use default DB
	})

	return client
}

func InitGoogleOAuth() (*oauth2.Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
		return nil, err
	}

	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}, nil
}
