package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wholesomeow/chatwrapper/cmd/app"
)

type SingleMessage struct {
	Message string
}

func (h *Handler) SendMessage(c *gin.Context) {
	// Get content from "form"
	id := c.Param("id")
	msg := c.PostForm("message")

	// Load existing conversation
	convo, err := app.GetConversation(h.convoDir(), id)
	if err != nil {
		msg := fmt.Sprintf("Failed to load conversation: %s", err)
		status, response := Response500(msg)
		c.JSON(status, response)

		return
	}

	// Get conversation context to send with message
	content, err := app.SendMessage(h.Client, convo, msg)
	if err != nil {
		msg := fmt.Sprintf("Failed to send message: %s", err)
		status, response := Response500(msg)
		c.JSON(status, response)

		return
	}

	data, _ := json.MarshalIndent(SingleMessage{Message: content}, "", "  ")

	// Populate the context
	c.JSON(http.StatusOK, Response{
		Status:    http.StatusText(http.StatusOK),
		Message:   "Message sent successfully",
		Data:      data,
		Timestamp: time.Now(),
	})
}

func (h *Handler) GetMessages(c *gin.Context) {
	id := c.Param("id")

	convo, err := app.GetConversation(h.convoDir(), id)
	if err != nil {
		msg := fmt.Sprintf("Failed to get messages from conversation: %s", err)
		status, response := Response500(msg)
		c.JSON(status, response)

		return
	}

	// Populate the context
	c.JSON(http.StatusOK, Response{
		Status:    http.StatusText(http.StatusOK),
		Message:   "Retrieved messages successfully",
		Data:      convo.Messages,
		Timestamp: time.Now(),
	})
}

func (h *Handler) DeleteMessage(c *gin.Context) {
	// TODO: Implement MessageID as passed param to then delete message from conversation
	// if err != nil {
	// 	msg := fmt.Sprintf("Failed to delete message: %s", err)
	// 	status, response := Response500(msg)
	// 	c.JSON(status, response)
	// }

	// Populate the context
	c.JSON(http.StatusOK, Response{
		Status:  http.StatusText(http.StatusOK),
		Message: "Conversations listed successfully",
		Data: FakeData{
			Data: "Wow, look at all this data",
		},
		Timestamp: time.Now(),
	})
}
