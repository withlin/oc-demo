package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// InputReader defines the interface for input reading
type InputReader interface {
	// ReadLine reads a line of input
	ReadLine(prompt string) (string, error)
}

// SecureInputReader defines the interface for secure input reading
type SecureInputReader interface {
	// ReadSecurely reads input securely (without displaying the actual input)
	ReadSecurely(prompt string) (string, error)
}

// defaultInputReader implements standard input
type defaultInputReader struct {
	reader *bufio.Reader
}

// terminalInputReader implements terminal input
type terminalInputReader struct {
	defaultInputReader
	fd int
}

// NewInputReader creates a new input reader
func NewInputReader() InputReader {
	return &defaultInputReader{
		reader: bufio.NewReader(os.Stdin),
	}
}

// NewSecureInputReader creates a new secure input reader
func NewSecureInputReader() SecureInputReader {
	return &terminalInputReader{
		defaultInputReader: defaultInputReader{
			reader: bufio.NewReader(os.Stdin),
		},
		fd: int(syscall.Stdin),
	}
}

// ReadLine implements standard input reading
func (r *defaultInputReader) ReadLine(prompt string) (string, error) {
	if _, err := fmt.Fprint(os.Stdout, prompt); err != nil {
		return "", fmt.Errorf("failed to write prompt: %w", err)
	}

	reader := bufio.NewReader(r.reader)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}

	return strings.TrimSpace(input), nil
}

// ReadLine implements terminal input reading
func (r *terminalInputReader) ReadLine(prompt string) (string, error) {
	return r.defaultInputReader.ReadLine(prompt)
}

// ReadSecurely implements secure input reading
func (r *terminalInputReader) ReadSecurely(prompt string) (string, error) {
	if _, err := fmt.Fprint(os.Stdout, prompt); err != nil {
		return "", fmt.Errorf("failed to write prompt: %w", err)
	}

	// Get current terminal state
	oldState, err := terminal.GetState(r.fd)
	if err != nil {
		return "", fmt.Errorf("failed to set terminal to raw mode: %w", err)
	}
	defer func() {
		_ = terminal.Restore(r.fd, oldState)
	}()

	var password []byte
	for {
		var char [1]byte
		n, err := r.reader.Read(char[:])
		if err != nil {
			return "", fmt.Errorf("failed to read character: %w", err)
		}
		if n == 0 {
			continue
		}

		c := char[0]

		switch c {
		case '\r', '\n':
			_, _ = fmt.Fprintln(os.Stdout)
			return string(password), nil
		case 3: // Ctrl+C
			return "", fmt.Errorf("interrupted by user")
		case '\b', 127: // Backspace and Delete
			if len(password) > 0 {
				password = password[:len(password)-1]
				_, _ = fmt.Fprint(os.Stdout, "\b \b")
			}
		default:
			if c >= 32 && c <= 126 { // Printable characters
				password = append(password, c)
				_, _ = fmt.Fprint(os.Stdout, "*")
			}
		}
	}
}

// For backward compatibility, keep the original functions
func ReadPassword(prompt string) (string, error) {
	reader := NewSecureInputReader()
	return reader.ReadSecurely(prompt)
}

func ReadInput(prompt string) (string, error) {
	reader := NewInputReader()
	return reader.ReadLine(prompt)
} 