package loader

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// FileInfo contains metadata about a loaded file
type FileInfo struct {
	Path      string
	Name      string
	Extension string
	Size      int64
}

// FileLoader handles loading and reading file contents
type FileLoader struct {
	filePath string
	fileInfo FileInfo
	reader   *bufio.Reader
	file     *os.File
}

// NewFileLoader creates a new file loader instance
func NewFileLoader(path string) (*FileLoader, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	fileInfo, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	if fileInfo.IsDir() {
		return nil, fmt.Errorf("cannot load directory, must be a file")
	}

	fileName := fileInfo.Name()
	extension := strings.ToLower(filepath.Ext(fileName))

	return &FileLoader{
		filePath: absPath,
		fileInfo: FileInfo{
			Path:      absPath,
			Name:      fileName,
			Extension: extension,
			Size:      fileInfo.Size(),
		},
	}, nil
}

// Open opens the file for reading
func (f *FileLoader) Open() error {
	file, err := os.Open(f.filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	f.file = file
	f.reader = bufio.NewReader(file)
	return nil
}

// Close closes the file
func (f *FileLoader) Close() error {
	if f.file != nil {
		return f.file.Close()
	}
	return nil
}

// ReadAll reads the entire file content
func (f *FileLoader) ReadAll() ([]byte, error) {
	if f.file == nil {
		if err := f.Open(); err != nil {
			return nil, err
		}
		defer f.Close()
	}

	return io.ReadAll(f.reader)
}

// ReadLines reads the file line by line using a callback
func (f *FileLoader) ReadLines(callback func(line string) bool) error {
	if f.file == nil {
		if err := f.Open(); err != nil {
			return err
		}
		defer f.Close()
	}

	scanner := bufio.NewScanner(f.file)
	for scanner.Scan() {
		line := scanner.Text()
		if !callback(line) {
			break
		}
	}

	return scanner.Err()
}

// GetFileInfo returns metadata about the loaded file
func (f *FileLoader) GetFileInfo() FileInfo {
	return f.fileInfo
}