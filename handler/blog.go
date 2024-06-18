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
	"github.com/skiba-mateusz/blog-nest/views/components"
)

type blogHandler struct {
	userStore 		types.UserStore
	blogStore 		types.BlogStore
	commentStore 	types.CommentStore
}

func NewBlogHanlder(userStore types.UserStore, blogStore types.BlogStore, commentStore types.CommentStore) *blogHandler {
	return &blogHandler{
		userStore: userStore,
		blogStore: blogStore,
		commentStore: commentStore,
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

	var userID int
	user, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		userID = 0
	} else {
		userID = user.ID
	}

	blog, err := h.blogStore.GetBlogByID(blogID)
	if err != nil {
		utils.ServerError(w, err)
		return
	}
	blogLikes, err := h.blogStore.GetBlogLikes(userID, blog.ID)
	if err != nil {
		utils.ServerError(w, err)
		return
	}
	comments, err := h.commentStore.GetCommentsByBlogID(blog.ID, userID)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	groupedComments := groupComments(comments, 0)
	blog.Likes = blogLikes

	utils.Render(w, blogs.Show(blogs.ShowData{
		Title: "Blog | BlogNest",
		Blog: blog,
		Comments: groupedComments,
		User: user,
		CommentForm: &forms.Form{},
	}))
}

func (h *blogHandler) HandleCreateBlog(w http.ResponseWriter, r *http.Request) {
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

func (h *blogHandler) HandleCreateLike(w http.ResponseWriter, r *http.Request) {
	blogID, value, err := parseLikeRequest(r)
	if err != nil {
		utils.ClientError(w, "invalid request data", http.StatusBadRequest)
		return
	}

	user, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		blogLikes, err := h.blogStore.GetBlogLikes(0, blogID)
		if err != nil {
			utils.ServerError(w, err)
			return
		}
		utils.Render(w, components.Reactions(blogLikes, "blog", blogID))
		auth.PermissionDenied(w, r)
		return
	}

	err = h.blogStore.CreateLike(user.ID, blogID, value)
	if err != nil {
		utils.ServerError(w, err)
		return
	}
	blogLikes, err := h.blogStore.GetBlogLikes(user.ID, blogID)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	utils.Render(w, components.Reactions(blogLikes, "blog", blogID))
}

func (h *blogHandler) HandleUpdateLike(w http.ResponseWriter, r *http.Request) {
	blogID, value, err := parseLikeRequest(r)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	user, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		auth.PermissionDenied(w, r)
		return
	}

	err = h.blogStore.UpdateLike(user.ID, blogID, value)
	if err != nil {
		utils.ServerError(w, err)
		return
	}
	blogLikes, err := h.blogStore.GetBlogLikes(user.ID, blogID)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	utils.Render(w, components.Reactions(blogLikes, "blog", blogID))
}

func parseLikeRequest(r *http.Request) (int, int, error) {
	if err := r.ParseForm(); err != nil {
		return 0, 0, err
	}

	str := chi.URLParam(r, "blogID")
	blogID, _ := strconv.Atoi(str)

	valueStr := r.PostForm.Get("value")
	value, _ := strconv.Atoi(valueStr)

	return blogID, value, nil
}
