package cmd

type ProcessInfo struct {
	ID            string `json:"id"`
	APIPort       string `json:"api_port"`
	DashboardPort string `json:"dashboard_port"`
	PID           int    `json:"pid"`
	DBName        string `json:"db_name"`
}

type Config struct {
	Version   string              `mapstructure:"version"`
	Databases map[string]DBConfig `mapstructure:"databases"`
}

type DBConfig struct {
	Domain        string `mapstructure:"domain"`
	Port          string `mapstructure:"port"`
	DashboardPort string `mapstructure:"dashboard_port"`
	DBName        string `mapstructure:"db_name"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
