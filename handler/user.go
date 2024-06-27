package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/skiba-mateusz/blog-nest/auth"
	"github.com/skiba-mateusz/blog-nest/forms"
	"github.com/skiba-mateusz/blog-nest/types"
	"github.com/skiba-mateusz/blog-nest/utils"
	"github.com/skiba-mateusz/blog-nest/views/user"
)

type userHandler struct {
	userStore types.UserStore
	s3Uploader types.S3Uploader
}

func NewUserHandler(userStore types.UserStore, s3Uploader types.S3Uploader) *userHandler {
	return &userHandler{
		userStore: userStore,
		s3Uploader: s3Uploader,
	}
}

func (h *userHandler) HandleRegisterShow(w http.ResponseWriter, r *http.Request) {
	u, ok := auth.GetUserFromContext(r.Context())
	if ok {
		w.Header().Set("HX-Redirect", "/")
	}

	utils.Render(w, user.Register(user.RegisterData{
		Title: "Register | BlogNest",
		RegisterForm: &forms.Form{},
		User: u,
	}))
} 

func (h *userHandler) HandleProfileShow(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		utils.ClientError(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	profile, err := h.userStore.GetProfileByID(userID)
	if err != nil {
		utils.ServerError(w, err)
		return 
	}

	authUser, _ := auth.GetUserFromContext(r.Context())
	profile.IsOwner = authUser.ID == userID

	utils.Render(w, user.Profile(user.ProfileData{
		Title: fmt.Sprintf("%s's Profile | BlogNest", profile.Username),
		User: authUser,
		Profile: profile,
	}))
}

func (h *userHandler) HandleSettingsShow(w http.ResponseWriter, r *http.Request) {
	u, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		w.Header().Set("HX-Redirect", "/")
	}

	utils.Render(w, user.Settings(user.SettingsData{
		Title: "Settings | BlogNest",
		User: u,
		Form: &forms.Form{},
	}))
}

func (h *userHandler) HandleLoginShow(w http.ResponseWriter, r *http.Request) {
	u, ok := auth.GetUserFromContext(r.Context())
	if ok {
		w.Header().Set("HX-Redirect", "/")
	}

	utils.Render(w, user.Login(user.LoginData{
		Title: "Login | BlogNest",
		LoginForm: &forms.Form{},
		User: u,
	}))
}

func (h *userHandler) HandleRegisterUserStep1(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		utils.ClientError(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("username", "email", "password", "password_repeat")
	form.MinLength("username", 4)
	form.MinLength("password", 8)
	form.Email("email")
	form.PasswordsMatch("password", "password_repeat")

	if !form.Valid() {
		utils.Render(w, user.RegisterFormStep1(form))
		return
	}

	u := types.User{
		Username: form.Values.Get("username"),
		Email: form.Values.Get("email"),
		Password: form.Values.Get("password"),
	}

	existingUser, err := h.userStore.GetUserByEmail(u.Email)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	if existingUser != nil {
		form.Errors.Add("email", "User with this email already exists")
		utils.Render(w, user.RegisterFormStep1(form))
		return
	}

	hashedPassword, err := auth.HashPassword(u.Password)
	if err != nil {
		utils.ServerError(w, err)
		return
	}
	u.Password = hashedPassword

	id, err := h.userStore.CreateUser(u)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	_, err = auth.GenerateToken(w, id)	
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	utils.Render(w, user.RegisterFormStep2(&forms.Form{}))
} 

func (h *userHandler) HandleRegisterUserStep2(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.ClientError(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("bio")
	form.MaxLength("bio", 200)
	
	file, header, err := r.FormFile("avatar")
	if err != nil {
		form.Errors.Add("avatar", "Avatar is required")
	} else {
		defer file.Close()
	}
	
	if !form.Valid() {
		utils.Render(w, user.RegisterFormStep2(form))
		return
	}

	s3Key, err := h.s3Uploader.PutObject(file, header.Filename, "avatars")
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	u, _ := auth.GetUserFromContext(r.Context())

	user := types.User{
		ID: u.ID,
		AvatarPath: s3Key,
		Bio: form.Values.Get("bio"),
	}

	_, err = h.userStore.UpdateUser(user)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	w.Header().Set("HX-Redirect", "/")
}

func (h *userHandler) HandleLoginUser(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		utils.ClientError(w, "Invalid from data", http.StatusBadRequest)
	}

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.Email("email")

	if !form.Valid() {
		utils.Render(w, user.LoginForm(form))
		return
	}

	u := types.User{
		Email: form.Values.Get("email"),
		Password: form.Values.Get("password"),
	}

	existingUser, err := h.userStore.GetUserByEmail(u.Email)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	if existingUser == nil {
		form.Errors.Add("password", "Invalid email or password")
		utils.Render(w, user.LoginForm(form))
		return
	}

	if ok := auth.ComparePasswords(existingUser.Password, u.Password); !ok {
		form.Errors.Add("password", "Invalid email or password")
		utils.Render(w, user.LoginForm(form))
		return
	}

	_, err = auth.GenerateToken(w, existingUser.ID)	
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	w.Header().Set("HX-Redirect", "/")
}

func (h *userHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.ClientError(w, "invalid request data", http.StatusBadRequest)
		return
	}

	authUser, _ := auth.GetUserFromContext(r.Context())

	form := forms.New(r.PostForm)
	form.Required("username", "bio")
	form.MaxLength("bio", 200)
	form.MinLength("username", 4)

	if !form.Valid() {
		utils.Render(w, user.SettingsForm(authUser, form))
		return
	}

	s3Key := ""
	file, header, err := r.FormFile("avatar")
	if err != nil {
		if err == http.ErrMissingFile {
			log.Println("File is missing")
		} else {
			utils.ServerError(w, err)
			return
		}
	}
	if file != nil {
		defer file.Close()
		s3Key, err = h.s3Uploader.PutObject(file, header.Filename, "avatars")
		if err != nil {
			utils.ServerError(w, err)
			return
		}
	}

	u := types.User{
		ID: authUser.ID,
		Username: form.Values.Get("username"),
		Email: form.Values.Get("email"),
		Bio: form.Values.Get("bio"),
		AvatarPath: s3Key,
	}

	_, err = h.userStore.UpdateUser(u)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	w.Header().Set("HX-Redirect", "/user/settings")
}

func (h *userHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	auth.DestroyJWT(w)
	http.Redirect(w, r , "/user/login", http.StatusSeeOther)	
}