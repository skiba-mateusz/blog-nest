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
	blogStore types.BlogStore
}

func NewHomeHandler(userStore types.UserStore, blogStore types.BlogStore) *homeHandler {
	return &homeHandler{
		userStore: userStore,
		blogStore: blogStore,
	}
}

func (h *homeHandler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.GetUserFromContext(r.Context())

	blogs, err := h.blogStore.GetBlogs()
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	utils.Render(w, home.Index(home.IndexData{
		Title: "BlogNest | Explore Blogs",
		Blogs: blogs,
		User: user,
	}))
}
