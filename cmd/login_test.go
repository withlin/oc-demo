package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoginCmd(t *testing.T) {
	// Create temporary directory for test
	tmpDir, err := os.MkdirTemp("", "skectl-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Set KUBECONFIG environment variable
	kubeconfigPath := filepath.Join(tmpDir, "config")
	os.Setenv("KUBECONFIG", kubeconfigPath)
	defer os.Unsetenv("KUBECONFIG")

	tests := []struct {
		name        string
		args        []string
		token       string
		expectError bool
	}{
		{
			name:        "login success with token",
			args:        []string{"--token", "test-token", "https://api.test.com:6443"},
			expectError: false,
		},
		{
			name:        "login failed without auth",
			args:        []string{"https://api.test.com:6443"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags
			token = ""
			username = ""
			password = ""
			server = ""

			// Delete existing kubeconfig file
			_ = os.Remove(kubeconfigPath)

			// Execute command
			cmd := NewLoginCmd()
			cmd.SetArgs(tt.args)
			err := cmd.Execute()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Verify kubeconfig was created
				_, err := os.Stat(kubeconfigPath)
				assert.False(t, os.IsNotExist(err))
			}
		})
	}
} 