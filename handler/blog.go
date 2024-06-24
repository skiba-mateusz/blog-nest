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

func (h *blogHandler) HandleSearchShow(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		utils.ClientError(w, "invalid request data", http.StatusBadRequest)
		return
	}

	user, _ := auth.GetUserFromContext(r.Context())

	searchQuery := r.URL.Query().Get("search_query")
	category := r.URL.Query().Get("category_name")

	blogsSlice, totalBlogs, err := h.blogStore.GetBlogs(0, types.DefaultPageSize, searchQuery, category)
	if err != nil {
		utils.ServerError(w, err)
		return 
	}

	totalPages := (totalBlogs + types.DefaultPageSize - 1) / types.DefaultPageSize

	utils.Render(w, blogs.Search(blogs.SearchData{
		Title: fmt.Sprintf("'%s' Results | BlogNest", searchQuery),
		User: user,
		SearchQuery: searchQuery,
		Category: category,
		Blogs: blogsSlice,
		TotalBlogs: totalBlogs,
		Page: 1,
		TotalPages: totalPages,
	}))
}

func (h *blogHandler) HandleBlogShow(w http.ResponseWriter, r *http.Request) {
	blogID, err := strconv.Atoi(chi.URLParam(r, "blogID"))
	if err != nil {
		utils.ClientError(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	user, _ := auth.GetUserFromContext(r.Context())

	blog, err := h.blogStore.GetBlogByID(blogID, user.ID)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	comments, err := h.commentStore.GetCommentsByBlogID(blog.ID, user.ID)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	groupedComments := groupComments(comments, 0)

	utils.Render(w, blogs.Index(blogs.IndexData{
		Title: fmt.Sprintf("%s | BlogNest", blog.Title),
		Blog: blog,
		Comments: groupedComments,
		User: user,
		CommentForm: &forms.Form{},
	}))
}

func (h *blogHandler) HandleGetBlogs(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(chi.URLParam(r, "page"))
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * 4
	searchQuery := r.URL.Query().Get("search_query")
	category := r.URL.Query().Get("category")

	blogsSlice, totalBlogs, err := h.blogStore.GetBlogs(offset, types.DefaultPageSize, searchQuery, category)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	totalPages := (totalBlogs + types.DefaultPageSize - 1) / types.DefaultPageSize

	utils.Render(w, blogs.List(blogs.ListData{
		Blogs: blogsSlice,
		TotalBlogs: totalBlogs,
		Page: page,
		TotalPages: totalPages,
		SearchQuery: searchQuery,
	}))
}

func (h *blogHandler) HandleCreateBlog(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		utils.ClientError(w, "invalid requesta data", http.StatusBadRequest)
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
	form.Required("category", "title")

	if !form.Valid() {
		utils.Render(w, blogs.CreateBlogForm(categories, form))
		return
	}

	categoryID, err := strconv.Atoi(form.Values.Get("category"))
	if err != nil {
		utils.ClientError(w, "invalid category ID", http.StatusBadRequest)
		return
	}

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
	blogID, value, err := parseCreateBlogLikeRequest(r)
	if err != nil {
		utils.ClientError(w, "invalid request data", http.StatusBadRequest)
		return
	}

	if !validateCreateBlogLikeRequest(r) {
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
	blogID, value, err := parseCreateBlogLikeRequest(r)
	if err != nil {
		utils.ClientError(w, "invalid request data", http.StatusBadRequest)
		return
	}

	if !validateCreateBlogLikeRequest(r) {
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

func parseCreateBlogLikeRequest(r *http.Request) (int, int, error) {
	if err := r.ParseForm(); err != nil {
		return 0, 0, err
	}

	blogID, err := strconv.Atoi(chi.URLParam(r, "blogID"))
	if err != nil {
		return 0, 0, err
	}

	value, err := strconv.Atoi(r.PostForm.Get("value"))
	if err != nil {
		return 0, 0, err
	}

	return blogID, value, nil
}

func validateCreateBlogLikeRequest(r *http.Request) bool {
	form := forms.New(r.PostForm)
	form.Required("value")
	form.MinValue("value", -1)
	form.MaxValue("value", 1)

	return form.Valid()
}
