package app_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/wholesomeow/chatwrapper/cmd/app"
	"gopkg.in/yaml.v2"
)

// Create mock client for testing purposes
// Temporarily commented out until it's implemented
// mock := &MockClient{Response: "This is a fake client response"}
// convo, err := StartNewConversation(mock, "Test Topic")

// Test config.go -------------------------------------------------------------

// helper to write a yaml config file
func writeConfigFile(t *testing.T, dir string, apiKey string) string {
	t.Helper()
	cfgPath := filepath.Join(dir, "sidebar-config.yaml")
	data := map[string]string{"API_KEY": apiKey}
	out, err := yaml.Marshal(data)
	if err != nil {
		t.Fatalf("failed to marshal yaml: %v", err)
	}
	if err := os.WriteFile(cfgPath, out, 0644); err != nil {
		t.Fatalf("failed to write yaml file: %v", err)
	}
	return cfgPath
}

func TestSetupConfig(t *testing.T) {
	tmp := t.TempDir()
	// Override path in function by faking current dir
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmp)

	// place template file
	templateDir := filepath.Join(tmp, "templates")
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		t.Fatalf("failed to make templates dir: %v", err)
	}
	templateFile := filepath.Join(templateDir, "sidebar-config.yaml")
	if err := os.WriteFile(templateFile, []byte("API_KEY: testkey"), 0644); err != nil {
		t.Fatalf("failed to write template file: %v", err)
	}

	// run
	if err := app.SetupConfig(); err != nil {
		t.Errorf("SetupConfig failed: %v", err)
	}

	// check created file exists
	if _, err := os.Stat("./.sidebar/sidebar-config.yaml"); os.IsNotExist(err) {
		t.Errorf("expected config file to be copied, but not found")
	}
}

func TestInitAPIKey_Success(t *testing.T) {
	tmp := t.TempDir()
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmp)

	// prepare config dir + file
	sidebarDir := filepath.Join(tmp, ".sidebar")
	if err := os.MkdirAll(sidebarDir, 0755); err != nil {
		t.Fatalf("failed to make sidebar dir: %v", err)
	}
	writeConfigFile(t, sidebarDir, "abc1234567")

	displayKey, err := app.InitAPIKey()
	if err != nil {
		t.Fatalf("InitAPIKey returned error: %v", err)
	}

	// check masked format
	if len(displayKey) == 0 || displayKey[:3] != "abc" {
		t.Errorf("unexpected displayKey: %s", displayKey)
	}
	// check env var set
	if got := os.Getenv("OPENAI_API_KEY"); got != "abc1234567" {
		t.Errorf("expected env var set, got %q", got)
	}
}

func TestInitAPIKey_MissingKey(t *testing.T) {
	tmp := t.TempDir()
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmp)

	sidebarDir := filepath.Join(tmp, ".sidebar")
	if err := os.MkdirAll(sidebarDir, 0755); err != nil {
		t.Fatalf("failed to make sidebar dir: %v", err)
	}
	// write file with missing API_KEY
	cfgPath := filepath.Join(sidebarDir, "sidebar-config.yaml")
	if err := os.WriteFile(cfgPath, []byte(""), 0644); err != nil {
		t.Fatalf("failed to write empty config: %v", err)
	}

	_, err := app.InitAPIKey()
	if err == nil {
		t.Errorf("expected error for missing API_KEY, got nil")
	}
}

func TestGetAPIKey(t *testing.T) {
	tmp := t.TempDir()
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmp)

	sidebarDir := filepath.Join(tmp, ".sidebar")
	if err := os.MkdirAll(sidebarDir, 0755); err != nil {
		t.Fatalf("failed to make sidebar dir: %v", err)
	}
	writeConfigFile(t, sidebarDir, "testKey999")

	got := app.GetAPIKey()
	if got != "testKey999" {
		t.Errorf("expected key testKey999, got %q", got)
	}
}

// Test client.go -------------------------------------------------------------

// Test conversation.go -------------------------------------------------------

func TestStartNewConversation_Success(t *testing.T) {
	// TODO: Fix this so it can see the sidebar-config.yaml. Might need to implement a test config file or something
	tmp := t.TempDir()

	// Set up fake .sidebar dir with template
	sidebarDir := filepath.Join(tmp, ".sidebar")
	if err := os.MkdirAll(sidebarDir, 0755); err != nil {
		t.Fatal(err)
	}
	templateDir := filepath.Join(tmp, "templates")
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		t.Fatal(err)
	}
	templateFile := filepath.Join(templateDir, "convo.json")
	if err := os.WriteFile(templateFile, []byte(`{"conversationID": "template"}`), 0644); err != nil {
		t.Fatal(err)
	}

	// Move cwd so StartNewConversation sees this temp .sidebar
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmp)

	// Fake client returns a canned ChatCompletion
	mockTopic := "This is an API test."
	expectedContent := "Got it — you said: \"This is an API test.\" Everything looks good!"

	// Run
	mock := &app.MockClient{}
	convo, err := app.StartNewConversation(mock, mockTopic)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 1. New conversation file created
	if convo.Path == "" {
		t.Errorf("expected convo.Path to be set")
	}
	if _, err := os.Stat(convo.Path); err != nil {
		t.Errorf("expected convo file on disk, got %v", err)
	}

	// 2. Client was used
	if convo.Topic != mockTopic {
		t.Errorf("expected topic %q, got %q", mockTopic, convo.Topic)
	}

	// 3. Message created
	msg, ok := convo.Messages[convo.LastMessageID]
	if !ok {
		t.Fatalf("expected message with LastMessageID in convo.Messages")
	}

	// 4. ChatCompletion parsed
	if msg.Content != expectedContent {
		t.Errorf("expected message content %q, got %q", expectedContent, msg.Content)
	}

	// 5. Commit applied (Head updated)
	if convo.Head == "" {
		t.Errorf("expected convo.Head to be set after commit")
	}
}

func TestStartNewConversation_ErrorResponseParsesOpenAIError(t *testing.T) {
	tmp := t.TempDir()

	// Set up dirs
	sidebarDir := filepath.Join(tmp, ".sidebar")
	if err := os.MkdirAll(sidebarDir, 0755); err != nil {
		t.Fatal(err)
	}
	templateDir := filepath.Join(tmp, "templates")
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		t.Fatal(err)
	}
	templateFile := filepath.Join(templateDir, "convo.json")
	if err := os.WriteFile(templateFile, []byte(`{"conversationID": "template"}`), 0644); err != nil {
		t.Fatal(err)
	}

	// Move cwd
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmp)

	// Fake client returns error with JSON payload
	openaiErr := app.OpenAIError{
		Message: "quota exceeded",
		Type:    "insufficient_quota",
		Code:    "403",
	}
	errJSON, _ := json.Marshal(openaiErr)

	// Mock client that always returns an error
	mock := &app.MockClient{
		Err: fmt.Errorf("OpenAI error: %s", string(errJSON)),
	}

	// Run
	convo, err := app.StartNewConversation(mock, "test-topic")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check that message content came from parsed error
	msg, ok := convo.Messages[convo.LastMessageID]
	if !ok {
		t.Fatalf("expected message created")
	}
	if msg.Content != openaiErr.Message {
		t.Errorf("expected message content %q, got %q", openaiErr.Message, msg.Content)
	}
}

// Test message.go ------------------------------------------------------------

// Test archive.go ------------------------------------------------------------

// Test pin.go ----------------------------------------------------------------

// Test toygit.go -------------------------------------------------------------

// Test utilities.go ----------------------------------------------------------
