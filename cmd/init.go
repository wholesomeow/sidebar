package cmd

import (
	"context"
	"fmt"
	"os"

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
			fmt.Printf("apiKey Variable: %s\n", apiKey)
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

		fmt.Printf("\n*** New Session initiated ***\n - Conversation ID: %s\n - Seed: %d\n - Topic: %s\n\n", conversationID, seed, topic)

		param := openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(topic),
			},
			Seed:  openai.Int(seed),      // TODO: Implement an int64 seed generator
			Model: openai.ChatModelGPT4o, // TODO: Implement a model parameter in settings
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
	rootCmd.AddCommand(sessionInitCmd)
}
