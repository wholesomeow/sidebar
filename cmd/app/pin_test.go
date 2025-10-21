package app_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wholesomeow/chatwrapper/cmd/app"
)

// Test pin.go ----------------------------------------------------------------

func TestPinMessage_TestTable(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(t *testing.T, baseDir string)
		mockClient    app.ChatClient
		expectErr     bool
		expectErrPart string
		verify        func(t *testing.T, convo *app.Conversation)
	}{
		{
			name: "Pin Message",
			setup: func(t *testing.T, baseDir string) {
				setupFakeSidebarEnv(t)
			},
			mockClient: &app.MockClient{
				MockResponse: &app.MockCompletion,
			},
			expectErr: false,
			verify: func(t *testing.T, convo *app.Conversation) {
				res := convo.PinMessage(convo.LastMessageID)
				require.True(t, res)
				require.Len(t, convo.Pinned, 1)
			},
		},
		{
			name: "Unpin Message",
			setup: func(t *testing.T, baseDir string) {
				setupFakeSidebarEnv(t)
			},
			mockClient: &app.MockClient{
				MockResponse: &app.MockCompletion,
			},
			expectErr: false,
			verify: func(t *testing.T, convo *app.Conversation) {
				// Pin Message first
				convo.PinMessage(convo.LastMessageID)

				res := convo.UnpinMessage(convo.LastMessageID)
				require.True(t, res)
				require.Len(t, convo.Pinned, 0)
			},
		},
		{
			name: "Pin Message Archived Convo",
			setup: func(t *testing.T, baseDir string) {
				setupFakeSidebarEnv(t)
			},
			mockClient: &app.MockClient{
				MockResponse: &app.MockCompletion,
			},
			expectErr: false,
			verify: func(t *testing.T, convo *app.Conversation) {
				// Archive conversation first
				convo.ArchiveConversation()
				convo.PinMessage(convo.LastMessageID)
				require.Len(t, convo.Pinned, 0)
			},
		},
		{
			name: "Unpin Message Archived Convo",
			setup: func(t *testing.T, baseDir string) {
				setupFakeSidebarEnv(t)
			},
			mockClient: &app.MockClient{
				MockResponse: &app.MockCompletion,
			},
			expectErr: false,
			verify: func(t *testing.T, convo *app.Conversation) {
				// Pin Message first
				convo.PinMessage(convo.LastMessageID)

				// Then archive conversation
				convo.ArchiveConversation()
				convo.UnpinMessage(convo.LastMessageID)
				require.Len(t, convo.Pinned, 1)
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
