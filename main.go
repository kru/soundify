package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kru/soundify/core/api"
	"github.com/kru/soundify/core/database"
	"github.com/kru/soundify/core/helper"
	"github.com/kru/soundify/core/middleware"
	"github.com/kru/soundify/core/worker"
)

func main() {

	err := helper.LoadEnv(".env")
	if err != nil {
		log.Fatalf("error while loading env file %v", err)
		return
	}
	// start db connection
	database.Init()
	defer database.DB.Close()

	if err != nil {
		log.Fatalf("error while querying users %v", err)
		return
	}

	// Run the worker in a separate goroutine
	go worker.Run()

	router := http.NewServeMux()

	// html renderer

	// JSON API
	router.HandleFunc("POST /links", api.HandleLink)
	router.HandleFunc("POST /v2/links", api.HandleLinkV2)

	middlewares := middleware.CombineMiddleware(
		middleware.Logger,
		middleware.Auth,
	)

	server := http.Server{
		Addr:    ":8080",
		Handler: middlewares(router),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v\n", err)
		}
	}()

	fmt.Println("Server listening to port 8080")

	// Keep the main function running for the HTTP server
	select {}
}
