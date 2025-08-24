package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"github.com/spf13/cobra"
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

		param := openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(msg),
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
	rootCmd.AddCommand(sessionMsgCmd)
}
