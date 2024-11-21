package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/balajiss36/cache-proxy/cache"
	"github.com/balajiss36/cache-proxy/proxy"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:     "start",
	Short:   "start",
	Aliases: []string{"start"},
	Run: func(cmd *cobra.Command, args []string) {
		portNumber, err := cmd.Flags().GetString("port")
		if err != nil {
			fmt.Println("Error while getting the port number")
		}
		url, err := cmd.Flags().GetString("url")
		if err != nil {
			fmt.Println("Error while getting the url")
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		cache := cache.NewCache()

		proxy := proxy.Proxy{
			Context: ctx,
			Port:    portNumber,
			URL:     url,
			Cache:   cache,
		}

		err = proxy.StartServer()
		if err != nil {
			fmt.Println("Error while starting the server")
			os.Exit(1)
		}
	},
}

func init() {
	startCmd.Flags().StringP("port", "p", "", "Port number to start the server")
	startCmd.Flags().StringP("url", "u", "", "URL of the server to proxy")
	err := startCmd.MarkFlagRequired("port")
	if err != nil {
		fmt.Println("Error while marking the port flag as required")
		return
	}
	err = startCmd.MarkFlagRequired("url")
	if err != nil {
		fmt.Println("Error while marking the url flag as required")
		return
	}
	rootCmd.AddCommand(startCmd)
}
