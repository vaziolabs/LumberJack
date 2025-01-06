package types

type ProcessInfo struct {
	ID            string `json:"id"`
	APIPort       string `json:"api_port"`
	DashboardPort string `json:"dashboard_port"`
	DashboardUp   bool   `json:"dashboard_up"`
	PID           int    `json:"pid"`
	DBName        string `json:"db_name"`
}

type Config struct {
	Version   string              `yaml:"version"`
	Databases map[string]DBConfig `yaml:"databases"`
}

type DBConfig struct {
	Domain        string `yaml:"domain"`
	Port          string `yaml:"port"`
	DashboardPort string `yaml:"dashboardport"`
	DBName        string `yaml:"dbname"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ServerConfig struct {
	Port    string `json:"port"`
	DBName  string `json:"db_name"`
	LogPath string `json:"log_path"`
	DbPath  string `json:"db_path"`
	User    User   `json:"user"`
}
