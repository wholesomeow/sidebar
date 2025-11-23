package cli

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/wholesomeow/chatwrapper/cmd/api"
)

var startServer = &cobra.Command{
	Use:   "start",
	Short: "Start the Sidebar backend",
	Long:  "Starts the Sidebar Gin backend for the Sidebar webapp",
	Run: func(cmd *cobra.Command, args []string) {
		StartBackendServer()
	},
}

func Start() {
	rootCmd.AddCommand(startServer)
}

func StartBackendServer() {
	fmt.Println("Server starting...")
	r := gin.Default()

	// Initializing CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{"GET"},
		AllowHeaders: []string{"Content-Type"},
	}))

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
				msg.GET("", api.GetMessages)
				msg.DELETE("/:msg_id", api.DeleteMessage)
			}
		}
	}

	r.Run(":8080")

	fmt.Println("Server successfully started!")
}
