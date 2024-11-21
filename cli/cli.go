package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "cache-proxy",
	Short:   "cache-proxy - a simple CLI to start the Caching Proxy server",
	Version: "v1.0.0",
	Long: `cache-proxy is CLI to start a caching server that caches responses from other servers.
One can use proxy to start or stop the server which will cache the responses from the server`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello from cache-proxy")
	},
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error while executing your CLI command '%s'", err)
		os.Exit(1)
		return err
	}
	return nil
}
