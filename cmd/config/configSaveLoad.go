package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

// Writes YAML data to a file
func writeYAML(path string, v any) error {
	data, err := yaml.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshalling config: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing config to %q: %w", path, err)
	}
	return nil
}

// projectNameFromRoot derives a human-readable project name from the
// directory name of the project root. Falls back to "untitled" if the
// path is empty or just slashes.
func projectNameFromRoot(projectRoot string) string {
	name := filepath.Base(projectRoot)
	if name == "." || name == "/" || name == "" {
		return "untitled"
	}
	return name
}

// SaveGlobalConfig marshals cfg to YAML and writes it to the global config
// path, creating the config directory if needed.
func SaveGlobalConfig(cfg *GlobalConfig) error {
	dir, err := EnsureGlobalConfigDir()
	if err != nil {
		return err
	}

	path := filepath.Join(dir, GlobalConfigFile)
	return writeYAML(path, cfg)
}

// SaveProjectConfig marshals cfg to YAML and writes it to
// <projectRoot>/.arbor/project.yaml, creating .arbor/ if needed.
func SaveProjectConfig(projectRoot string, cfg *ProjectConfig) error {
	if _, err := EnsureProjectConfigDir(projectRoot); err != nil {
		return err
	}

	path := GetProjectConfigPath(projectRoot)
	return writeYAML(path, cfg)
}

// LoadGlobalConfig loads the global config from disk. If no file exists,
// it writes defaults and returns them - first-run is handled transparently.
func LoadGlobalConfig() (*GlobalConfig, error) {
	path, err := GetGlobalConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("reading global config: %w", err)
		}
		// First run — write defaults and return them.
		cfg := DefaultGlobalConfig()
		if err := SaveGlobalConfig(cfg); err != nil {
			return nil, fmt.Errorf("writing default global config: %w", err)
		}
		return cfg, nil
	}

	var cfg GlobalConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing global config at %q: %w", path, err)
	}

	// Future: if cfg.Version < CurrentConfigVersion, run migration here.

	return &cfg, nil
}

// LoadProjectConfig loads a project config from <projectRoot>/.arbor/project.yaml.
// If no file exists, it writes defaults named after the project root's
// directory name and returns them.
func LoadProjectConfig(projectRoot string) (*ProjectConfig, error) {
	path := GetProjectConfigPath(projectRoot)

	data, err := os.ReadFile(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("reading project config: %w", err)
		}
		// First run — derive project name from directory name.
		name := projectNameFromRoot(projectRoot)
		cfg := DefaultProjectConfig(name)
		if err := SaveProjectConfig(projectRoot, cfg); err != nil {
			return nil, fmt.Errorf("writing default project config: %w", err)
		}
		return cfg, nil
	}

	var cfg ProjectConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing project config at %q: %w", path, err)
	}

	return &cfg, nil
}
