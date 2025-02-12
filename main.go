package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iamgak/todo/models"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type ToDo struct {
	Title       string
	Description string
	Visibility  bool
	Editable    bool
	Deleted     bool
	Updated     bool
	version     int
	// Created     date
}

type Application struct {
	// Config          Config
	Model           *models.Init
	Uid             int
	isAuthenticated bool

	Username string
}

// type Config struct {
// 	port int
// 	env  string
// 	// db   struct {
// 	// 	dsn          string
// 	// 	maxOpenConns int
// 	// 	maxIdleConns int
// 	// 	maxIdleTime  string
// 	// }
// }

// const version = "1.0.0.0"

func main() {
	fmt.Print("To do Web App startet \n")
	// var cfg Config

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	dbUser := os.Getenv("DB_USERNAME")
	dbName := os.Getenv("DB_DATABASE")
	Port := os.Getenv("PORT")
	dbPassword := os.Getenv("DB_PASSWORD")

	addr := flag.String("addr", ":"+Port, "HTTP network address")
	// dsn := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName)
	dsn := flag.String("dsn", fmt.Sprintf("%s:%s@/%s?parseTime=true", dbUser, dbPassword, dbName), "MySQL data source name")

	flag.Parse()
	db, err := openDB(*dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	ctx := context.Background()
	redis_name := "localhost"
	redis_password := ""
	client := InitRedis(redis_name, redis_password)
	app := Application{
		Model: models.Constructor(db, client),
		// Config: cfg,
	}

	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:         *addr,
		Handler:      app.InitRouter(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		// MaxHeaderBytes: 1 << 20,
		MaxHeaderBytes: maxHeaderBytes,
	}

	log.Printf("[info] start http server listening %s", *addr)

	go func() {
		// msg := []byte("New to-do item added")

		sub := client.Subscribe(ctx, "todo.notifications")
		ch := sub.Channel()
		for msg := range ch {
			fmt.Println("Received message:", msg)
		}

	}()

	server.ListenAndServe()
}
