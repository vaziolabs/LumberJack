package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/vaziolabs/lumberjack/types"
)

var (
	processLock sync.RWMutex
)

func getProcessFilePath(id string) string {
	return filepath.Join("/var/lib/lumberjack", id+".dat")
}

func spawnServer(config types.DBConfig, withDashboard bool) error {
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

	// Prepare command with proper arguments
	args := []string{"start", config.DBName}
	if withDashboard {
		args = append(args, "-d")
	}
	cmd := exec.Command(os.Args[0], args...)
	cmd.Env = append(os.Environ(), "LUMBERJACK_SPAWNED=1")
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	// Detach process completely from parent
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Pgid:    0, // Force into new process group
	}

	// Ensure process continues running after parent exits
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}

	// Don't wait for the process
	go cmd.Process.Release()

	// Save process info to individual .dat file
	proc := types.ProcessInfo{
		ID:            id,
		APIPort:       config.Port,
		DashboardPort: config.DashboardPort,
		PID:           cmd.Process.Pid,
		DBName:        config.DBName,
		DashboardUp:   withDashboard,
	}

	data, err := json.Marshal(proc)
	if err != nil {
		return fmt.Errorf("failed to marshal process info: %v", err)
	}

	processFile := getProcessFilePath(id)
	if err := os.WriteFile(processFile, data, 0644); err != nil {
		return fmt.Errorf("failed to save process info: %v", err)
	}

	// Create live directory if it doesn't exist
	if err := os.MkdirAll("/etc/lumberjack/live", 0755); err != nil {
		return fmt.Errorf("failed to create live directory: %v", err)
	}

	// Create a file in /etc/lumberjack/live to track running process
	liveFilePath := filepath.Join("/etc/lumberjack/live", id)
	if err := os.WriteFile(liveFilePath, []byte{}, 0644); err != nil {
		return fmt.Errorf("failed to create live process file: %v", err)
	}

	return nil
}

func getRunningServers() ([]types.ProcessInfo, error) {
	processLock.RLock()
	defer processLock.RUnlock()

	// Check live processes directory
	liveFiles, err := os.ReadDir("/etc/lumberjack/live")
	if err != nil {
		if os.IsNotExist(err) {
			return []types.ProcessInfo{}, nil
		}
		return nil, err
	}

	var processes []types.ProcessInfo
	for _, file := range liveFiles {
		// Get process info from .dat file
		data, err := os.ReadFile(filepath.Join("/var/lib/lumberjack", file.Name()+".dat"))
		if err != nil {
			continue
		}

		var proc types.ProcessInfo
		if err := json.Unmarshal(data, &proc); err != nil {
			continue
		}

		processes = append(processes, proc)
	}

	return processes, nil
}

func removeProcess(id string) error {
	processLock.Lock()
	defer processLock.Unlock()

	// Remove the .dat file
	if err := os.Remove(getProcessFilePath(id)); err != nil {
		return err
	}

	// Remove the live process file
	liveFilePath := filepath.Join("/etc/lumberjack/live", id)
	return os.Remove(liveFilePath)
}
