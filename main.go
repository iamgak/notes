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
	"github.com/iamgak/todo/pkg/others"
	_ "github.com/lib/pq"
)

type Application struct {
	// Config          Config
	Model           *models.Init
	UserID          int
	isAuthenticated bool
	Username        string
}

func main() {
	fmt.Print("To do Web App startet \n")

	err := others.LoadEnvVariables()
	if err != nil {
		panic(err)
	}

	Port := os.Getenv("PORT")
	addr := flag.String("addr", ":"+Port, "HTTP network address")
	flag.Parse()
	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	ctx := context.Background()
	client := InitRedis()
	if client == nil {
		panic(fmt.Errorf("redis client is %T", client))
	}

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

	go func() {
		//later on create a function to init all the function
		// app.Model.Redis.Subscribe(ctx)
		sub := client.Subscribe(ctx, "todo.notifications")
		ch := sub.Channel()
		for msg := range ch {
			fmt.Println("Received message:", msg)
		}
	}()

	log.Printf("[info] start http server listening %s", *addr)
	server.ListenAndServe()
}
