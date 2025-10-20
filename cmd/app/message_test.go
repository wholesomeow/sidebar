package app_test

// Test Helpers ---------------------------------------------------------------

// Test message.go ------------------------------------------------------------

// Pre-Test Steps
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
