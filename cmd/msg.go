package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var sessionMsgCmd = &cobra.Command{
	Use:  "msg",
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
		msg := args[0]

		// Get params from last conversation
		directoryToSearch := "./.sidebar"
		configFilePath := fmt.Sprintf("%s/sidebar-config.yaml", directoryToSearch)
		f, err := os.ReadFile(configFilePath)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", configFilePath, err)
		}

		var config Config
		err = yaml.Unmarshal([]byte(f), &config)
		if err != nil {
			fmt.Printf("Error unmarshaling YAML: %v\n", err)
		}
		fileToFind := fmt.Sprintf("convo_%s.json", config.LastConversationID)

		// Find last conversation file in directory
		dir, err := os.Open(directoryToSearch)
		if err != nil {
			fmt.Printf("Error opening directory: %v\n", err)
			return
		}
		defer dir.Close()

		fileInfos, err := dir.Readdir(-1) // -1 means read all entries
		if err != nil {
			fmt.Printf("Error reading directory: %v\n", err)
			return
		}

		found := false
		for _, fileInfo := range fileInfos {
			if !fileInfo.IsDir() && fileInfo.Name() == fileToFind {
				fmt.Printf("File '%s' found in directory '%s'\n", fileToFind, directoryToSearch)
				found = true
				break
			}
		}

		if !found {
			fmt.Printf("File '%s' not found in directory '%s'\n", fileToFind, directoryToSearch)
		}

		// Unmarshal json for params
		prevConvo, err := ReadJSON(fileToFind)
		if err != nil {
			fmt.Printf("Error reading previous conversation data: %v\n", err)
			return
		}
		var convoHistory ConversationHistory
		if err := json.Unmarshal([]byte(prevConvo), &convoHistory); err != nil {
			panic(err)
		}

		prevParams := openai.ChatCompletionNewParams{
			Messages: convoHistory.Messages[0].Param,
		}

		// Append user message
		prevParams.Messages = append(prevParams.Messages, openai.UserMessage(msg))

		// Pass params into new param
		param := openai.ChatCompletionNewParams{
			Messages: prevParams.Messages,
			Seed:     openai.Int(seed),      // TODO: Implement an int64 seed generator
			Model:    openai.ChatModelGPT4o, // TODO: Implement a model parameter in settings
		}

		// TODO: Check out "Responses"
		completion, err := client.Chat.Completions.New(context.Background(), param)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in Response: %s\n", err)
			os.Exit(1)
		}

		param.Messages = append(param.Messages, completion.Choices[0].Message.ToParam())

		fmt.Fprintf(os.Stdout, "Assistant: %s\n", completion.Choices[0].Message.Content)
	},
}

func init() {
	rootCmd.AddCommand(sessionMsgCmd)
}
