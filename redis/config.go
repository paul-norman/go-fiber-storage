package redis

import (
	"crypto/tls"
	"fmt"
	"runtime"

	redis "github.com/redis/go-redis/v9"
)

// Config defines the config for storage.
type Config struct {
	// DB object will override connection uri and other connection fields
	//
	// Optional. Default is nil
	DB redis.UniversalClient

	// The standard format redis url to parse all other options. If this is set all other config options, Host, Port, Username, Password, Database have no effect.
	//
	// Example: redis://<user>:<pass>@localhost:6379/<db>
	// Optional. Default is ""
	ConnectionURI string

	// <Host>:<Port> pairs for connection
	//
	// Optional. Default is empty
	Addresses []string

	// Host name where the DB is hosted
	//
	// Optional. Default is "127.0.0.1"
	Host string

	// Port where the DB is listening on
	//
	// Optional. Default is 6379
	Port int

	// Server username
	//
	// Optional. Default is ""
	Username string

	// Server password
	//
	// Optional. Default is ""
	Password string

	// Database to be selected after connecting to the server.
	//
	// Optional. Default is 0
	Database int

	// Key prefix to allow different types of information in the same table
	//
	// Optional. Default is ""
	Namespace string

	// Reset clears any existing keys in existing Collection
	//
	// Optional. Default is false
	Reset bool

	// TLS Config to use. When set TLS will be negotiated.
	TLSConfig *tls.Config

	// Maximum number of socket connections.
	//
	// Optional. Default is 10 connections per every available CPU as reported by runtime.GOMAXPROCS.
	PoolSize int

	////////////////////////////////////
	// Adaptor related config options //
	////////////////////////////////////

	// https://pkg.go.dev/github.com/go-redis/redis/v8#Options
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	DB: 			nil,
	ConnectionURI:	"",
	Host:			"127.0.0.1",
	Port:			6379,
	Username:		"",
	Password:		"",
	Addresses:		[]string{},
	Database:		0,
	Namespace:		"",
	Reset:			false,
	TLSConfig:		nil,
	PoolSize:		10 * runtime.GOMAXPROCS(0),
}

func (c *Config) getUniversalOptions() *redis.UniversalOptions {
	if c.ConnectionURI != "" {
		options, err := redis.ParseURL(c.ConnectionURI)
		if err != nil {
			panic(err)
		}

		// Update the config values with the parsed URL values
		c.Username	= options.Username
		c.Password	= options.Password
		c.Database	= options.DB
		c.Addresses	= []string{ options.Addr }
	} else if len(c.Addresses) > 0 {
		// Fallback to Host and Port values if Addrs is empty
		c.Addresses = []string{ fmt.Sprintf("%s:%d", c.Host, c.Port) }
	}

	return &redis.UniversalOptions{
		Addrs:				c.Addresses,
		//MasterName:		cfg.MasterName,
		//ClientName:		cfg.ClientName,
		//SentinelUsername:	cfg.SentinelUsername,
		//SentinelPassword:	cfg.SentinelPassword,
		DB:					c.Database,
		Username:			c.Username,
		Password:			c.Password,
		TLSConfig:			c.TLSConfig,
		PoolSize:			c.PoolSize,
	}
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

	if cfg.PoolSize <= 0 {
		cfg.PoolSize = ConfigDefault.PoolSize
	}

	return cfg
}