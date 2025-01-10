package types

import "github.com/vaziolabs/lumberjack/internal/core"

type ProcessInfo struct {
	ID            string `json:"id"`
	PID           int    `json:"pid"`
	ServerURL     string `json:"server_url"`
	ServerPort    string `json:"api_port"`
	DashboardPort string `json:"dashboard_port"`
	DashboardURL  string `json:"dashboard_url,omitempty"`
	DashboardUp   bool   `json:"dashboard_up"`
	DbName        string `json:"db_name"`
	LogPath       string `json:"log_path"`
	DatabasePath  string `json:"database_path"`
}

type Config struct {
	Version   string                 `yaml:"version"`
	Databases map[string]ProcessInfo `yaml:"databases"`
}

// User represents a user in the system
type User struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Username     string            `json:"username"`
	Email        string            `json:"email"`
	Password     string            `json:"password"`
	Organization string            `json:"organization"`
	Phone        string            `json:"phone"`
	Permissions  []core.Permission `json:"permissions"`
}

type ServerConfig struct {
	DatabaseName string      `json:"database_name"`
	Organization string      `json:"organization"`
	Phone        string      `json:"phone,omitempty"`
	Process      ProcessInfo `json:"process"`
}
