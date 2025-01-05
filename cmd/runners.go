package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"golang.org/x/exp/rand"
)

var (
	processLock sync.RWMutex
)

func getProcessFilePath(id string) string {
	return filepath.Join("/var/lib/lumberjack", id+".dat")
}

func spawnServer(config DBConfig, withDashboard bool) error {
	// Create unique ID for this instance
	id := generateID()

	// Ensure /var/lib/lumberjack exists for .dat files
	if err := os.MkdirAll("/var/lib/lumberjack", 0755); err != nil {
		return fmt.Errorf("failed to create process directory: %v", err)
	}

	// Create log file in /var/log/lumberjack
	logPath := filepath.Join("/var/log/lumberjack", fmt.Sprintf("lumberjack-%s.log", id))
	if err := os.MkdirAll("/var/log/lumberjack", 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	logFile, err := os.Create(logPath)
	if err != nil {
		return fmt.Errorf("failed to create log file: %v", err)
	}
	defer logFile.Close()

	// Prepare command
	cmd := exec.Command(os.Args[0], "start")
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	// Set process to run in background
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}

	// Save process info to individual .dat file
	proc := ProcessInfo{
		ID:            id,
		APIPort:       config.Port,
		DashboardPort: config.DashboardPort,
		PID:           cmd.Process.Pid,
		DBName:        config.DBName,
	}

	data, err := json.Marshal(proc)
	if err != nil {
		return fmt.Errorf("failed to marshal process info: %v", err)
	}

	processFile := getProcessFilePath(id)
	if err := os.WriteFile(processFile, data, 0644); err != nil {
		return fmt.Errorf("failed to save process info: %v", err)
	}

	return nil
}

func getRunningServers() ([]ProcessInfo, error) {
	processLock.RLock()
	defer processLock.RUnlock()

	files, err := os.ReadDir("/var/lib/lumberjack")
	if err != nil {
		if os.IsNotExist(err) {
			return []ProcessInfo{}, nil
		}
		return nil, err
	}

	var processes []ProcessInfo
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".dat") {
			continue
		}

		data, err := os.ReadFile(filepath.Join("/var/lib/lumberjack", file.Name()))
		if err != nil {
			continue
		}

		var proc ProcessInfo
		if err := json.Unmarshal(data, &proc); err != nil {
			continue
		}

		if processExists(proc.PID) {
			processes = append(processes, proc)
		} else {
			// Clean up .dat file for dead process
			os.Remove(filepath.Join("/var/lib/lumberjack", file.Name()))
		}
	}

	return processes, nil
}

func removeProcess(id string) error {
	processLock.Lock()
	defer processLock.Unlock()

	return os.Remove(getProcessFilePath(id))
}

func processExists(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// Send signal 0 to check if process exists
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

func generateID() string {
	// Generate random 8 character string
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, 8)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
