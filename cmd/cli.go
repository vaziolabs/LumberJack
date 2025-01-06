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
	deleteAll    bool
	forceDelete  bool
	rootCmd      = &cobra.Command{
		Use:   "lumberjack",
		Short: "LumberJack - Event tracking and management system",
		Long: `LumberJack - Event tracking and management system

Directory Structure:
    Configs:     /etc/lumberjack/config.yaml
    Processes:   /etc/lumberjack/live/
    Logs:        /var/log/lumberjack/
    Data Files:  /var/lib/lumberjack/

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
		Use:   "delete [database-name]",
		Short: "Delete configuration",
		Long: `Delete a database configuration.
Requires database name unless --all flag is provided.
Use --force to skip confirmation prompt.

Example:
    lumberjack delete mydb
    lumberjack delete --all
    lumberjack delete --all --force`,
		Run: deleteConfig,
	}
	killCmd = &cobra.Command{
		Use:   "kill [server-id]",
		Short: "Kill a running server",
		Long: `Kill a running server by its ID.
Use 'list running' to see available server IDs.
Use -d flag to kill only the dashboard.

Example:
    lumberjack kill abc123xyz
    lumberjack kill abc123xyz -d`,
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

	killCmd.Flags().BoolVarP(&dashboardSet, "dashboard", "d", false, "Kill only the dashboard")

	deleteCmd.Flags().BoolVar(&deleteAll, "all", false, "Delete all configurations")
	deleteCmd.Flags().BoolVar(&forceDelete, "force", false, "Force delete without confirmation")
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
		// If this is a spawned process, run the server directly
		if os.Getenv("LUMBERJACK_SPAWNED") == "1" {
			runServer(cmd, args)
			return
		}

		// Otherwise, spawn a new process
		dbName := "default"
		if len(args) > 0 {
			dbName = args[0]
		}

		config := loadConfig(dbName)
		if err := spawnServer(config, dashboardSet); err != nil {
			fmt.Printf("Error spawning server: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("LumberJack %s server started in background\n", dbName)
		if dashboardSet {
			fmt.Printf("Dashboard available at http://%s:%s\n", config.Domain, config.DashboardPort)
		}
	},
}

func runServer(cmd *cobra.Command, args []string) {
	dbName := "default"
	if len(args) > 0 {
		dbName = args[0]
	}

	config := loadConfig(dbName)
	serverConfig := types.ServerConfig{
		Port:    config.Port,
		DBName:  dbName,
		LogPath: defaultLogDir,
		DbPath:  defaultLibDir,
	}

	server, err := internal.NewServer(serverConfig)
	if err != nil {
		fmt.Printf("Error creating server: %v\n", err)
		os.Exit(1)
	}

	if err := server.Start(); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}

	if dashboardSet {
		dash := dashboard.NewDashboard(config.Domain, config.DashboardPort)
		dash.Start()
	}

	select {}
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

	defaultDBName := "default"
	if len(args) > 0 {
		defaultDBName = args[0]
	}

	prompts := []struct {
		label    string
		field    *string
		default_ string
	}{
		{"LumberJack Database Name", &dbConfig.DBName, defaultDBName},
		{"LumberJack Host Domain", &dbConfig.Domain, "localhost"},
		{"LumberJack API Port", &dbConfig.Port, "8080"},
		{"LumberJack Dashboard Port", &dbConfig.DashboardPort, "8081"},
		{"Admin Username", &user.Username, user.Username},
		{"Admin Password", &user.Password, user.Password},
		{"Re-enter Admin Password", &user.Password, user.Password},
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

	// Check for matching passwords
	if prompts[5].default_ != prompts[6].default_ {
		fmt.Println("Passwords do not match")
		os.Exit(1)
	}

	if user.Password != prompts[5].default_ {
		user.Username = prompts[4].default_
		user.Password = prompts[5].default_
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

	// Create initial server to save admin info
	serverConfig := types.ServerConfig{
		Port:    dbConfig.Port,
		DBName:  dbName,
		LogPath: defaultLogDir,
		DbPath:  defaultLibDir,
		User:    user,
	}

	// Initialize server just to save admin info
	_, err := internal.NewServer(serverConfig)
	if err != nil {
		fmt.Printf("Error creating server: %v\n", err)
		os.Exit(1)
	}
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

	if dashboardSet {
		if !proc.DashboardUp {
			fmt.Println("Dashboard is not running for this server")
			return
		}
		// TODO: Implement dashboard process kill
		// For now, just update the process info
		proc.DashboardUp = false
		if err := updateProcessInfo(proc); err != nil {
			fmt.Printf("Error updating process info: %v\n", err)
			return
		}
		fmt.Printf("Dashboard killed for server %s\n", args[0])
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

	fmt.Printf("%-20s %-10s %-20s %-10s\n", "Name", "API Port", "Dashboard Port", "Status")
	fmt.Println(strings.Repeat("-", 60))

	for name, db := range config.Databases {
		apiPort := db.Port
		dashPort := db.DashboardPort
		status := "Not Running"

		for _, proc := range runningServers {
			if proc.DBName == name {
				if proc.DashboardUp {
					dashPort = fmt.Sprintf("%s (Running)", db.DashboardPort)
				}
				status = fmt.Sprintf("Running (%s)", proc.ID)
				break
			}
		}

		fmt.Printf("%-20s %-10s %-20s %-10s\n",
			name,
			apiPort,
			dashPort,
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

func startServer(cmd *cobra.Command, args []string, user types.User) {
	dbName := "default"
	if len(args) > 0 {
		dbName = args[0]
	}

	config := loadConfig(dbName)
	if err := spawnServer(config, dashboardSet); err != nil {
		fmt.Printf("Error spawning server: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("LumberJack %s server started at http://%s:%s\n", dbName, config.Domain, config.Port)
	if dashboardSet {
		fmt.Printf("Dashboard available at http://%s:%s\n", config.Domain, config.DashboardPort)
	}
}
