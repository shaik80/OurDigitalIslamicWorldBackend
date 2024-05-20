/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/shaik80/ODIW/config"
	connect "github.com/shaik80/ODIW/internal/db/opensearch"
	"github.com/shaik80/ODIW/internal/server"

	lp "github.com/shaik80/ODIW/utils/logger"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile   string
	debugMode    bool
	configPrefix = "./config/"
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
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// If configFile is specified, use it; otherwise, use default configuration file
	configFileName := configFile
	if configFileName == "" {
		configFileName = "config.yaml"
	}
	viper.SetConfigFile(configPrefix + configFileName)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config file:", err)
		os.Exit(1)
	}

	// Unmarshal config into struct
	if err := viper.Unmarshal(&config.Cfg); err != nil {
		fmt.Println("Can't unmarshal config:", err)
		os.Exit(1)
	}

	// Set log level from config or default to "info"
	logLevel := strings.ToLower(config.Cfg.Logging.LogLevel)
	enum, ok := map[string]lp.LogLevel{
		"debug": lp.Debug,
		"info":  lp.Info,
		"warn":  lp.Warn,
		"error": lp.Error,
	}[logLevel]
	if !ok {
		fmt.Printf("Invalid log level: %s\n", logLevel)
		os.Exit(1)
	}

	// Set up logger with the configured log level
	lp.Logs = lp.ConfigurableLogger{
		LogLevel: enum,
	}

	err := connect.InitOpenSearchClient(config.Cfg)
	if err != nil {
		fmt.Printf("failed to connect to database", err)
		os.Exit(1)
	}
	server.SetupGofiber()
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
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "d", false, "enable debug mode")
}
