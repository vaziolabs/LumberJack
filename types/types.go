package types

import "time"

type ProcessInfo struct {
	ID            string `json:"id"`
	APIPort       string `json:"api_port"`
	DashboardPort string `json:"dashboard_port"`
	DashboardUp   bool   `json:"dashboard_up"`
	PID           int    `json:"pid"`
	DbName        string `json:"db_name"`
	LogDirectory  string `json:"log_directory"`
}

type Config struct {
	Version   string              `yaml:"version"`
	Databases map[string]DBConfig `yaml:"databases"`
}

type DBConfig struct {
	Domain        string `yaml:"domain"`
	Port          string `yaml:"port"`
	DashboardPort string `yaml:"dashboardport"`
	DbName        string `yaml:"dbname"`
	LogDirectory  string `yaml:"logdirectory"`
}

type User struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Organization string `json:"organization"`
}

type ServerConfig struct {
	DatabaseName  string      `json:"database_name"`
	DatabasePath  string      `json:"database_path"`
	LogDirectory  string      `json:"log_directory"`
	Organization  string      `json:"organization"`
	ServerURL     string      `json:"server_url"`
	ServerPort    string      `json:"server_port"`
	DashboardURL  string      `json:"dashboard_url,omitempty"`
	DashboardPort string      `json:"dashboard_port"`
	Phone         string      `json:"phone,omitempty"`
	ProcessInfo   ProcessInfo `json:"process_info"`
}

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Trace     string    `json:"trace,omitempty"`
}
