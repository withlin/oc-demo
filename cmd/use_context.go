package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

// NewUseContextCmd creates a new use-context command
func NewUseContextCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use-context <context>",
		Short: "Switch to a different context",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("context name is required")
			}
			contextName := args[0]

			// Get kubeconfig path
			kubeconfigPath := os.Getenv("KUBECONFIG")
			if kubeconfigPath == "" {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("failed to get home directory: %w", err)
				}
				kubeconfigPath = filepath.Join(homeDir, ".kube", "config")
			}

			// Load kubeconfig
			config, err := clientcmd.LoadFromFile(kubeconfigPath)
			if err != nil {
				return fmt.Errorf("failed to load kubeconfig: %w", err)
			}

			// Check if context exists
			if _, exists := config.Contexts[contextName]; !exists {
				return fmt.Errorf("context %q does not exist", contextName)
			}

			// Switch context
			config.CurrentContext = contextName

			// Save kubeconfig
			if err := clientcmd.WriteToFile(*config, kubeconfigPath); err != nil {
				return fmt.Errorf("failed to write kubeconfig: %w", err)
			}

			fmt.Printf("Switched to context %q\n", contextName)
			return nil
		},
	}

	return cmd
}

var useContextCmd = NewUseContextCmd()

func init() {
	rootCmd.AddCommand(useContextCmd)
} 