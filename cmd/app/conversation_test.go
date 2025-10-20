package app_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wholesomeow/chatwrapper/cmd/app"
)

// Test Helpers ---------------------------------------------------------------

func setupFakeSidebarEnv(t *testing.T) string {
	t.Helper()

	// Create temp directories for testing
	tmp := t.TempDir()
	sidebarDir := filepath.Join(tmp, ".sidebar")
	templateDir := filepath.Join(tmp, "templates")

	require.NoError(t, os.MkdirAll(sidebarDir, 0755))
	require.NoError(t, os.MkdirAll(templateDir, 0755))

	// Fake template file
	templateFile := filepath.Join(sidebarDir, "convo.json")
	require.NoError(t, os.WriteFile(templateFile, []byte(`{"conversationID":"template"}`), 0644))

	// Fake config file
	configFile := filepath.Join(sidebarDir, "sidebar-config.yaml")
	require.NoError(t, os.WriteFile(configFile, []byte("API_KEY: fake-key\nlastConversationID: template\n"), 0644))

	// Change current working dir so relative paths work
	orignWD, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orignWD) })
	require.NoError(t, os.Chdir(tmp))

	return tmp
}

// Test conversation.go -------------------------------------------------------

// Pre-Test Steps
// Create Fake Config Struct and File
// Create Fake Conversation Struct and File

// Tests to write
// 5. Commit applied
//     A. File saved
//     B. Config updated
//     C. Head updated
// 6. Client returns error (403 - insufficient_quota)
//     A. Compile the internet error codes that OpenAI returns

func TestStartNewConversation_TestTable(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(t *testing.T, baseDir string)
		mockClient    app.ChatClient
		expectErr     bool
		expectErrPart string
		verify        func(t *testing.T, convo *app.Conversation)
	}{
		{
			name: "Happy Path",
			setup: func(t *testing.T, baseDir string) {
				setupFakeSidebarEnv(t)
			},
			mockClient: &app.MockClient{
				MockResponse: &app.MockCompletion,
			},
			expectErr: false,
			verify: func(t *testing.T, convo *app.Conversation) {
				msg := convo.Messages[convo.LastMessageID]
				require.NotNil(t, msg)
				require.Contains(t, msg.Content, "Chatbot")
				require.Equal(t, "assistant", msg.Role)
			},
		},
		{
			name: "Missing template file",
			setup: func(t *testing.T, baseDir string) {
				sidebar := filepath.Join(baseDir, ".sidebar")
				require.NoError(t, os.MkdirAll(sidebar, 0755))
				// Note: we do *not* create templates dir → should fail
				configFile := filepath.Join(sidebar, "sidebar-config.yaml")
				require.NoError(t, os.WriteFile(configFile, []byte("API_KEY: fake\n"), 0644))
			},
			mockClient:    &app.MockClient{},
			expectErr:     true,
			expectErrPart: "error copying template",
		},
		{
			name: "Missing .sidebar dir",
			setup: func(t *testing.T, baseDir string) {
				// Create only templates
				require.NoError(t, os.MkdirAll(filepath.Join(baseDir, "templates"), 0755))
				os.WriteFile(filepath.Join(baseDir, "templates/convo.json"), []byte("{}"), 0644)
			},
			mockClient:    &app.MockClient{},
			expectErr:     true,
			expectErrPart: "config directory missing",
		},
		{
			name: "OpenAI JSON error",
			setup: func(t *testing.T, baseDir string) {
				setupFakeSidebarEnv(t)
			},
			mockClient: &app.MockClient{
				Err: fmt.Errorf("OpenAI error: {\"message\":\"invalid model\",\"type\":\"invalid_request_error\"}"),
			},
			expectErr: false, // still returns convo
			verify: func(t *testing.T, convo *app.Conversation) {
				msg := convo.Messages[convo.LastMessageID]
				require.NotNil(t, msg)
				require.Contains(t, msg.ErrorResponse, "invalid model")
			},
		},
		{
			name: "OpenAI plain error",
			setup: func(t *testing.T, baseDir string) {
				setupFakeSidebarEnv(t)
			},
			mockClient: &app.MockClient{
				Err: fmt.Errorf("network timeout"),
			},
			expectErr: false,
			verify: func(t *testing.T, convo *app.Conversation) {
				msg := convo.Messages[convo.LastMessageID]
				require.Empty(t, msg.ErrorResponse)
			},
		},
	}

	// Test the test cases in the table
	for _, testcase := range tests {
		t.Run(testcase.name, func(t *testing.T) {
			baseDir := t.TempDir()
			testcase.setup(t, baseDir)

			// Check the errors we are expecting
			topic := "testing topic for chatbot wrapper application"
			convo, err := app.StartNewConversation(testcase.mockClient, topic)
			if testcase.expectErr {
				require.Error(t, err)
				if testcase.expectErrPart != "" {
					require.Contains(t, err.Error(), testcase.expectErrPart)
				}
				return
			}

			// Make sure the test case is ready to be verified
			// AKA the conversation built correctly
			require.NoError(t, err)
			require.NotNil(t, convo)
			require.NotEmpty(t, convo.ConversationID)
			require.NotEmpty(t, convo.LastMessageID)

			// If there is a verification step, run it
			if testcase.verify != nil {
				testcase.verify(t, convo)
			}
		})
	}
}
