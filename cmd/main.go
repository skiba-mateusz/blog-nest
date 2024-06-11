package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/skiba-mateusz/blog-nest/auth"
	"github.com/skiba-mateusz/blog-nest/config"
	"github.com/skiba-mateusz/blog-nest/db"
	"github.com/skiba-mateusz/blog-nest/handler"
	"github.com/skiba-mateusz/blog-nest/store"
)

func main() {
	router := chi.NewRouter()
	
	db, err := db.NewPostgreSQLStorage()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("DB connected")

	userStore := store.NewUserStore(db)
	blogStore := store.NewBlogStore(db)

	homeHandler := handler.NewHomeHandler(userStore, blogStore)
	userHandler := handler.NewUserHandler(userStore)
	blogHandler := handler.NewBlogHanlder(userStore, blogStore)

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(auth.WithJWT(userStore))

	router.Get("/", homeHandler.HandleIndex)
	router.Get("/blog/create", blogHandler.HandleCreateShow)
	router.Get("/blog/{blogID}", blogHandler.HandleBlogShow)
	router.Get("/user/register", userHandler.HandleRegisterShow)
	router.Get("/user/login", userHandler.HandleLoginShow)
	router.Get("/user/logout", userHandler.HandleLogout)

	router.Post("/user/register", userHandler.HandleRegisterUser)
	router.Post("/user/login", userHandler.HandleLoginUser)
	router.Post("/blog/create", blogHandler.HandleCreateBlog)


	fs := http.FileServer(http.Dir("./static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	log.Println("Server starting on", config.Envs.ListenAddr)

	log.Fatal(http.ListenAndServe(config.Envs.ListenAddr, router))
}