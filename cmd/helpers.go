package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vaziolabs/lumberjack/types"
	"golang.org/x/exp/rand"
)

func generateID() string {
	// Generate random 8 character string
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, 8)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func processExists(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// Send a signal 0 to the process to check if it is running
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

func deleteConfig(cmd *cobra.Command, args []string) {
	if !deleteAll && len(args) == 0 {
		fmt.Println("Error: database name required unless --all flag is provided")
		return
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(defaultProcDir)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		return
	}

	var config types.Config
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Printf("Error parsing config: %v\n", err)
		return
	}

	if !deleteAll {
		dbName := args[0]
		if _, exists := config.Databases[dbName]; !exists {
			fmt.Printf("Database '%s' not found in configuration\n", dbName)
			return
		}

		if !forceDelete {
			fmt.Printf("Are you sure you want to delete database '%s'? [y/N]: ", dbName)
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				fmt.Println("Operation cancelled")
				return
			}
		}

		delete(config.Databases, dbName)
		if err := saveConfig(config); err != nil {
			fmt.Printf("Error saving configuration: %v\n", err)
			return
		}
		fmt.Printf("Database '%s' configuration deleted successfully\n", dbName)
		return
	}

	// Handle --all flag
	if !forceDelete {
		fmt.Print("Are you sure you want to delete ALL configurations? [y/N]: ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Operation cancelled")
			return
		}
	}

	configFile := filepath.Join(defaultProcDir, "config.yaml")
	if err := os.Remove(configFile); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No configuration file found")
			return
		}
		fmt.Printf("Error deleting config: %v\n", err)
		return
	}

	fmt.Println("All configurations deleted successfully")
}

func saveConfig(config types.Config) error {
	viper.Set("version", config.Version)
	viper.Set("databases", config.Databases)
	return viper.WriteConfig()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func configExists() bool {
	_, err := os.Stat(filepath.Join(defaultProcDir, "config.yaml"))
	return err == nil
}

func getConfigDir() string {
	return defaultProcDir
}

func loadConfig(dbName string) types.DBConfig {
	configDir := getConfigDir()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		os.Exit(1)
	}

	var config types.Config
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Printf("Error parsing config: %v\n", err)
		os.Exit(1)
	}

	dbConfig, exists := config.Databases[dbName]
	if !exists {
		fmt.Printf("Database %s not found in config\n", dbName)
		os.Exit(1)
	}

	return dbConfig
}

func updateProcessInfo(proc types.ProcessInfo) error {
	data, err := json.Marshal(proc)
	if err != nil {
		return fmt.Errorf("failed to marshal process info: %v", err)
	}

	processFile := getProcessFilePath(proc.ID)
	return os.WriteFile(processFile, data, 0644)
}
