package mysql

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// Config defines the config for storage.
type Config struct {
	// DB Will override ConnectionURI and all other authentication values if used
	//
	// Optional. Default is nil
	DB *sqlx.DB

	// Connection string to use for DB. Will override all other authentication values if used
	//
	// Optional. Default is ""
	ConnectionURI string

	// Host name where the DB is hosted
	//
	// Optional. Default is "127.0.0.1"
	Host string

	// Port where the DB is listening on
	//
	// Optional. Default is 3306
	Port int

	// Server username
	//
	// Optional. Default is ""
	Username string

	// Server password
	//
	// Optional. Default is ""
	Password string

	// Database name
	//
	// Optional. Default is "fiber"
	Database string

	// Table name
	//
	// Optional. Default is "fiber_storage"
	Table string

	// Reset clears any existing keys in existing Table
	//
	// Optional. Default is false
	Reset bool

	// Time before deleting expired keys
	//
	// Optional. Default is 10 * time.Second
	GCInterval time.Duration

	// Namespace to allow different types of information in the same table
	//
	// Optional. Default is ""
	Namespace string

	////////////////////////////////////
	// Adaptor related config options //
	////////////////////////////////////

	// MaxIdleConns sets the maximum number of connections in the idle connection pool.
	//
	// Optional. Default is 100.
	MaxIdleConns int

	// MaxOpenConns sets the maximum number of open connections to the database.
	//
	// Optional. Default is 100.
	MaxOpenConns int

	// ConnMaxLifetime sets the maximum amount of time a connection may be reused.
	//
	// Optional. Default is 1 second.
	ConnMaxLifetime time.Duration
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	DB:					nil,
	ConnectionURI:		"",
	Host:				"127.0.0.1",
	Port:				3306,
	Database:			"fiber",
	Table:				"fiber_storage",
	Reset:				false,
	GCInterval:			10 * time.Second,
	MaxOpenConns:		100,
	MaxIdleConns:		100,
	ConnMaxLifetime:	1 * time.Second,
	Namespace:			"",
}

func (c Config) getDSN() string {
	if c.ConnectionURI != "" {
		return c.ConnectionURI
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.Username, c.Password, c.Host, c.Port, c.Database)
}

// Helper function to set default values
func configDefault(config ...Config) Config {
	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	// Set default values
	if cfg.Host == "" {
		cfg.Host = ConfigDefault.Host
	}

	if cfg.Port <= 0 {
		cfg.Port = ConfigDefault.Port
	}

	if cfg.Database == "" {
		cfg.Database = ConfigDefault.Database
	}

	if cfg.Table == "" {
		cfg.Table = ConfigDefault.Table
	}

	if int(cfg.GCInterval.Seconds()) <= 0 {
		cfg.GCInterval = ConfigDefault.GCInterval
	}

	if cfg.MaxIdleConns <= 0 {
		cfg.MaxIdleConns = ConfigDefault.MaxIdleConns
	}

	if cfg.MaxOpenConns <= 0 {
		cfg.MaxOpenConns = ConfigDefault.MaxOpenConns
	}

	if cfg.ConnMaxLifetime == 0 {
		cfg.ConnMaxLifetime = ConfigDefault.ConnMaxLifetime
	}

	return cfg
}