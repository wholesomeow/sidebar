package cli

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/wholesomeow/chatwrapper/cmd/api"
)

var startServer = &cobra.Command{
	Use:   "start",
	Short: "Start the Sidebar backend",
	Long:  "Starts the Sidebar Gin backend for the Sidebar webapp",
	Run: func(cmd *cobra.Command, args []string) {
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
	},
}

func Start() {
	rootCmd.AddCommand(startServer)
}
