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

func SendMessage(c *gin.Context) {
	// Get content from "form"
	msg := c.PostForm("message")

	// Create new conversation
	client := app.NewOpenAIClient()
	content, err := app.SendMessage(client, msg)
	if err != nil {
		msg := fmt.Sprintf("Failed to send message: %s", err)
		status, response := Response500(msg)
		c.JSON(status, response)

		return
	}

	message := SingleMessage{
		Message: content,
	}

	data, _ := json.MarshalIndent(message, "", "  ")

	// Populate the context
	c.JSON(http.StatusOK, Response{
		Status:    http.StatusText(http.StatusOK),
		Message:   "Conversations listed successfully",
		Data:      data,
		Timestamp: time.Now(),
	})
}

func GetMessages(c *gin.Context) {
	// Call the function and process any errors
	// Frontend can send /:id/...
	id := c.Param("id")
	convo, err := app.GetConversation(id)
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
		Data:      convo.Messages,
		Timestamp: time.Now(),
	})
}

func DeleteMessage(c *gin.Context) {
	// TODO: Implement MessageID as passed param to then delete message from conversation
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
