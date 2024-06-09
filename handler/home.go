package handler

import (
	"net/http"

	"github.com/skiba-mateusz/blog-nest/auth"
	"github.com/skiba-mateusz/blog-nest/types"
	"github.com/skiba-mateusz/blog-nest/utils"
	"github.com/skiba-mateusz/blog-nest/views/home"
)

type homeHandler struct {
	userStore types.UserStore
}

func NewHomeHandler(userStore types.UserStore) *homeHandler {
	return &homeHandler{userStore: userStore}
}

func (h *homeHandler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.GetUserFromContext(r.Context())

	utils.Render(w, home.Index(home.IndexData{
		Title: "BlogNest | Explore Blogs",
		User: user,
	}))
}