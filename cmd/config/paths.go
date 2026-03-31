package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const (
	AppName           = "arbor"
	GlobalConfigFile  = "arbor-config.yaml"
	ProjectConfigDir  = ".arbor"
	ProjectConfigFile = "project.yaml"
	CanvasFile        = "canvas.json"
)

// Returns a human-readable OS name for debugging and troubleshooting
func OSLabel() string {
	switch runtime.GOOS {
	case "darwin":
		return "macOS"
	case "linux":
		return "Linux"
	case "windows":
		return "Windows"
	default:
		return runtime.GOOS
	}
}

// GetGlobalConfigDir returns the OS-appropriate directory for Arbor's
// global config. It never creates the directory — call EnsureDir() after.
//
//	macOS:   ~/Library/Application Support/arbor
//	Linux:   $XDG_CONFIG_HOME/arbor  (fallback: ~/.config/arbor)
//	Windows: %APPDATA%\arbor
func GetGlobalConfigDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine user config directory: %w", err)
	}

	return filepath.Join(base, AppName), nil
}

// Returns the full path to arbor-config.yaml
func GetGlobalConfigPath() (string, error) {
	dir, err := GetGlobalConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, GlobalConfigFile), nil
}

// Returns the path to the .arbor/ directory inside a given
// project root. Does not verify the project root exists.
func GetProjectConfigDir(projectRoot string) string {
	return filepath.Join(projectRoot, ProjectConfigDir)
}

// Returns the full path to .arbor/project.yaml for a given project root
func GetProjectConfigPath(projectRoot string) string {
	return filepath.Join(GetProjectConfigDir(projectRoot), ProjectConfigFile)
}

// Returns the full path to .arbor/canvas.json for a given
// project root.
func GetCavnasPath(projectRoot string) string {
	return filepath.Join(GetProjectConfigDir(projectRoot), CanvasFile)
}

// Creates a directory and all parents if they don't exist
func EnsureDir(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create directory %q: %w", path, err)
	}

	return nil
}

// Creates the global config directory if needed
func EnsureGlobalConfigDir() (string, error) {
	dir, err := GetGlobalConfigDir()
	if err != nil {
		return "", err
	}

	if err := EnsureDir(dir); err != nil {
		return "", err
	}

	return dir, nil
}

// Creates .arbor/ inside the project root, if needed
func EnsureProjectConfigDir(projectRoot string) (string, error) {
	dir := GetProjectConfigDir(projectRoot)
	if err := EnsureDir(dir); err != nil {
		return "", err
	}

	return dir, nil
}
