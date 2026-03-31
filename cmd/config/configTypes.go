package config

import "time"

const CurrentConfigVersion = 1

// Global config types

type GlobalConfig struct {
	Version   int
	Providers ProvidersConfig `yaml:"providers"`
	Defaults  DefaultsConfig  `yaml:"defaults"`
	UI        UIConfig        `yaml:"ui"`
	Keybinds  KeybindsConfig  `yaml:"keybindings"`
}

type ProvidersConfig struct {
	OpenAI    OpenAIConfig    `yaml:"openai"`
	Anthropic AnthropicConfig `yaml:"anthropic"`
	Ollama    OllamaConfig    `yaml:"ollama"`
}

type OpenAIConfig struct {
	APIKey string `yaml:"api_key"`
}

type AnthropicConfig struct {
	APIKey string `yaml:"api_key"`
}

type OllamaConfig struct {
	// Endpoint is overridable so users can point at a remote Ollama server
	// (e.g. a home GPU box). Defaults to http://localhost:11434.
	Endpoint string `yaml:"endpoint"`
}

type DefaultsConfig struct {
	Provider string `yaml:"provider"` // "openai" | "anthropic" | "ollama"
	Model    string `yaml:"model"`
}

type UIConfig struct {
	Theme string  `yaml:"theme"` // "dark" | "light" | "system"
	Zoom  float64 `yaml:"zoom"`
}

type KeybindsConfig struct {
	Branch  string `yaml:"branch"`
	Merge   string `yaml:"merge"`
	NewNode string `yaml:"new_node"`
	Export  string `yaml:"export"`
}

// Per-project config types

type ProjectConfig struct {
	Version  int            `yaml:"version"`
	Project  ProjectMeta    `yaml:"project"`
	LLMRoles LLMRolesConfig `yaml:"llm_roles"`
	Memory   MemoryConfig   `yaml:"memory"`
}

type ProjectMeta struct {
	Name        string    `yaml:"name"`
	Description string    `yaml:"description"`
	CreatedAt   time.Time `yaml:"created_at"`
}

type LLMRolesConfig struct {
	// Chat is the primary conversational model.
	Chat LLMRole `yaml:"chat"`
	// DocParser handles background document ingestion
	DocParser LLMRole `yaml:"doc_parser"`
	// Compressor auto-summarises long branch histories
	Compressor LLMRole `yaml:"compressor"`
}

type LLMRole struct {
	Provider string `yaml:"provider"` // "openai" | "anthropic" | "ollama"
	Model    string `yaml:"model"`
}

type MemoryConfig struct {
	Enabled    bool `yaml:"enabled"`
	MaxEntries int  `yaml:"max_entries"`
}

// Default config types

func DefaultGlobalConfig() *GlobalConfig {
	return &GlobalConfig{
		Version: CurrentConfigVersion,
		Providers: ProvidersConfig{
			OpenAI:    OpenAIConfig{APIKey: ""},
			Anthropic: AnthropicConfig{APIKey: ""},
			Ollama:    OllamaConfig{Endpoint: "http://localhost:11434"},
		},
		Defaults: DefaultsConfig{
			Provider: "ollama",
			Model:    "llama3.2",
		},
		UI: UIConfig{
			Theme: "system",
			Zoom:  1.0,
		},
		Keybinds: KeybindsConfig{
			Branch:  "mod+b",
			Merge:   "mod+m",
			NewNode: "mod+enter",
			Export:  "mod+shift+e",
		},
	}
}

func DefaultProjectConfig(name string) *ProjectConfig {
	return &ProjectConfig{
		Version: CurrentConfigVersion,
		Project: ProjectMeta{
			Name:        name,
			Description: "",
			CreatedAt:   time.Now().UTC(),
		},
		LLMRoles: LLMRolesConfig{
			Chat:       LLMRole{Provider: "ollama", Model: "llama3.2"},
			DocParser:  LLMRole{Provider: "ollama", Model: "llama3.2"},
			Compressor: LLMRole{Provider: "ollama", Model: "llama3.2"},
		},
		Memory: MemoryConfig{
			Enabled:    true,
			MaxEntries: 100,
		},
	}
}
