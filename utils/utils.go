package utils

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

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

func FormatDate(date time.Time) string {
	if date.IsZero() {
		date = time.Now()
	}

	duration := time.Since(date).Abs().Hours()
	days := duration / 24
	
    switch {
		case days < 1:
			return fmt.Sprintf("Today, %s", date.Format("15:04"))
		case days < 2:
			return fmt.Sprintf("Yesterday, %s", date.Format("15:04"))	
		case days < 7:
			return date.Format("Monday, 15:04")
		default:
			return date.Format("2006-01-02 15:04")
    }
}