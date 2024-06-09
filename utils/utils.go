package utils

import (
	"context"
	"log"
	"net/http"

	"github.com/a-h/templ"
)

func Render(w http.ResponseWriter, c templ.Component) {
	c.Render(context.Background(), w)
}

func ClientError(w http.ResponseWriter, message string, code int) {
	http.Error(w, message, code)
}

func ServerError(w http.ResponseWriter, err error) {
	log.Printf("Server error: %v", err)
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

