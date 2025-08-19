package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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
