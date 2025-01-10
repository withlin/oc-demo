package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewRootCmd creates a new root command
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skectl",
		Short: "skectl is a command line tool similar to OpenShift CLI",
		Long: `skectl is a command line tool that provides similar functionality to OpenShift CLI (oc).
It supports cluster login and context management.

Available Commands:
  login       Log in to a server
  use-context Switch to a different context

Use "skectl <command> --help" for more information about a command.`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	// Add subcommands
	cmd.AddCommand(loginCmd)
	cmd.AddCommand(useContextCmd)

	return cmd
}

var rootCmd = NewRootCmd()

// Execute executes the root command
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	return nil
} 