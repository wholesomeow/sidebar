package api

import (
	"path/filepath"

	"github.com/wholesomeow/chatwrapper/cmd/app"
	"github.com/wholesomeow/chatwrapper/cmd/config"
)

// Holds all the shared dependencies for all API handlers.
// Constructed once at server startup.
type Handler struct {
	Global      *config.GlobalConfig
	Project     *config.ProjectConfig
	Client      app.ChatClient
	ProjectRoot string
}

func NewHandler(global *config.GlobalConfig, project *config.ProjectConfig, projectRoot string) *Handler {
	client := app.NewOpenAIClient(global)
	return &Handler{
		Global:      global,
		Project:     project,
		Client:      client,
		ProjectRoot: projectRoot,
	}
}

// convoDir is a convenience method so handlers don't construct paths manually.
func (h *Handler) convoDir() string {
	return filepath.Join(h.ProjectRoot, ".arbor", "conversations")
}
