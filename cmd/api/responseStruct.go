package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status    string
	Message   string
	Data      interface{}
	Error     interface{}
	Timestamp time.Time
}

func Response404(context *gin.Context) {
	status := http.StatusNotFound
	response := Response{
		Status:    http.StatusText(status),
		Message:   "Page not found",
		Timestamp: time.Now(),
	}
	context.JSON(status, response)
}

func Response500(msg string) (int, Response) {
	// Return error message if querystring parameters don't pass
	status := http.StatusInternalServerError
	response := Response{
		Status:    http.StatusText(status),
		Message:   msg,
		Timestamp: time.Now(),
	}

	return status, response
}
