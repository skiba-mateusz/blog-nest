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

	latestBlogs, err := h.blogStore.GetLatestBlogs()
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	blogs, totalBlogs, err := h.blogStore.GetBlogs(0, types.DefaultPageSize, "", "")
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	categories, err := h.blogStore.GetCategories()
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	totalPages := (totalBlogs + types.DefaultPageSize - 1) / 4

	utils.Render(w, home.Index(home.IndexData{
		Title: "BlogNest | Explore Blogs",
		Categories: categories,
		Blogs: blogs,
		LatestBlogs: latestBlogs,
		TotalBlogs: totalBlogs,
		Page: 1,
		TotalPages: totalPages,
		User: user,
	}))
}
