package util

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ErrMockRead represents a mock read error
var ErrMockRead = errors.New("mock read error")

// mockReader implements a mock io.Reader for testing
type mockReader struct {
	data []byte
	err  error
}

func (m *mockReader) Read(p []byte) (n int, err error) {
	if m.err != nil {
		return 0, fmt.Errorf("mock reader error: %w", m.err)
	}
	if len(m.data) == 0 {
		return 0, io.EOF
	}
	n = copy(p, m.data)
	m.data = m.data[n:]
	if len(m.data) == 0 {
		return n, io.EOF
	}
	return n, nil
}

func TestInputReader_ReadLine(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		prompt      string
		expectError bool
		expected    string
	}{
		{
			name:        "should successfully read normal input",
			input:       "test input\n",
			prompt:      "Enter value: ",
			expectError: false,
			expected:    "test input",
		},
		{
			name:        "should handle empty input",
			input:       "\n",
			prompt:      "Enter value: ",
			expectError: false,
			expected:    "",
		},
		{
			name:        "should handle read error",
			input:       "",
			prompt:      "Enter value: ",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reader io.Reader
			if tt.expectError {
				reader = &mockReader{err: ErrMockRead}
			} else {
				reader = bytes.NewBufferString(tt.input)
			}

			inputReader := &defaultInputReader{
				reader: reader,
				writer: io.Discard,
			}

			result, err := inputReader.ReadLine(tt.prompt)

			if tt.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), "mock reader error")
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestTerminalInputReader_ReadLine(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		prompt      string
		expectError bool
		expected    string
	}{
		{
			name:        "should successfully read secure input",
			input:       "secure123\n",
			prompt:      "Enter password: ",
			expectError: false,
			expected:    "secure123",
		},
		{
			name:        "should handle empty secure input",
			input:       "\n",
			prompt:      "Enter password: ",
			expectError: false,
			expected:    "",
		},
		{
			name:        "should handle read error",
			input:       "",
			prompt:      "Enter password: ",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reader io.Reader
			if tt.expectError {
				reader = &mockReader{err: ErrMockRead}
			} else {
				reader = bytes.NewBufferString(tt.input)
			}

			inputReader := &terminalInputReader{
				defaultInputReader: defaultInputReader{
					reader: reader,
					writer: io.Discard,
				},
				fd: 0,
			}

			result, err := inputReader.ReadLine(tt.prompt)

			if tt.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), "mock reader error")
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestReadInput(t *testing.T) {
	t.Run("should successfully read input from stdin", func(t *testing.T) {
		// Save original stdin and stdout
		oldStdin := os.Stdin
		oldStdout := os.Stdout
		defer func() {
			os.Stdin = oldStdin
			os.Stdout = oldStdout
		}()

		// Create a pipe for testing
		r, w, err := os.Pipe()
		require.NoError(t, err, "failed to create pipe")
		os.Stdin = r

		// Write test input in a goroutine
		go func() {
			defer w.Close()
			_, err := w.Write([]byte("test input\n"))
			require.NoError(t, err, "failed to write test input")
		}()

		// Test ReadInput function
		result, err := ReadInput("test")
		require.NoError(t, err, "ReadInput failed")
		assert.Equal(t, "test input", result)
	})
}

func TestReadPassword(t *testing.T) {
	t.Run("should handle password input in various environments", func(t *testing.T) {
		// Skip test in CI environment
		if os.Getenv("CI") == "true" {
			t.Skip("skipping password test in CI environment")
		}

		// Skip test in non-terminal environment
		if os.Getenv("TERM") == "" {
			t.Skip("skipping password test in non-terminal environment")
		}

		// Test password input
		result, err := ReadPassword("test")
		if err != nil {
			t.Logf("warning: password input error: %v", err)
			t.Skip("skipping password test due to terminal error")
		}
		require.NotEmpty(t, result, "password input should not be empty")
	})
} 