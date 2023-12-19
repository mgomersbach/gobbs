package config

import (
	"flag"

	"github.com/spf13/viper"
)

// CommandLineArgs holds the values of command-line arguments
type CommandLineArgs struct {
	ConfigPath   string
	DatabaseType string
	DatabaseDSN  string
	BBSType      string
	BBSAddress   string
	AuthMethod   string
}

// Config represents the application configuration
type Config struct {
	// Add fields corresponding to configuration
	Database struct {
		Type string
		DSN  string
	}
	BBS struct {
		Type    string // "tcp" or "unix"
		Address string // IP:Port for TCP, file path for UNIX
	}
	Authentication struct {
		Method string
		// Other fields like DB connection info, LDAP server address, etc.
	}
	// Add other configuration fields as needed
}

// ProcessFlags processes command line flags and returns a CommandLineArgs struct
func ProcessFlags() CommandLineArgs {
	args := CommandLineArgs{}

	flag.StringVar(&args.ConfigPath, "config", "config.yaml", "path to config file")
	flag.StringVar(&args.DatabaseType, "database-type", "", "database type")
	flag.StringVar(&args.DatabaseDSN, "database-dsn", "", "database DSN")
	flag.StringVar(&args.BBSType, "bbs-type", "", "BBS type")
	flag.StringVar(&args.BBSAddress, "bbs-address", "", "BBS address")
	flag.StringVar(&args.AuthMethod, "auth-method", "", "authentication method")

	flag.Parse()
	return args
}

// LoadConfig reads the configuration file and merges command line arguments
func LoadConfig(args CommandLineArgs) (*Config, error) {
	viper.SetConfigFile(args.ConfigPath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Override configuration with command line arguments if they are provided
	if args.DatabaseType != "" {
		cfg.Database.Type = args.DatabaseType
	}
	if args.DatabaseDSN != "" {
		cfg.Database.DSN = args.DatabaseDSN
	}
	if args.BBSType != "" {
		cfg.BBS.Type = args.BBSType
	}
	if args.BBSAddress != "" {
		cfg.BBS.Address = args.BBSAddress
	}
	if args.AuthMethod != "" {
		cfg.Authentication.Method = args.AuthMethod
	}

	return &cfg, nil
}