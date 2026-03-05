package bootstrap

import (
	"fmt"

	"github.com/cloudwego/hertz/pkg/app/server"

	appcontainer "scaffold/internal/app"
	taskservice "scaffold/internal/app/task"
	tmplservice "scaffold/internal/app/template"
	"scaffold/internal/conf"
	"scaffold/internal/pkg/logger"
	"scaffold/internal/repo/db"
	"scaffold/internal/repo/db/dao"
	"scaffold/internal/repo/db/model"
	redisclient "scaffold/internal/repo/redis"
	"scaffold/internal/repo/redis/cache"

	genrouter "scaffold/gen/http/router"
)

func Bootstrap() (*server.Hertz, error) {
	// Step 1: Init config
	if err := initConfig(); err != nil {
		return nil, fmt.Errorf("init config: %w", err)
	}

	// Step 2: Init logger
	if err := initLogger(); err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}

	// Step 3: Init database (optional)
	if err := initDatabase(); err != nil {
		logger.Warn("Database init failed, running without persistence: " + err.Error())
	}

	// Step 4: Init Redis (optional)
	if err := initRedis(); err != nil {
		logger.Warn("Redis init failed, running without Redis cache: " + err.Error())
	}

	// Step 5: Init services
	initServices()

	// Step 6: Init HTTP server
	h := initServer()

	return h, nil
}

func Cleanup() {
	logger.Info("Cleaning up resources...")
	db.Close()
	redisclient.Close()
	logger.Sync()
}

func initConfig() error {
	return conf.Init()
}

func initLogger() error {
	cfg := conf.GlobalConfig
	return logger.Init(&logger.Config{
		Level:      cfg.Log.Level,
		Filename:   cfg.Log.Filename,
		MaxSize:    cfg.Log.MaxSize,
		MaxBackups: cfg.Log.MaxBackups,
		MaxAge:     cfg.Log.MaxAge,
		Compress:   cfg.Log.Compress,
	})
}

func initDatabase() error {
	cfg := conf.GlobalConfig
	if !cfg.Database.Enable {
		return nil
	}

	if err := db.Init(&cfg.Database); err != nil {
		return err
	}

	if err := db.AutoMigrate(&model.TemplateModel{}); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}

	logger.Info("Database connected successfully")
	return nil
}

func initRedis() error {
	cfg := conf.GlobalConfig
	if !cfg.Redis.Enable {
		return nil
	}

	if err := redisclient.Init(&cfg.Redis); err != nil {
		return err
	}

	logger.Info("Redis connected successfully")
	return nil
}

func initServices() {
	cfg := conf.GlobalConfig

	var tmplOpts []tmplservice.ServiceOption

	if db.DB != nil {
		repo := dao.NewGormTemplateRepository(db.DB)
		tmplOpts = append(tmplOpts, tmplservice.WithRepository(repo))
	}

	templateCache := cache.NewTemplateCache(cfg.Server.CacheDir, cfg.Redis.Enable)
	tmplOpts = append(tmplOpts, tmplservice.WithCache(templateCache))

	tmplSvc := tmplservice.NewService(cfg.Server.CacheDir, tmplOpts...)
	taskSvc := taskservice.NewService()

	if err := tmplSvc.LoadFromDatabase(); err != nil {
		logger.Warn("Failed to load templates from database: " + err.Error())
	}

	appcontainer.Container = &appcontainer.ServiceContainer{
		TemplateService: tmplSvc,
		TaskService:     taskSvc,
	}

	logger.Info("Services initialized successfully")
}

func initServer() *server.Hertz {
	cfg := conf.GlobalConfig
	addr := cfg.GetServerAddr()

	h := server.Default(server.WithHostPorts(addr))

	// Use Hz-generated router registration (routes + middleware defined in gen/)
	genrouter.GeneratedRegister(h)

	logger.Info("HTTP server configured at " + addr)
	return h
}
