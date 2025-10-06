package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wholesomeow/chatwrapper/cmd/app"
)

type FakeData struct {
	Data string
}

type Conversations struct {
	Conversations []string
}

func ListConversations(c *gin.Context) {
	// Call the function and process any errors
	conversations, err := app.ListConversations()
	if err != nil {
		msg := fmt.Sprintf("Conversation collection failed: %s", err)
		status, response := Response500(msg)
		c.JSON(status, response)
	}

	conversations_list := Conversations{
		Conversations: conversations,
	}

	data, _ := json.MarshalIndent(conversations_list, "", "  ")

	// Populate the context
	c.JSON(http.StatusOK, Response{
		Status:    http.StatusText(http.StatusOK),
		Message:   "Conversations listed successfully",
		Data:      data,
		Timestamp: time.Now(),
	})
}

func CreateConversation(c *gin.Context) {
	// Get content from "form"
	topic := c.PostForm("topic")

	// Create new conversation
	client := app.NewOpenAIClient()
	conversation, err := app.StartNewConversation(client, topic)
	if err != nil {
		msg := fmt.Sprintf("Failed to create new Conversation: %s", err)
		status, response := Response500(msg)
		c.JSON(status, response)
	}

	// Process data
	data, err := app.StructToJSON(*conversation)
	if err != nil {
		msg := fmt.Sprintf("Failed to parse conversation to JSON: %s", err)
		status, response := Response500(msg)
		c.JSON(status, response)
	}

	// Populate the context
	c.JSON(http.StatusOK, Response{
		Status:    http.StatusText(http.StatusOK),
		Message:   "Conversations listed successfully",
		Data:      data,
		Timestamp: time.Now(),
	})
}

func GetConversation(c *gin.Context) {
	// TODO: Implement ConversationID as passed param to then get conversation data
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

func DeleteConversation(c *gin.Context) {
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
