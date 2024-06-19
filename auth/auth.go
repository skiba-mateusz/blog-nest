package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/skiba-mateusz/blog-nest/config"
	"github.com/skiba-mateusz/blog-nest/types"
	"github.com/skiba-mateusz/blog-nest/utils"
	"github.com/skiba-mateusz/blog-nest/views/user"
)

type userContextKey string
const UserKey userContextKey = "user"

func GenerateToken(w http.ResponseWriter, id int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpiration)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": strconv.Itoa(id),
		"exp": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.Envs.JWTSecret))
	if err != nil {
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name: "token",
		Value: tokenString,
		Expires: time.Now().Add(expiration),
		Path: "/",
		HttpOnly: true,
	})

	return tokenString, nil
}

func DestroyJWT(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name: "token",
		Value: "",
		Expires: time.Unix(0, 0),
		Path: "/",
		HttpOnly: true,
	})
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return []byte(config.Envs.JWTSecret), nil
	})
}

func WithJWT(userStore types.UserStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := getTokenFromCookie(r)

			if err != nil {
				log.Printf("failed to get token from cookie: %v", err)
				next.ServeHTTP(w, r)
				return
			}
	
			token, err := ParseToken(tokenString)
			if err != nil {
				log.Printf("failed to parse token: %v", err)
				next.ServeHTTP(w, r)
				return
			}
	
			if !token.Valid {
				log.Println("invalid token")
				next.ServeHTTP(w, r)
				return
			}
	
			claims := token.Claims.(jwt.MapClaims)
			str := claims["userID"].(string)
			userID, _ := strconv.Atoi(str)
	
			user, err := userStore.GetUserByID(userID)
			if err != nil {
				log.Println("failed to get user")
				next.ServeHTTP(w, r)
				return
			}
			user.Password = ""

			ctx := r.Context()
			ctx = context.WithValue(ctx, UserKey, user)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func GetUserFromContext(ctx context.Context) (*types.User, bool) {
	user, ok := ctx.Value(UserKey).(*types.User)
	if !ok {
		return &types.User{
			ID: 0,
		}, ok
	}
	return user, ok
}

func PermissionDenied(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	} else {
		utils.Render(w, user.LoginRequiredModal())
	}
}

func getTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("token")

	if err != nil {
		return "", err
	}
	
	return string(cookie.Value), nil
}

