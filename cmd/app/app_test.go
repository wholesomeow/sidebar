package app_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/wholesomeow/chatwrapper/cmd/app"
)

// Create mock client for testing purposes
// Temporarily commented out until it's implemented
// mock := &MockClient{Response: "This is a fake client response"}
// convo, err := StartNewConversation(mock, "Test Topic")

// Writing a test notes
// Each test function should only test one thing
// - Should name tests test purpose with prefix "Test" and pass in testing pointer
// - Full example
//   - func TestTestName(t *testing.T) {}

// - Should create a temp directory with tmp := t.TempDir()
//   - Full example:
// tmp := t.TempDir()
// origWd, _ := os.Getwd()
// defer os.Chdir(origWd)
// os.Chdir(tmp)

// - Should show error in test with t.Errorf() or t.Fatalf()
//   - Full example
// if err := os.MkdirAll(sidebarDir, 0755); err != nil {
// 	t.Fatalf("failed to make sidebar dir: %v", err)
// }
// if got != "testKey999" {
// 	t.Errorf("expected key testKey999, got %q", got)
// }

// - Test helper functions should start with t.Helper() in function
//   - Full example
// func TestTestName(t *testing.T) {
//     t.Helper()
// }

// Test Helpers ---------------------------------------------------------------

func CreateFakeConfig(t *testing.T) (*app.Config, string) {
	t.Helper()

	// Create fake config
	tmp := t.TempDir()

	sidebarDir := filepath.Join(tmp, ".sidebar")
	if err := os.MkdirAll(sidebarDir, 0755); err != nil {
		t.Fatal(err)
	}

	uuid, _ := app.CreateUUIDv4()
	convoName := fmt.Sprintf("convo_%s.json", uuid)
	convoPath := filepath.Join(sidebarDir, convoName)

	configPath := filepath.Join(sidebarDir, "sidebar-config.yaml")
	config := app.Config{
		APIKEY:                   "teeheeAPIKey",
		LastConversationID:       uuid,
		ConversationFileLocation: convoPath,
	}

	return &config, configPath
}

func CreateFakeMessage(t *testing.T) *app.Message {
	t.Helper()

	// Create fake message
	uuid, _ := app.CreateUUIDv4()
	parentUUID, _ := app.CreateUUIDv4()
	message := app.Message{
		MessageID: uuid,
		ParentIDs: []string{parentUUID},
		Timestamp: time.Now(),
		Content:   "Content Message Content",
	}

	return &message
}

// Test config.go -------------------------------------------------------------

// Test client.go -------------------------------------------------------------

// Test message.go ------------------------------------------------------------

// Pre-Test Steps
// Create Fake Config Struct and File
// Create Fake Covnersation Struct
// Populate Conversation Struct with Message(s)
// Write Conversation Struct to File

// Tests to write
// 1. Message isn't malformed on creation
// 2. ChatCompletion parsed
// 3. Validate LastMessageID is updated
// 4. Commit applied
//     A. File saved
//     B. Config updated
//     C. Head updated
// 5. Client returns error (403 - insufficient_quota)

// Test conversation.go -------------------------------------------------------

// Pre-Test Steps
// Create Fake Config Struct and File
// Create Fake Conversation Struct and File

// Previously used mock conversation lines
// mock := &app.MockClient{}
// convo, err := app.StartNewConversation(mock, mockTopic)
// if err != nil {
// 	t.Fatalf("unexpected error: %v", err)
// }

// Tests to write
// 1. New conversation file created
// 2. Client was used
// 3. Message created
// 4. ChatCompletion parsed
// 5. Commit applied
//     A. File saved
//     B. Config updated
//     C. Head updated
// 6. Client returns error (403 - insufficient_quota)

func TestStartNewConversation_UsesMockCompletion(t *testing.T) {
	tmp := t.TempDir()

	// Create .sidebar dir
	sidebarDir := filepath.Join(tmp, ".sidebar")
	if err := os.MkdirAll(sidebarDir, 0755); err != nil {
		t.Fatalf("failed to make sidebar dir: %v", err)
	}

	// Create a fake convo template file
	templateDir := filepath.Join(tmp, "templates")
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		t.Fatalf("failed to make templates dir: %v", err)
	}
	templateFile := filepath.Join(templateDir, "convo.json")
	if err := os.WriteFile(templateFile, []byte(`{"conversationID": "template"}`), 0644); err != nil {
		t.Fatalf("failed to write template file: %v", err)
	}

	// Create a fake config file
	configFile := filepath.Join(sidebarDir, "sidebar-config.yaml")
	if err := os.WriteFile(configFile, []byte("API_KEY: fake-key\nlastConversationID: template\n"), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	// Change into temp dir so StartNewConversation finds .sidebar
	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}

	client := &app.MockClient{}
	topic := "testing topic for chatbot wrapper application"

	convo, err := app.StartNewConversation(client, topic)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Assert last message was set from MockCompletion
	msg, ok := convo.Messages[convo.LastMessageID]
	if !ok {
		t.Fatalf("expected last message to exist in convo")
	}

	want := app.MockCompletion.Choices[0].Message.Content
	if msg.Content != want {
		t.Errorf("expected message content %q, got %q", want, msg.Content)
	}

	// Optional: check timestamp sanity
	if time.Since(msg.Timestamp) > 2*time.Second {
		t.Errorf("expected recent timestamp, got %v", msg.Timestamp)
	}
}

// Test archive.go ------------------------------------------------------------

// Pre-Test Steps
// Create Fake Config Struct and File
// Create Fake Conversation Struct and File

// Tests to write
// 1. Archive conversation
//     A. As Owner (success)
//     B. As Non-Owner (failure)
// 2. Attempt to write to conversation (should fail)
// 3. Unarchive conversation
//     A. As Owner (success)
//     B. As Non-Owner (failure)

// Test pin.go ----------------------------------------------------------------

// Pre-Test Steps
// Create Fake Config Struct and File
// Create Fake Conversation Struct and File

// Tests to write
// 1. Pin new message
// 2. List all pins
// 3. Reload conversation and list all pins
// 4. Unpin message
// 5. List all pins
// 6. Reload conversation and list all pins
