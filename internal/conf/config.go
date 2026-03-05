package conf

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig           `mapstructure:"server"`
	Database DatabaseConfig         `mapstructure:"database"`
	Redis    RedisConfig            `mapstructure:"redis"`
	Security SecurityConfig         `mapstructure:"security"`
	Log      LogConfig              `mapstructure:"log"`
	Template TemplateConfigSettings `mapstructure:"template"`
	App      AppMetaConfig          `mapstructure:"app"`
}

type ServerConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	CacheDir     string `mapstructure:"cache_dir"`
	Debug        bool   `mapstructure:"debug"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
	Enable          bool   `mapstructure:"enable"`
	Driver          string `mapstructure:"driver"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"db_name"`
	SSLMode         string `mapstructure:"ssl_mode"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

func (c *DatabaseConfig) DSN() string {
	switch c.Driver {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			c.User, c.Password, c.Host, c.Port, c.DBName)
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
	case "sqlite":
		return c.DBName
	default:
		return ""
	}
}

type RedisConfig struct {
	Enable   bool   `mapstructure:"enable"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type SecurityConfig struct {
	AccessKey string `mapstructure:"access_key"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

type TemplateConfigSettings struct {
	StoragePath string   `mapstructure:"storage_path"`
	AllowedDirs []string `mapstructure:"allowed_dirs"`
}

type AppMetaConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
}

var GlobalConfig *Config

// Init initializes configuration. It accepts an optional config file path.
//
// Resolution order:
//  1. Explicit path passed as argument (if any)
//  2. CONFIG_PATH environment variable
//  3. Default paths: configs/config.yaml, config.yaml
//  4. Pure defaults (if no config file found)
func Init(configPaths ...string) error {
	setDefaults()
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	// Build candidate paths
	var paths []string
	paths = append(paths, configPaths...)
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		paths = append(paths, envPath)
	}
	paths = append(paths, "configs/config.yaml", "config.yaml")

	// Try each path until one succeeds
	for _, p := range paths {
		if _, err := os.Stat(p); err != nil {
			continue
		}
		viper.SetConfigFile(p)
		if err := viper.ReadInConfig(); err == nil {
			break
		}
	}

	// Unmarshal into GlobalConfig (uses defaults for any missing fields)
	GlobalConfig = &Config{}
	if err := viper.Unmarshal(GlobalConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return nil
}

func setDefaults() {
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 9090)
	viper.SetDefault("server.cache_dir", "./cache")
	viper.SetDefault("server.debug", false)
	viper.SetDefault("server.read_timeout", 30)
	viper.SetDefault("server.write_timeout", 30)

	viper.SetDefault("database.enable", false)
	viper.SetDefault("database.driver", "sqlite")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.user", "root")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.db_name", "./data/scaffold.db")
	viper.SetDefault("database.ssl_mode", "disable")
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.conn_max_lifetime", 3600)

	viper.SetDefault("redis.enable", false)
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	viper.SetDefault("security.access_key", "")

	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.filename", "")
	viper.SetDefault("log.max_size", 100)
	viper.SetDefault("log.max_backups", 3)
	viper.SetDefault("log.max_age", 7)
	viper.SetDefault("log.compress", true)

	viper.SetDefault("template.storage_path", "./templates")
	viper.SetDefault("template.allowed_dirs", []string{"./templates"})

	viper.SetDefault("app.name", "scaffold")
	viper.SetDefault("app.version", "1.0.0")
}

// GetServerAddr returns the full server address string
func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// HasAccessKey returns whether an access key is configured
func (c *Config) HasAccessKey() bool {
	return c.Security.AccessKey != ""
}
