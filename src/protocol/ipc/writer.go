package ipc

import (
	"fmt"
	"os"
	"sync"
)

type Writer struct {
	Path string
	File *os.File
	mu   sync.Mutex
}

func NewWriter(path string) (*Writer, error) {
	// Open file to write
	file, err := os.OpenFile(path, os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open fifo %s for writing: %w", path, err)
	}

	return &Writer{
		Path: path,
		File: file,
	}, nil
}

func (w *Writer) WriteLine(payload []byte) error {
	// Handle Thread safety
	w.mu.Lock()
	defer w.mu.Unlock()

	// Check line is terminated by a \n
	data := payload
	if len(data) == 0 || data[len(data)-1] != '\n' {
		data = append(append([]byte(nil), data...), '\n')
	}

	_, err := w.File.Write(data)
	return err
}

func (w *Writer) Close() error {
	return w.File.Close()
}
