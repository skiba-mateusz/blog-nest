package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/skiba-mateusz/blog-nest/auth"
	"github.com/skiba-mateusz/blog-nest/forms"
	"github.com/skiba-mateusz/blog-nest/types"
	"github.com/skiba-mateusz/blog-nest/utils"
	"github.com/skiba-mateusz/blog-nest/views/blogs"
)

type blogHandler struct {
	userStore types.UserStore
	blogStore types.BlogStore
}

func NewBlogHanlder(userStore types.UserStore, blogStore types.BlogStore) *blogHandler {
	return &blogHandler{
		userStore: userStore,
		blogStore: blogStore,
	}
}

func (h *blogHandler) HandleCreateShow(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.GetUserFromContext(r.Context())

	categories, err := h.blogStore.GetCategories()
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	utils.Render(w, blogs.Create(blogs.CreateData{
		Title: "Create Blog | BlogNest",
		BlogForm: &forms.Form{},
		User: user,
		Categories: categories,
	}))
}

func (h *blogHandler) HandleBlogShow(w http.ResponseWriter, r *http.Request) {
	str := chi.URLParam(r, "blogID")
	blogID, _ := strconv.Atoi(str)

	blog, err := h.blogStore.GetBlogByID(blogID)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	user, _ := auth.GetUserFromContext(r.Context())
	
	utils.Render(w, blogs.Show(blogs.ShowData{
		Title: "Blog | BlogNest",
		Blog: blog,
		User: user,
	}))
}

func (h *blogHandler) HandleCreateBlog(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		utils.ClientError(w, "invalid form data", http.StatusBadRequest)
		return
	}

	categories, err := h.blogStore.GetCategories()
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	user, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		utils.Render(w, blogs.CreateBlogForm(categories, &forms.Form{}))
		auth.PermissionDenied(w, r)
		return
	}

	form := forms.New(r.PostForm)
	form.MinLength("content", 200)
	form.Required("category")
	form.Required("title")

	if !form.Valid() {
		utils.Render(w, blogs.CreateBlogForm(categories, form))
		return
	}

	categoryID, _ := strconv.Atoi(form.Values.Get("category"))

	blog := types.Blog{
		Title: form.Values.Get("title"),
		Content: form.Values.Get("content"),
		Category: types.Category{
			ID: categoryID,
		},
		User: types.User{
			ID: user.ID,
		},
	}

	blogID, err := h.blogStore.CreateBlog(blog)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	w.Header().Set("HX-Redirect", fmt.Sprintf("/blog/%d", blogID))
}