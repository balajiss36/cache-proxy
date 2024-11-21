package main

import (
	"log"
	"os"

	"github.com/balajiss36/cache-proxy/cli"
)

func main() {
	log.Println("Starting the CLI for cache-proxy")
	err := cli.Execute()
	if err != nil {
		log.Fatalf("Error while executing the CLI: %v", err)
		os.Exit(1)
		return
	}
}
