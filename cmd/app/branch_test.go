package app_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wholesomeow/chatwrapper/cmd/app"
)

// Test Helpers ---------------------------------------------------------------

// Test branch.go -------------------------------------------------------------

// Tests to write
// 3. Message in branch
// 4. Change branch and message to that branch (Should fail)
// 5. Archive conversation and create branch (Should fail)

func TestBranchConversation_TestTable(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(t *testing.T, baseDir string)
		mockClient    app.ChatClient
		expectErr     bool
		expectErrPart string
		verify        func(t *testing.T, convo *app.Conversation)
	}{
		{
			name: "Branch Created",
			setup: func(t *testing.T, baseDir string) {
				setupFakeSidebarEnv(t)
			},
			mockClient: &app.MockClient{
				MockResponse: &app.MockCompletion,
			},
			expectErr: false,
			verify: func(t *testing.T, convo *app.Conversation) {
				branchID := convo.Branch("test", convo.LastMessageID, false)
				require.Len(t, convo.Branches, 2)

				branch := convo.Branches[branchID]
				require.Contains(t, branch.Name, "test")
				require.False(t, branch.Main)
			},
		},
		{
			name: "Archive Branch Created",
			setup: func(t *testing.T, baseDir string) {
				setupFakeSidebarEnv(t)
			},
			mockClient: &app.MockClient{
				MockResponse: &app.MockCompletion,
			},
			expectErr: false,
			verify: func(t *testing.T, convo *app.Conversation) {
				convo.ArchiveConversation()
				convo.Branch("test", convo.LastMessageID, false)
				require.Len(t, convo.Branches, 1)
			},
		},
		{
			name: "Delete Branch",
			setup: func(t *testing.T, baseDir string) {
				setupFakeSidebarEnv(t)
			},
			mockClient: &app.MockClient{
				MockResponse: &app.MockCompletion,
			},
			expectErr: false,
			verify: func(t *testing.T, convo *app.Conversation) {
				branchID := convo.Branch("test", convo.LastMessageID, false)
				require.Len(t, convo.Branches, 2)

				branch := convo.Branches[branchID]
				require.Contains(t, branch.Name, "test")
				require.False(t, branch.Main)

				res := convo.DeleteBranch(branchID)
				require.True(t, res)
				require.Len(t, convo.Branches, 1)
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
