package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vaziolabs/LumberJack/dashboard"
	"github.com/vaziolabs/LumberJack/internal"
)

type Config struct {
	Domain        string `mapstructure:"domain"`
	Port          string `mapstructure:"port"`
	DashboardPort string `mapstructure:"dashboard_port"`
	AdminUser     string `mapstructure:"admin_user"`
	AdminPass     string `mapstructure:"admin_pass"`
	DBName        string `mapstructure:"db_name"`
	StorageDir    string `mapstructure:"storage_dir"`
}

var (
	configFile   string
	dashboardSet bool
	rootCmd      = &cobra.Command{
		Use:   "lumberjack",
		Short: "LumberJack - Event tracking and management system",
		Run: func(cmd *cobra.Command, args []string) {
			if exists := configExists(); !exists {
				createConfig(cmd, args)
			} else {
				startServer(cmd, args)
			}
		},
	}
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List current configuration",
		Run:   listConfig,
	}
	deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete current configuration",
		Run:   deleteConfig,
	}
)

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(deleteCmd)

	startCmd.Flags().BoolVarP(&dashboardSet, "dashboard", "d", false, "Start with dashboard")
	startCmd.Flags().StringVarP(&configFile, "config", "c", "", "Config file to use")
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new LumberJack configuration",
	Run:   createConfig,
}

var startCmd = &cobra.Command{
	Use:   "start [database-name]",
	Short: "Start LumberJack server",
	Run:   startServer,
}

func configExists() bool {
	configDir := getConfigDir()
	_, err := os.Stat(filepath.Join(configDir, "config.yaml"))
	return err == nil
}

func getConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		os.Exit(1)
	}
	return filepath.Join(home, ".config", "lumberjack")
}

func createConfig(cmd *cobra.Command, args []string) {
	config := Config{
		Domain:        "localhost",
		Port:          "8080",
		DashboardPort: "8081",
		StorageDir:    ".",
	}

	prompts := []struct {
		label    string
		field    *string
		default_ string
	}{
		{"LumberJack Domain [localhost]", &config.Domain, "localhost"},
		{"LumberJack API Port [8080]", &config.Port, "8080"},
		{"LumberJack Dashboard Port [8081]", &config.DashboardPort, "8081"},
		{"LumberJack Admin Username", &config.AdminUser, ""},
		{"LumberJack Admin Password", &config.AdminPass, ""},
		{"LumberJack Database Name", &config.DBName, ""},
		{"LumberJack Storage Directory [.]", &config.StorageDir, "."},
	}

	for _, p := range prompts {
		prompt := promptui.Prompt{
			Label:   p.label,
			Default: p.default_,
		}
		if strings.Contains(strings.ToLower(p.label), "password") {
			prompt.Mask = '*'
		}
		result, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed: %v\n", err)
			os.Exit(1)
		}
		*p.field = result
	}

	configDir := getConfigDir()
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Printf("Error creating config directory: %v\n", err)
		os.Exit(1)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	viper.Set("domain", config.Domain)
	viper.Set("port", config.Port)
	viper.Set("dashboard_port", config.DashboardPort)
	viper.Set("admin_user", config.AdminUser)
	viper.Set("admin_pass", config.AdminPass)
	viper.Set("db_name", config.DBName)
	viper.Set("storage_dir", config.StorageDir)

	if err := viper.SafeWriteConfig(); err != nil {
		fmt.Printf("Error writing config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Configuration created successfully!")
}

func startServer(cmd *cobra.Command, args []string) {
	configDir := getConfigDir()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		os.Exit(1)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Printf("Error parsing config: %v\n", err)
		os.Exit(1)
	}

	// If database name provided as argument, use it
	if len(args) > 0 {
		config.DBName = args[0]
	}

	errChan := make(chan error, 2)

	apiServer := internal.NewServer(config.Port)
	apiServer.Start()

	var dashboardServer *dashboard.DashboardServer

	if dashboardSet {
		dashboardServer = dashboard.NewDashboardServer(
			fmt.Sprintf("http://%s:%s", config.Domain, config.Port),
			config.DashboardPort,
		)
		go func() {
			if err := dashboardServer.Start(); err != nil {
				errChan <- err
			}
		}()
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errChan:
		log.Printf("Server error: %v", err)
	case <-quit:
		log.Printf("Shutting down servers...")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("API server forced to shutdown: %v", err)
	}

	if dashboardSet {
		if err := dashboardServer.Shutdown(ctx); err != nil {
			log.Printf("Dashboard server forced to shutdown: %v", err)
		}
	}

	log.Printf("Servers exited properly")
	os.Exit(0)
}

func listConfig(cmd *cobra.Command, args []string) {
	configDir := getConfigDir()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("No configuration found: %v\n", err)
		return
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Printf("Error parsing config: %v\n", err)
		return
	}

	fmt.Printf("Current Configuration:\n")
	fmt.Printf("Domain: %s\n", config.Domain)
	fmt.Printf("API Port: %s\n", config.Port)
	fmt.Printf("Dashboard Port: %s\n", config.DashboardPort)
	fmt.Printf("Admin User: %s\n", config.AdminUser)
	fmt.Printf("Database Name: %s\n", config.DBName)
	fmt.Printf("Storage Directory: %s\n", config.StorageDir)
}

func deleteConfig(cmd *cobra.Command, args []string) {
	configDir := getConfigDir()
	configFile := filepath.Join(configDir, "config.yaml")

	if err := os.Remove(configFile); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No configuration file found.")
			return
		}
		fmt.Printf("Error deleting config: %v\n", err)
		return
	}

	fmt.Println("Configuration deleted successfully!")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
