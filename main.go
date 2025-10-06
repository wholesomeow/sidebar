package main

import (
	"github.com/gin-gonic/gin"
	"github.com/wholesomeow/chatwrapper/cmd/api"
	"github.com/wholesomeow/chatwrapper/cmd/cli"
)

func main() {
	r := gin.Default()

	endpoints := r.Group("/api/v1")
	{
		convo := endpoints.Group("/conversations")
		{
			convo.GET("", api.ListConversations)
			convo.POST("", api.CreateConversation)
			convo.GET("/:id", api.GetConversation)
			convo.DELETE("/:id", api.DeleteConversation)

			msg := convo.Group("/:id/messages")
			{
				msg.POST("", api.SendMessage)
				msg.GET("/:msg_id", api.GetMessage)
				msg.DELETE("/:msg_id", api.DeleteMessage)
			}
		}
	}

	r.Run(":8080")

	cli.Execute()
}
