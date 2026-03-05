package main

import (
	"fmt"
	"os"

	"scaffold/internal/app/task"
	tmplservice "scaffold/internal/app/template"
	"scaffold/internal/conf"
	"scaffold/internal/repo/db"
	"scaffold/internal/repo/db/dao"
	"scaffold/internal/repo/db/model"
	"scaffold/internal/repo/redis/cache"

	"scaffold/cmd/server/bootstrap"

	"github.com/spf13/cobra"
)

var (
	serveCmd     bool
	templateFlag string
	outputFlag   string
	configFlag   string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "scaffold",
		Short: "A universal project scaffolding tool",
		Long:  `Scaffold helps you generate projects from templates via web UI or CLI.`,
		Run: func(cmd *cobra.Command, args []string) {
			if serveCmd {
				startServer()
			} else {
				cmd.Help()
			}
		},
	}

	rootCmd.Flags().BoolVarP(&serveCmd, "serve", "s", false, "Start web server")
	rootCmd.Flags().StringVar(&configFlag, "config", "", "Configuration file path")
	rootCmd.Flags().StringP("port", "p", "9090", "Server port")
	rootCmd.Flags().String("host", "0.0.0.0", "Server host")
	rootCmd.Flags().String("cache-dir", "./cache", "Cache directory")
	rootCmd.Flags().Bool("debug", false, "Enable debug mode")
	rootCmd.Flags().String("access-key", "", "Access key for API authentication")
	rootCmd.Flags().Bool("db-enable", false, "Enable database storage")
	rootCmd.Flags().String("db-driver", "sqlite", "Database driver (sqlite/mysql/postgres)")
	rootCmd.Flags().String("db-host", "localhost", "Database host")
	rootCmd.Flags().Int("db-port", 3306, "Database port")
	rootCmd.Flags().String("db-user", "root", "Database user")
	rootCmd.Flags().String("db-password", "", "Database password")
	rootCmd.Flags().String("db-name", "./data/scaffold.db", "Database name or path")

	initCmd := &cobra.Command{
		Use:   "init [project-name]",
		Short: "Initialize a new project interactively",
		Args:  cobra.MaximumNArgs(1),
		Run:   runInit,
	}
	initCmd.Flags().StringVarP(&templateFlag, "template", "t", "", "Template ID to use")
	initCmd.Flags().StringVarP(&outputFlag, "output", "o", ".", "Output directory")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List available templates",
		Run:   runList,
	}

	rootCmd.AddCommand(initCmd, listCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func startServer() {
	// Set config path from CLI flag if provided
	if configFlag != "" {
		os.Setenv("CONFIG_PATH", configFlag)
	}

	h, err := bootstrap.Bootstrap()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start server: %v\n", err)
		os.Exit(1)
	}

	cfg := conf.GlobalConfig
	fmt.Printf("Scaffold server running at http://%s\n", cfg.GetServerAddr())
	h.Spin()
}

func runInit(cmd *cobra.Command, args []string) {
	projectName := ""
	if len(args) > 0 {
		projectName = args[0]
	}

	svc := createTemplateService()
	templates := svc.ListTemplates()
	if len(templates) == 0 {
		fmt.Fprintln(os.Stderr, "No templates available. Add templates first.")
		os.Exit(1)
	}

	if err := runBubbleTeaUI(svc, projectName); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func runList(cmd *cobra.Command, args []string) {
	svc := createTemplateService()
	templates := svc.ListTemplates()

	if len(templates) == 0 {
		fmt.Println("No templates available.")
		return
	}

	fmt.Println("Available templates:")
	for _, t := range templates {
		fmt.Printf("  - %s (%s): %s\n", t.ID, t.Name, t.Description)
	}
}

func createTemplateService() *tmplservice.Service {
	if err := conf.Init(configFlag); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to load config, using defaults: %v\n", err)
	}
	cfg := conf.GlobalConfig

	var opts []tmplservice.ServiceOption

	if cfg.Database.Enable {
		if err := db.Init(&cfg.Database); err == nil {
			db.AutoMigrate(&model.TemplateModel{})
			repo := dao.NewGormTemplateRepository(db.DB)
			opts = append(opts, tmplservice.WithRepository(repo))
		}
	}

	templateCache := cache.NewTemplateCache(cfg.Server.CacheDir, false)
	opts = append(opts, tmplservice.WithCache(templateCache))

	svc := tmplservice.NewService(cfg.Server.CacheDir, opts...)
	svc.LoadFromDatabase()

	// Ensure task service exists for CLI too
	_ = task.NewService()

	return svc
}
