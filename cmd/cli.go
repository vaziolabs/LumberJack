package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultLogDir  = "/var/log/lumberjack"
	defaultLibDir  = "/var/lib/lumberjack"
	defaultProcDir = "/etc/lumberjack"
)

var (
	configFile   string
	dashboardSet bool
	rootCmd      = &cobra.Command{
		Use:   "lumberjack",
		Short: "LumberJack - Event tracking and management system",
		Long: `LumberJack - Event tracking and management system

Directory Structure:
    Configs:     /etc/lumberjack/config.yaml
    Logs:        /var/log/lumberjack/
    Data Files:  /var/lib/lumberjack/
    Processes:   /etc/lumberjack/

Commands:
    create              Create a new configuration
    start [db-name]     Start server with optional database name
    list running        List all running servers
    list configs        List current configuration
    kill [server-id]    Kill a running server
    logs [server-id]    View server logs
    delete             Delete current configuration

Flags:
    -d, --dashboard    Start with dashboard enabled
    -c, --config       Specify config file to use`,
		Run: func(cmd *cobra.Command, args []string) {
			if exists := configExists(); !exists {
				asciiArt, err := os.ReadFile("cmd/cli.ascii")
				if err == nil {
					fmt.Println(string(asciiArt))
				}
				createConfig(cmd, args)
			} else {
				startServer(cmd, args)
			}
		},
	}
	listCmd = &cobra.Command{
		Use:   "list [running|configs]",
		Short: "List running servers or configurations",
		Long: `List running servers or view current configuration.

Example:
    lumberjack list running
    lumberjack list configs`,
		Run: listHandler,
	}
	deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete current configuration",
		Long: `Delete the current configuration file.

Example:
    lumberjack delete`,
		Run: deleteConfig,
	}
	killCmd = &cobra.Command{
		Use:   "kill [server-id]",
		Short: "Kill a running server",
		Long: `Kill a running server by its ID.
Use 'list running' to see available server IDs.

Example:
    lumberjack kill abc123xyz`,
		Run: killServer,
	}
	logsCmd = &cobra.Command{
		Use:   "logs [server-id]",
		Short: "View server logs",
		Long: `View logs for a specific server by its ID.
Use 'list running' to see available server IDs.

Example:
    lumberjack logs abc123xyz`,
		Run: viewLogs,
	}
)

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(killCmd)
	rootCmd.AddCommand(logsCmd)

	createCmd.AddCommand(newHelpCmd(createCmd))
	startCmd.AddCommand(newHelpCmd(startCmd))
	listCmd.AddCommand(newHelpCmd(listCmd))
	deleteCmd.AddCommand(newHelpCmd(deleteCmd))
	killCmd.AddCommand(newHelpCmd(killCmd))
	logsCmd.AddCommand(newHelpCmd(logsCmd))

	startCmd.Flags().BoolVarP(&dashboardSet, "dashboard", "d", false, "Start with dashboard")
	startCmd.Flags().StringVarP(&configFile, "config", "c", "", "Config file to use")
}

var createCmd = &cobra.Command{
	Use:   "create [config-name]",
	Short: "Create a new LumberJack configuration",
	Long: `Create a new LumberJack configuration file.
If no config name is provided, 'default' will be used.

Example:
    lumberjack create
    lumberjack create myconfig`,
	Run: createConfig,
}

var startCmd = &cobra.Command{
	Use:   "start [database-name]",
	Short: "Start LumberJack server",
	Long: `Start LumberJack server with optional database name.
Use -d flag to start with dashboard enabled.

Example:
    lumberjack start
    lumberjack start mydb
    lumberjack start -d
    lumberjack start mydb -d`,
	Run: startServer,
}

func configExists() bool {
	_, err := os.Stat(filepath.Join(defaultProcDir, "config.yaml"))
	return err == nil
}

func getConfigDir() string {
	return defaultProcDir
}

