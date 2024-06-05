package handler

import (
	"net/http"

	"github.com/skiba-mateusz/blog-nest/views/home"
)

type homeHandler struct {}

func NewHomeHandler() *homeHandler {
	return &homeHandler{}
}

func (h *homeHandler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	Render(w, home.Index(home.IndexData{
		Title: "BlogNest | Explore Blogs",
	}))
}