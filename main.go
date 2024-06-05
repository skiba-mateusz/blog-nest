package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/joho/godotenv/autoload"
	"github.com/skiba-mateusz/blog-nest/config"
	"github.com/skiba-mateusz/blog-nest/handler"
)

func main() {
	router := chi.NewRouter()
	
	fs := http.FileServer(http.Dir("./static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	homeHandler := handler.NewHomeHandler()

	router.Get("/", homeHandler.HandleIndex)

	log.Println("Server starting on", config.Envs.ListenAddr)

	log.Fatal(http.ListenAndServe(config.Envs.ListenAddr, router))
}