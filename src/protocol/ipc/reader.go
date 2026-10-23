package ipc

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

type Reader struct {
	Path   string
	File   *os.File
	Reader *bufio.Reader
}

func NewReader(path string) (*Reader, error) {
	// Create FIFO
	err := unix.Mkfifo(path, 0666)
	if err != nil && !os.IsExist(err) {
		return nil, fmt.Errorf("failed to create fifo %s: %w", path, err)
	}

	// Open file to read
	file, err := os.OpenFile(path, os.O_RDWR, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open fifo %s for reading: %w", path, err)
	}

	return &Reader{
		Path:   path,
		File:   file,
		Reader: bufio.NewReader(file),
	}, nil
}

func (r *Reader) ReadLine() ([]byte, error) {
	// Get all the bytes before the first \n
	line, err := r.Reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	// Remove \n
	if len(line) > 0 && line[len(line)-1] == '\n' {
		line = line[:len(line)-1]
	}

	return line, nil
}

func (r *Reader) Close() error {
	// Close and Remove FIFO
	_ = r.File.Close()
	return os.Remove(r.Path)
}
