package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var (
	ErrFileNotFound = errors.New("file not found in form")
	ErrFileRead     = errors.New("error reading file")
)

// ReadFileFromForm reads a file from the multipart form and returns its contents as a slice of strings.
// Each line in the file becomes a separate string in the slice.
// Empty lines are skipped.
func ReadFileFromForm(r *http.Request, fieldName string) ([]string, error) {
	// Parse the multipart form first
	err := r.ParseMultipartForm(10 << 20) // 10 MB max file size
	if err != nil {
		return nil, fmt.Errorf("failed to parse form: %w", err)
	}

	// Get the uploaded file
	file, _, err := r.FormFile(fieldName)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFileNotFound, err)
	}
	defer file.Close()

	// Read lines from the file
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFileRead, err)
	}

	return lines, nil
}

// ReadFileFromFormWithReader is an alternative function that takes an io.Reader instead of *http.Request.
// This is useful for testing or when you already have the file content.
func ReadFileFromFormWithReader(reader io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFileRead, err)
	}

	return lines, nil
}
