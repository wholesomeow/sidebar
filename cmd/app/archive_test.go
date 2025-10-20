package app_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wholesomeow/chatwrapper/cmd/app"
)

// Test archive.go ------------------------------------------------------------

// Tests to write
// 1. Archive conversation
//     A. As Owner (success)
//     B. As Non-Owner (failure)
// 2. Attempt to write to conversation (should fail)
// 3. Unarchive conversation
//     A. As Owner (success)
//     B. As Non-Owner (failure)

func TestArchiveConversation_TestTable(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(t *testing.T, baseDir string)
		mockClient    app.ChatClient
		expectErr     bool
		expectErrPart string
		verify        func(t *testing.T, convo *app.Conversation)
	}{
		{
			name: "Archive Conversation",
			setup: func(t *testing.T, baseDir string) {
				setupFakeSidebarEnv(t)
			},
			mockClient: &app.MockClient{
				MockResponse: &app.MockCompletion,
			},
			expectErr: false,
			verify: func(t *testing.T, convo *app.Conversation) {
				convo.ArchiveConversation()
				require.True(t, convo.Archive)
			},
		},
		{
			name: "Unarchive Conversation",
			setup: func(t *testing.T, baseDir string) {
				setupFakeSidebarEnv(t)
			},
			mockClient: &app.MockClient{
				MockResponse: &app.MockCompletion,
			},
			expectErr: false,
			verify: func(t *testing.T, convo *app.Conversation) {
				convo.UnarchiveConversation()
				require.False(t, convo.Archive)
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
