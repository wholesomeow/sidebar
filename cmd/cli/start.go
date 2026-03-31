package cli

import (
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/wholesomeow/chatwrapper/cmd/api"
	"github.com/wholesomeow/chatwrapper/cmd/config"
)

var startServer = &cobra.Command{
	Use:   "start",
	Short: "Start the Sidebar backend",
	Long:  "Starts the Sidebar Gin backend for the Sidebar webapp",
	Run: func(cmd *cobra.Command, args []string) {
		projectRoot, _ := cmd.Flags().GetString("project")
		if projectRoot == "" {
			projectRoot, _ = os.Getwd()
		}
		StartBackendServer(projectRoot)
	},
}

func Start() {
	startServer.Flags().String("project", "", "Path to the project root (default: current directory)")
	rootCmd.AddCommand(startServer)
}

func StartBackendServer(projectRoot string) {
	// Load both configs
	globalCfg, err := config.LoadGlobalConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load global config: %v\n", err)
	}

	projectCfg, err := config.LoadProjectConfig(projectRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load project config: %v\n", err)
	}

	fmt.Println("Server starting...")
	fmt.Printf("    Project: %s\n", projectCfg.Project.Name)
	fmt.Printf("    Default provider: %s / %s\n", globalCfg.Defaults.Provider, globalCfg.Defaults.Model)
	h := api.NewHandler(globalCfg, projectCfg, projectRoot)

	// Initializing CORS
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{"GET", "POST", "DELETE"},
		AllowHeaders: []string{"Content-Type"},
	}))

	endpoints := r.Group("/api/v1")
	{
		convo := endpoints.Group("/conversations")
		{
			convo.GET("", h.ListConversations)
			convo.POST("", h.CreateConversation)
			convo.GET("/:id", h.GetConversation)
			convo.DELETE("/:id", h.DeleteConversation)

			msg := convo.Group("/:id/messages")
			{
				msg.POST("", h.SendMessage)
				msg.GET("", h.GetMessages)
				msg.DELETE("/:msg_id", h.DeleteMessage)
			}
		}
	}

	r.Run(":8080")

	fmt.Println("Server successfully started!")
}
