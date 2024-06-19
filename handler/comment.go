package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/skiba-mateusz/blog-nest/auth"
	"github.com/skiba-mateusz/blog-nest/forms"
	"github.com/skiba-mateusz/blog-nest/types"
	"github.com/skiba-mateusz/blog-nest/utils"
	"github.com/skiba-mateusz/blog-nest/views/comments"
	"github.com/skiba-mateusz/blog-nest/views/components"
)

type commentHandler struct {
	commentStore types.CommentStore
}

func NewCommentHandler(commentStore types.CommentStore) *commentHandler {
	return &commentHandler{commentStore: commentStore}
}

func (h commentHandler) HandleCreateLike(w http.ResponseWriter, r *http.Request) {
	commentID, value, err := parseCreateCommentLikeRequest(r)
	if err != nil {
		utils.ClientError(w, "invalid request data", http.StatusBadRequest)
		return
	}

	user, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		likes, err := h.commentStore.GetCommentLikes(commentID, user.ID)
		if err != nil {
			utils.ServerError(w, err)
			return 
		}
		auth.PermissionDenied(w, r)
		utils.Render(w, components.Reactions(likes, "comment", commentID))
		return
	}

	err = h.commentStore.CreateLike(value, commentID, user.ID)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	likes, err := h.commentStore.GetCommentLikes(commentID, user.ID)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	utils.Render(w, components.Reactions(likes, "comment", commentID))
}

func (h *commentHandler) HandleUpdateLike(w http.ResponseWriter, r *http.Request) {
	commentID, value, err := parseCreateCommentLikeRequest(r)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	user, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		auth.PermissionDenied(w, r)
		return
	}

	err = h.commentStore.UpdateLike(value, commentID, user.ID)
	if err != nil {
		utils.ServerError(w, err)
		return
	}
	likes, err := h.commentStore.GetCommentLikes(commentID, user.ID)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	utils.Render(w, components.Reactions(likes, "comment", commentID))
}

func (h commentHandler) HandleCreateComment(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		utils.ClientError(w, "invalid request data", http.StatusBadRequest)
		return
	}

	user, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		auth.PermissionDenied(w, r)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("parent_id", "content")

	if !form.Valid() {
		return
	}

	blogID, _ := strconv.Atoi(chi.URLParam(r, "blogID"))
	parentID, _ := strconv.Atoi(form.Values.Get("parent_id"))
	
	comment := types.Comment{
		Content: form.Values.Get("content"),
		ParentID: parentID,
		User: types.User{
			ID: user.ID,
			Username: user.Username,
		},
		Blog: types.Blog{
			ID: blogID,
		},
	} 

	commentID, err := h.commentStore.CreateComment(comment)
	if err != nil {
		utils.ServerError(w, err)
		return
	}
	
	comment.ID = commentID

	if parentID != 0 {
		c, err := h.commentStore.GetCommentsByBlogID(blogID, user.ID)
		if err != nil {
			utils.ServerError(w, err)
			return
		}

		comment := getComment(c, parentID)
		comment.Replies = getReplies(c, parentID)

		utils.Render(w, comments.Comment(comment, comment.Replies))
		return
	}

	utils.Render(w, comments.Comment(&comment, comment.Replies))
}

func getComment(comments []types.Comment, commentID int) *types.Comment {
	for _, comment := range comments {
		if comment.ID == commentID {
			return &comment
		}
	}
	return nil
}

func getReplies(comments []types.Comment, parentID int) []types.Comment {
	replies := []types.Comment{}
	for _, comment := range comments {
		if comment.ParentID == parentID {
			comment.Replies = getReplies(comments, comment.ID)
			replies = append(replies, comment)
		}
	} 
	return replies
}

func groupComments(comments []types.Comment, parentID int) []types.Comment {
	return getReplies(comments, parentID)
}

func parseCreateCommentLikeRequest(r *http.Request) (int, int, error) {
	if err := r.ParseForm(); err != nil {
		return 0, 0, err
	}

	commentID, err := strconv.Atoi(chi.URLParam(r, "commentID"))
	if err != nil {
		return 0, 0, err
	}

	value, err := strconv.Atoi(r.PostForm.Get("value"))
	if err != nil {
		return 0, 0, err
	}

	return commentID, value, nil
}