package handler

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
)

func Render(w http.ResponseWriter, c templ.Component) {
	c.Render(context.Background(), w)
}