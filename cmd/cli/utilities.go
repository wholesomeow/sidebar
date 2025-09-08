package cmd

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/openai/openai-go/v2"
	"gopkg.in/yaml.v2"
)

// CopyFile copies a file from src to dst.
func CopyFile(src, dstDir string) error {
	// Open the source file for reading
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close() // Ensure the source file is closed

	// Get the base name of the source file
	fileName := filepath.Base(src)
	// Construct the full path for the destination file
	dstPath := filepath.Join(dstDir, fileName)

	// Create the destination file in the target directory
	// 0644 grants read/write for owner, read-only for others
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close() // Ensure the destination file is closed

	// Copy the contents from the source to the destination
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	// Optionally, sync the destination file to ensure data is written to disk
	err = dstFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync destination file: %w", err)
	}

	return nil
}

// Appropriated from github.com/google/uuid
func encodeHex(dst []byte, uuid [16]byte) {
	hex.Encode(dst, uuid[:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], uuid[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], uuid[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], uuid[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:], uuid[10:])
}

// Appropriated from github.com/google/uuid
func CreateUUIDv4() (string, error) {
	uuid := [16]byte{}
	rand_reader := rand.Reader

	// Reads in uuid len amount of bytes of random numbers
	_, err := io.ReadFull(rand_reader, uuid[:])
	if err != nil {
		return "", err
	}

	// Changes the bits specified for the Version & Variant
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x3f // Variant 10

	// Creates a buffer and takes the bytes from the uuid and
	// encodes them into hexadecimal, replacing the appropriate
	// char with "-" to output the uuid as a string
	var buf [36]byte
	encodeHex(buf[:], uuid)

	return string(buf[:]), nil
}

func ReadJSON(path string) ([]byte, error) {
	var nil_data []byte

	// Open JSON File
	f, err := os.Open(path)
	if err != nil {
		return nil_data, err
	}
	defer f.Close()

	// Read JSON File
	byte_value, err := io.ReadAll(f)
	if err != nil {
		return nil_data, err
	}

	return byte_value, nil
}

// WriteJSON writes JSON data to a file
func WriteJSON(path string, data []byte) error {
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("error writing JSON to file: %w", err)
	}
	return nil
}

func parseJSONToSlice(data interface{}) [][]string {
	var result [][]string

	// Check the type of the data and handle accordingly
	switch v := data.(type) {
	case map[string]interface{}:
		// If it's a map, traverse the keys and values recursively
		for _, value := range v {
			// Recursively process each value in the map
			result = append(result, parseJSONToSlice(value)...)
		}

	case []interface{}:
		// If it's a slice, iterate through the elements and process each one
		for _, item := range v {
			// Recursively process each item in the slice
			result = append(result, parseJSONToSlice(item)...)
		}

	case string:
		// If it's a string, add it as a single-element slice
		result = append(result, []string{v})

	case float64:
		// If it's a number, convert it to string and add it as a single-element slice
		result = append(result, []string{fmt.Sprintf("%f", v)})

	case bool:
		// If it's a boolean, convert it to string and add it as a single-element slice
		result = append(result, []string{fmt.Sprintf("%t", v)})

	default:
		// For any other data type, handle it by converting it to string
		result = append(result, []string{fmt.Sprintf("%v", v)})
	}

	return result
}

func JSONToSlice(data []byte) ([][]string, error) {
	// Use a generic map to unmarshal the JSON data
	var temp interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	// Parse the data into a nested slice of strings
	return parseJSONToSlice(data), nil
}

type Message struct {
	Timestamp string                                   `json:"timestamp"`
	MessageID string                                   `json:"messageID"`
	Content   string                                   `json:"content"`
	Param     []openai.ChatCompletionMessageParamUnion `json:"param"`
}

type Conversation struct {
	ConversationID string    `json:"conversationID"`
	Seed           int       `json:"seed"`
	Topic          string    `json:"topic"`
	Messages       []Message `json:"messages"`
}

type ConversationHistory struct {
	ConversationID string `json:"conversationID"`
	Seed           int    `json:"seed"`
	Topic          string `json:"topic"`
	Messages       []struct {
		Timestamp string                                   `json:"timestamp,omitempty"`
		MessageID string                                   `json:"messageID,omitempty"`
		Content   string                                   `json:"content,omitempty"`
		Param     []openai.ChatCompletionMessageParamUnion `json:"param,omitempty"`
	} `json:"messages"`
}

// JSONToStruct takes raw JSON bytes and unmarshals into a Conversation struct
func JSONToStruct(data []byte) (*Conversation, error) {
	var conv *Conversation
	if err := json.Unmarshal(data, &conv); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON to struct: %w", err)
	}
	return conv, nil
}

// StructToJSON marshals a Conversation struct back to pretty-printed JSON
func StructToJSON(conv Conversation) ([]byte, error) {
	data, err := json.MarshalIndent(conv, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshaling struct to JSON: %w", err)
	}
	return data, nil
}

type OpenAIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param,omitempty"`
	Code    string `json:"code"`
}

func CommitCoversation(conv *Conversation, path string) error {
	writeData, err := StructToJSON(*conv)
	if err != nil {
		return err
	}
	err = WriteJSON(path, writeData)
	if err != nil {
		return err
	}

	return nil
}

func UpdateConversationID(yamlFile string, conversationID string) error {
	f, err := os.ReadFile(yamlFile)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", yamlFile, err)
	}

	var config Config
	err = yaml.Unmarshal([]byte(f), &config)
	if err != nil {
		fmt.Printf("Error unmarshaling YAML: %v\n", err)
	}
	config.LastConversationID = conversationID

	updatedYAML, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}
	err = os.WriteFile(yamlFile, updatedYAML, 0644)
	if err != nil {
		return err
	}

	return nil
}
