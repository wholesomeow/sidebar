package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wholesomeow/chatwrapper/cmd/app"
)

type FakeData struct {
	Data string
}

func (h *Handler) ListConversations(c *gin.Context) {
	// Optional override from query param
	// otherwise, use project config path
	dir := c.Query("path")
	if dir == "" {
		dir = h.convoDir()
	}

	conversations, err := app.ListConversations(dir)
	if err != nil {
		msg := fmt.Sprintf("Conversation collection failed: %s", err)
		status, response := Response500(msg)
		c.JSON(status, response)

		return
	}

	// Populate the context
	c.JSON(http.StatusOK, Response{
		Status:    http.StatusText(http.StatusOK),
		Message:   "Conversations listed successfully",
		Data:      conversations,
		Timestamp: time.Now(),
	})
}

func (h *Handler) CreateConversation(c *gin.Context) {
	// Get content from "form"
	topic := c.PostForm("topic")

	// Create new conversation
	conversation, err := app.StartNewConversation(h.Client, topic, h.ProjectRoot)
	if err != nil {
		msg := fmt.Sprintf("Failed to create new Conversation: %s", err)
		status, response := Response500(msg)
		c.JSON(status, response)

		return
	}

	// Process data
	data, err := app.StructToJSON(*conversation)
	if err != nil {
		msg := fmt.Sprintf("Failed to parse conversation to JSON: %s", err)
		status, response := Response500(msg)
		c.JSON(status, response)

		return
	}

	// Populate the context
	c.JSON(http.StatusOK, Response{
		Status:    http.StatusText(http.StatusOK),
		Message:   "Conversations listed successfully",
		Data:      data,
		Timestamp: time.Now(),
	})
}

func (h *Handler) GetConversation(c *gin.Context) {
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
		Message:   "Conversations listed successfully",
		Data:      convo,
		Timestamp: time.Now(),
	})
}

func (h *Handler) DeleteConversation(c *gin.Context) {
	// TODO: Implement conversation deletion
	// if err != nil {
	// 	msg := fmt.Sprintf("NPC name generation failed: %s", err)
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
