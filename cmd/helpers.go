package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
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

func configExists() bool {
	_, err := os.Stat(filepath.Join(defaultProcDir, "config.yaml"))
	return err == nil
}

func getConfigDir() string {
	return defaultProcDir
}
