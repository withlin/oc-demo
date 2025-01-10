package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/withlin/oc-demo/pkg/auth"
	"github.com/withlin/oc-demo/pkg/util"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

var (
	username string
	password string
	server   string
)

// NewLoginCmd creates a new login command
func NewLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login [flags] <server>",
		Short: "Log in to a server",
		Long:  "Log in to a server using username and password authentication",
		Example: `  # Log in to a server with username
  skectl login https://api.example.com -u admin
  
  # Log in to a server with username and password
  skectl login https://api.example.com -u admin -p password123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("server URL is required")
			}
			server = args[0]

			// Validate server URL
			if server == "" {
				return fmt.Errorf("server URL cannot be empty")
			}

			// Get username if not provided
			if username == "" {
				fmt.Print("Enter username: ")
				var err error
				username, err = util.ReadInput("")
				if err != nil {
					return fmt.Errorf("failed to read username: %w", err)
				}
				if username == "" {
					return fmt.Errorf("username cannot be empty")
				}
			}

			// Get password if not provided
			if password == "" {
				var err error
				password, err = util.ReadPassword("Enter password: ")
				if err != nil {
					return fmt.Errorf("failed to read password: %w", err)
				}
				if password == "" {
					return fmt.Errorf("password cannot be empty")
				}
			}

			// Create authenticator
			config := &auth.Config{
				Server: server,
			}
			authenticator, err := auth.NewAuthenticator(config)
			if err != nil {
				return fmt.Errorf("failed to create authenticator: %w", err)
			}

			// Authenticate user
			authToken, err := authenticator.Authenticate(username, password)
			if err != nil {
				return fmt.Errorf("authentication failed: %w", err)
			}

			// Create kubeconfig
			kubeconfig := api.NewConfig()

			// Create cluster
			cluster := api.NewCluster()
			cluster.Server = server
			cluster.InsecureSkipTLSVerify = true
			kubeconfig.Clusters[server] = cluster

			// Create auth info
			authInfo := api.NewAuthInfo()
			authInfo.Token = authToken
			kubeconfig.AuthInfos[server] = authInfo

			// Create context
			context := api.NewContext()
			context.Cluster = server
			context.AuthInfo = server
			kubeconfig.Contexts[server] = context

			// Set current context
			kubeconfig.CurrentContext = server

			// Get kubeconfig path
			kubeconfigPath := os.Getenv("KUBECONFIG")
			if kubeconfigPath == "" {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("failed to get home directory: %w", err)
				}
				kubeconfigPath = filepath.Join(homeDir, ".kube", "config")
			}

			// Create directory if not exists
			if err := os.MkdirAll(filepath.Dir(kubeconfigPath), 0755); err != nil {
				return fmt.Errorf("failed to create kubeconfig directory: %w", err)
			}

			// Write kubeconfig
			if err := clientcmd.WriteToFile(*kubeconfig, kubeconfigPath); err != nil {
				return fmt.Errorf("failed to write kubeconfig: %w", err)
			}

			fmt.Printf("Successfully logged in as %s to %s\n", username, server)
			return nil
		},
	}

	cmd.Flags().StringVarP(&username, "username", "u", "", "Username for authentication")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Password for authentication")
	cmd.Flags().Bool("insecure-skip-tls-verify", false, "Skip TLS certificate verification")

	return cmd
}

var loginCmd = NewLoginCmd()

func init() {
	rootCmd.AddCommand(loginCmd)
} 