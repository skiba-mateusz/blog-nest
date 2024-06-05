package handler

import (
	"log"
	"net/http"
)

type homeHandler struct {}

func NewHomeHandler() *homeHandler {
	return &homeHandler{}
}

func (h *homeHandler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	log.Println("index page")
}