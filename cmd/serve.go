/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"shaik80/ODIW/config"
	connect "shaik80/ODIW/internal/db/opensearch"
	"shaik80/ODIW/internal/server"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Long: `Starts the application server. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: ServeFunc,
}

func ServeFunc(cmd *cobra.Command, args []string) {
	// Load configuration
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Get configuration values
	cfg := config.GetConfig()
	addr := cfg.Server.Host
	port := cfg.Server.Port

	// Initialize OpenSearch client
	if err := connect.InitOpenSearchClient(); err != nil {
		log.Fatalf("Failed to create OpenSearch client: %v", err)
	}

	// Initialize OpenSearch in controllers
	// controllers.InitOpenSearch(connect.Client)

	// Initialize the server
	app := server.New()
	log.Printf("Server is running on %s:%s", addr, port)
	if err := app.Listen(fmt.Sprintf("%s:%s", addr, port)); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
