package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"github.com/spf13/cobra"
)

var sessionInitCmd = &cobra.Command{
	Use:  "init",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize the OpenAI client
		apiKey := os.Getenv("OPENAI_API_KEY")
		if len(apiKey) == 0 {
			apiKey = GetAPIKey()
		}

		client := openai.NewClient(
			option.WithAPIKey(apiKey),
		)

		var seed int64 = 1
		topic := args[0]

		conversationID, err := CreateUUIDv4()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating converstaionID: %s\n", err)
		}

		// Create new conversation file
		path := "./.sidebar"
		_, err = os.Stat(path)
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Copy convo.json template into directory
		fileName := fmt.Sprintf("convo_%s.json", conversationID)
		sourceFilePath := "templates/convo.json"
		err = CopyFile(sourceFilePath, path)
		if err != nil {
			fmt.Printf("Error copying file: %v\n", err)
		}
		oldPath := fmt.Sprintf("%s/convo.json", path)
		newPath := fmt.Sprintf("%s/%s", path, fileName)
		err = os.Rename(oldPath, newPath)
		if err != nil {
			fmt.Printf("Error renaming file: %v\n", err)
		}

		// Read conversation file
		jsonPath := fmt.Sprintf("%s/%s", path, fileName)
		data, err := ReadJSON(jsonPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		convo, err := JSONToStruct(data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Intentionally cursed formatting
		fmt.Printf(`
	                *** New Session initiated ***
*** To respond back to the assistant, use 'sidebar msg "your response"' ***

    - Conversation ID: %s
    - Seed: %d
    - Topic: %s

`,
			conversationID, seed, topic)

		param := openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(topic),
			},
			Seed:  openai.Int(seed),      // TODO: Implement an int64 seed generator
			Model: openai.ChatModelGPT4o, // TODO: Implement a model parameter in settings
		}

		// Write data from the API to file
		convo.ConversationID = conversationID
		convo.Seed = int(seed)
		convo.Topic = topic

		convo.Messages[0].MessageID, err = CreateUUIDv4() // TODO: Double-check that this isn't something returned from OpenAI already (probably is)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating converstaionID: %s\n", err)
		}

		// TODO: Check out "Responses"
		completion, err := client.Chat.Completions.New(context.Background(), param)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in Response: %s\n", err)
			errString := err.Error()

			// Extract only the JSON part
			idx := strings.Index(errString, "{")
			if idx == -1 {
				panic("no JSON found in error string")
			}
			jsonPart := errString[idx:]

			// Parse into struct
			var errResp OpenAIError
			if err := json.Unmarshal([]byte(jsonPart), &errResp); err != nil {
				panic(err)
			}

			// Write data from the API to file
			convo.Messages[0].Content = errResp.Message
		} else {
			// Write data from the API to file
			// convo.Messages[0].Timestamp = timestamp logic lol
			convo.Messages[0].Content = completion.Choices[0].Message.Content
			convo.Messages[0].Param = append(convo.Messages[0].Param, completion.Choices[0].Message.ToParam())

			// Print output to terminal
			fmt.Fprintf(os.Stdout, "Assistant: %s\n", completion.Choices[0].Message.Content)
		}

		// Initial file commit
		writeData, err := StructToJSON(*convo)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing struct to byte array: %s\n", err)
		}
		err = WriteJSON(newPath, writeData)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing byte array to file: %s\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(sessionInitCmd)
}
