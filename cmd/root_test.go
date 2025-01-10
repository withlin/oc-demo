package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/withlin/oc-demo/pkg/testutil"
)

func TestRootCmd(t *testing.T) {
	// Create output capturer
	capture := testutil.NewCaptureOutput()
	require.NoError(t, capture.Start(), "Failed to start output capture")
	defer capture.Stop()

	// Create new command instance
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"--help"})

	// Execute command
	err := cmd.Execute()
	assert.NoError(t, err, "Execute() failed")

	// Get output
	output := capture.Combined()

	// Verify output
	expectedSubstrings := []string{
		"skectl",
		"login",
		"use-context",
		"Use \"skectl <command> --help\" for more information",
	}

	for _, substr := range expectedSubstrings {
		assert.Contains(t, output, substr, "Output does not contain expected text: %s", substr)
	}
}

func TestExecute(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectError bool
		errContains string
	}{
		{
			name:        "help command",
			args:        []string{"--help"},
			expectError: false,
		},
		{
			name:        "unknown command",
			args:        []string{"unknown"},
			expectError: true,
			errContains: "unknown command",
		},
		{
			name:        "version command",
			args:        []string{"version"},
			expectError: true,
			errContains: "unknown command",
		},
		{
			name:        "no command",
			args:        []string{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create output capturer
			capture := testutil.NewCaptureOutput()
			require.NoError(t, capture.Start(), "Failed to start output capture")
			defer capture.Stop()

			// Set command arguments
			rootCmd.SetArgs(tt.args)

			// Execute command
			err := Execute()

			// Get output
			output := capture.Combined()

			// Verify error
			if tt.expectError {
				assert.Error(t, err, "Expected error but got none")
				if err != nil {
					assert.Contains(t, err.Error(), tt.errContains, "Error message does not contain expected text")
				}
			} else {
				assert.NoError(t, err, "Expected no error but got: %v", err)
				assert.NotEmpty(t, output, "Expected output but got none")
			}
		})
	}
} 