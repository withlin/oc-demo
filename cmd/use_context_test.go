package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func TestUseContextCmd(t *testing.T) {
	// Create temporary directory for test
	tmpDir, err := os.MkdirTemp("", "skectl-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Set KUBECONFIG environment variable
	kubeconfigPath := filepath.Join(tmpDir, "config")
	os.Setenv("KUBECONFIG", kubeconfigPath)
	defer os.Unsetenv("KUBECONFIG")

	// Create test kubeconfig
	config := api.NewConfig()
	config.Contexts["context1"] = api.NewContext()
	config.Contexts["context2"] = api.NewContext()
	config.CurrentContext = "context1"

	require.NoError(t, clientcmd.WriteToFile(*config, kubeconfigPath))

	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "switch context success",
			args:        []string{"context2"},
			expectError: false,
		},
		{
			name:        "switch to non-existent context",
			args:        []string{"unknown"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute command
			cmd := NewUseContextCmd()
			cmd.SetArgs(tt.args)
			err := cmd.Execute()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Verify context was switched
				newConfig, err := clientcmd.LoadFromFile(kubeconfigPath)
				require.NoError(t, err)
				assert.Equal(t, tt.args[0], newConfig.CurrentContext)
			}
		})
	}
} 