func createConfig(cmd *cobra.Command, args []string) {
	config := Config{
		Version:   "0.1.1-alpha",
		Databases: make(map[string]DBConfig),
	}

	dbConfig := DBConfig{
		Domain:        "localhost",
		Port:          "8080",
		DashboardPort: "8081",
		DBName:        "",
	}

	user := User{
		Username: "admin",
		Password: "admin",
	}

	prompts := []struct {
		label    string
		field    *string
		default_ string
	}{
		{"LumberJack Database Name", &dbConfig.DBName, ""},
		{"LumberJack Host Domain [localhost]", &dbConfig.Domain, "localhost"},
		{"LumberJack API Port [8080]", &dbConfig.Port, "8080"},
		{"LumberJack Dashboard Port [8081]", &dbConfig.DashboardPort, "8081"},
		{"Admin Username [admin]", &user.Username, "admin"},
		{"Admin Password [admin]", &user.Password, "admin"},
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

	configKey := dbConfig.DBName
	if configKey == "" {
		configKey = "default"
	}

	config.Databases[configKey] = dbConfig

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(defaultProcDir)

	viper.Set("version", config.Version)
	viper.Set("databases", config.Databases)

	if err := viper.SafeWriteConfig(); err != nil {
		fmt.Printf("Error writing config: %v\n", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(filepath.Join(defaultLibDir, configKey), 0755); err != nil {
		fmt.Printf("Error creating database directory: %v\n", err)
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

	dbName := "default"
	if len(args) > 0 {
		dbName = args[0]
	}

	dbConfig, exists := config.Databases[dbName]
	if !exists {
		fmt.Printf("Database configuration '%s' not found\n", dbName)
		os.Exit(1)
	}

	if err := spawnServer(dbConfig, dashboardSet); err != nil {
		fmt.Printf("Error spawning server: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Server started successfully in background")
}

func killServer(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Server ID required")
		return
	}

	processes, err := getRunningServers()
	if err != nil {
		fmt.Printf("Error getting running servers: %v\n", err)
		return
	}

	var proc ProcessInfo
	found := false
	for _, p := range processes {
		if p.ID == args[0] {
			proc = p
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("Server %s not found\n", args[0])
		return
	}

	process, err := os.FindProcess(proc.PID)
	if err != nil {
		fmt.Printf("Error finding process: %v\n", err)
		return
	}

	if err := process.Kill(); err != nil {
		fmt.Printf("Error killing process: %v\n", err)
		return
	}

	if err := removeProcess(args[0]); err != nil {
		fmt.Printf("Error removing process from tracking: %v\n", err)
		return
	}

	fmt.Printf("Server %s killed successfully\n", args[0])
}

func viewLogs(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Server ID required")
		return
	}

	processes, err := getRunningServers()
	if err != nil {
		fmt.Printf("Error getting running servers: %v\n", err)
		return
	}

	var logFile string
	for _, p := range processes {
		if p.ID == args[0] {
			logFile = getProcessFilePath(p.ID)
			break
		}
	}

	if logFile == "" {
		fmt.Printf("Server %s not found\n", args[0])
		return
	}

	logs, err := os.ReadFile(logFile)
	if err != nil {
		fmt.Printf("Error reading log file: %v\n", err)
		return
	}

	fmt.Printf("=== Logs for server %s ===\n", args[0])
	fmt.Println(string(logs))
}

func listHandler(cmd *cobra.Command, args []string) {
	if len(args) == 0 || args[0] == "configs" {
		listConfig(cmd, args)
		return
	}

	if args[0] == "running" {
		listRunning()
		return
	}

	fmt.Println("Invalid argument. Use 'running' or 'configs'")
}

func listRunning() {
	// Get list of running servers from process file
	processes, err := getRunningServers()
	if err != nil {
		fmt.Printf("Error getting running servers: %v\n", err)
		return
	}

	if len(processes) == 0 {
		fmt.Println("No running servers found")
		return
	}

	fmt.Println("Running Servers:")
	for _, p := range processes {
		fmt.Printf("ID: %s | API Port: %s | Dashboard Port: %s\n",
			p.ID, p.APIPort, p.DashboardPort)
	}
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

	// Show detailed config for specific database
	if len(args) > 1 {
		dbConfig, exists := config.Databases[args[1]]
		if !exists {
			fmt.Printf("Database configuration '%s' not found\n", args[1])
			return
		}

		fmt.Printf("\nConfiguration for %s:\n", args[1])
		fmt.Printf("Domain: %s\n", dbConfig.Domain)
		fmt.Printf("API Port: %s\n", dbConfig.Port)
		fmt.Printf("Dashboard Port: %s\n", dbConfig.DashboardPort)
		fmt.Printf("Database Name: %s\n", dbConfig.DBName)
		return
	}

	// List all databases in table format
	fmt.Printf("\nConfigured Databases:\n")
	fmt.Printf("%-20s %-15s %-10s %-15s\n", "ID", "Name", "API Port", "Dashboard Port")
	fmt.Println(strings.Repeat("-", 65))

	for id, db := range config.Databases {
		fmt.Printf("%-20s %-15s %-10s %-15s\n",
			id,
			db.DBName,
			db.Port,
			db.DashboardPort)
	}
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

func newHelpCmd(parent *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "help",
		Short: fmt.Sprintf("Help about the %s command", parent.Name()),
		Long:  parent.Long,
		Run: func(cmd *cobra.Command, args []string) {
			parent.Help()
		},
	}
}
