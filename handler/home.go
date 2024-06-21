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

	blogs, totalBlogs, err := h.blogStore.GetBlogs(0, "")
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	latestBlogs, err := h.blogStore.GetLatestBlogs()
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	totalPages := (totalBlogs + 3) / 4

	utils.Render(w, home.Index(home.IndexData{
		Title: "BlogNest | Explore Blogs",
		Blogs: blogs,
		LatestBlogs: latestBlogs,
		TotalBlogs: totalBlogs,
		Page: 1,
		TotalPages: totalPages,
		User: user,
	}))
}
