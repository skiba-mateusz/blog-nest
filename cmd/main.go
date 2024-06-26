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
	"github.com/skiba-mateusz/blog-nest/s3uploader"
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

	s3Uploader, err := s3uploader.New(config.Envs.AWSRegion)
	if err != nil {
		log.Fatal("failed to initialize S3Uploader")		
	}

	userStore := store.NewUserStore(db)
	blogStore := store.NewBlogStore(db)
	commentStore := store.NewCommentStore(db)

	userHandler := handler.NewUserHandler(userStore, s3Uploader)
	commentHandler := handler.NewCommentHandler(commentStore)
	homeHandler := handler.NewHomeHandler(userStore, blogStore)
	blogHandler := handler.NewBlogHanlder(userStore, blogStore, commentStore, s3Uploader)

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(auth.WithJWT(userStore))

	router.Get("/", homeHandler.HandleIndex)

	router.Get("/user/register", userHandler.HandleRegisterShow)
	router.Get("/user/login", userHandler.HandleLoginShow)
	router.Get("/user/logout", userHandler.HandleLogout)
	router.Get("/user/profile/{userID}", userHandler.HandleProfileShow)
	router.Get("/user/settings", userHandler.HandleSettingsShow)
	router.Post("/user/register/step1", userHandler.HandleRegisterUserStep1)
	router.Post("/user/register/step2", userHandler.HandleRegisterUserStep2)
	router.Post("/user/login", userHandler.HandleLoginUser)
	router.Put("/user/update", userHandler.HandleUpdate)

	router.Get("/blog/create", blogHandler.HandleCreateShow)
	router.Get("/blog/{blogID}", blogHandler.HandleBlogShow)
	router.Get("/blog/page/{page}", blogHandler.HandleGetBlogs)
	router.Get("/blog/search", blogHandler.HandleSearchShow)
	router.Post("/blog/create", blogHandler.HandleCreateBlog)
	router.Post("/blog/{blogID}/like", blogHandler.HandleCreateLike)
	router.Post("/blog/{blogID}/comment", commentHandler.HandleCreateComment)
	router.Put("/blog/{blogID}/like", blogHandler.HandleUpdateLike)

	router.Post("/comment/{commentID}/like", commentHandler.HandleCreateLike)
	router.Put("/comment/{commentID}/like", commentHandler.HandleUpdateLike)

	

	fs := http.FileServer(http.Dir("./static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	log.Println("Server starting on", config.Envs.ListenAddr)

	log.Fatal(http.ListenAndServe(config.Envs.ListenAddr, router))
}