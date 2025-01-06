package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vaziolabs/lumberjack/dashboard"
	"github.com/vaziolabs/lumberjack/internal"
	"github.com/vaziolabs/lumberjack/types"
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
			user := types.User{}
			if exists := configExists(); !exists {
				asciiArt, err := os.ReadFile("cmd/cli.ascii")
				if err == nil {
					fmt.Println(string(asciiArt))
				}
				createConfig(cmd, args)
			} else {
				startServer(cmd, args, user)
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
	Run: func(cmd *cobra.Command, args []string) {
		user := types.User{}
		startServer(cmd, args, user)
	},
}

func createConfig(cmd *cobra.Command, args []string) {
	config := types.Config{
		Version:   "0.1.1-alpha",
		Databases: make(map[string]types.DBConfig),
	}

	dbConfig := types.DBConfig{
		Domain:        "localhost",
		Port:          "8080",
		DashboardPort: "8081",
		DBName:        "",
	}

	user := types.User{
		Username: "admin",
		Password: "admin",
	}

	prompts := []struct {
		label    string
		field    *string
		default_ string
	}{
		{"LumberJack Database Name", &dbConfig.DBName, "default"},
		{"LumberJack Host Domain", &dbConfig.Domain, "localhost"},
		{"LumberJack API Port", &dbConfig.Port, "8080"},
		{"LumberJack Dashboard Port", &dbConfig.DashboardPort, "8081"},
		{"Admin Username", &user.Username, "admin"},
		{"Admin Password", &user.Password, "admin"},
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

	dbName := "default"
	if len(args) > 0 {
		dbName = args[0]
	}

	dbConfig.DBName = dbName
	config.Databases[dbName] = dbConfig

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(defaultProcDir)

	viper.Set("version", config.Version)
	viper.Set("databases", config.Databases)

	if err := viper.SafeWriteConfig(); err != nil {
		fmt.Printf("Error writing config: %v\n", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(filepath.Join(defaultLibDir, dbName), 0755); err != nil {
		fmt.Printf("Error creating database directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Configuration created successfully!")
}

func startServer(cmd *cobra.Command, args []string, user types.User) {
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

	dbName := "default"
	if len(args) > 0 {
		dbName = args[0]
	}

	dbConfig, exists := config.Databases[dbName]
	if !exists {
		fmt.Printf("Database configuration '%s' not found\n", dbName)
		os.Exit(1)
	}

	// Initialize and start the server
	server := internal.NewServer(dbConfig.Port, user)
	server.Start()

	// Initialize and start dashboard if enabled
	if dashboardSet {
		dashboardServer := dashboard.NewDashboardServer(
			fmt.Sprintf("http://%s:%s", dbConfig.Domain, dbConfig.Port),
			dbConfig.DashboardPort,
		)
		if err := dashboardServer.Start(); err != nil {
			fmt.Printf("Error starting dashboard: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("LumberJack %s server running on %s:%s\n", dbConfig.DBName, dbConfig.Domain, dbConfig.Port)
	if dashboardSet {
		fmt.Printf("LumberJack %s dashboard running on %s:%s\n", dbConfig.DBName, dbConfig.Domain, dbConfig.DashboardPort)
	}

	// Keep the process running
	select {}
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

	var proc types.ProcessInfo
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
		liveFilePath := filepath.Join("/etc/lumberjack/live", p.ID)
		status := "Not Running"
		if _, err := os.Stat(liveFilePath); err == nil {
			status = "Running"
		}

		fmt.Printf("ID: %s | Name: %s | API Port: %s | Dashboard Port: %s | Status: %s\n",
			p.ID, p.DBName, p.APIPort, p.DashboardPort, status)
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

	var config types.Config
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Printf("Error parsing config: %v\n", err)
		return
	}

	// Get running servers to check status
	runningServers, err := getRunningServers()
	if err != nil {
		fmt.Printf("Error getting running servers: %v\n", err)
		return
	}

	fmt.Printf("%-15s %-10s %-15s %-10s\n", "Name", "API Port", "Dashboard Port", "Status")
	fmt.Println(strings.Repeat("-", 55))

	for name, db := range config.Databases {
		status := "Not Running"
		for _, proc := range runningServers {
			if proc.DBName == name {
				status = fmt.Sprintf("Running (%s)", proc.ID)
				break
			}
		}

		fmt.Printf("%-15s %-10s %-15s %-10s\n",
			name,
			db.Port,
			db.DashboardPort,
			status)
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
