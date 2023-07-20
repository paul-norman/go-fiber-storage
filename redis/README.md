# Redis

A Redis storage driver for [Fiber](https://gofiber.io/) using [go-redis/redis](https://github.com/go-redis/redis) which also supports namespacing to allow multiple storage spaces to run independently using the same connection pool.

This **IS NOT** directly compatible with the standard storage drivers for Fiber, but only requires minimal code tweaks for it to work with code that uses those. This differs from the standard Fiber versions in that it allows any data type to be entered and retrieved, and allows values to be recalled either as an `interface{}` or as its original type.

## Table of Contents

- [Signatures](#signatures)
- [Installation](#installation)
- [Examples](#examples)
- [Config](#config)
- [Default Config](#default-config)

## Signatures

```go
func New(config ...Config) Storage
func (s *Storage) Get(key string) *Result
func (s *Storage) Set(key string, value any, expiry ...time.Duration) error
func (s *Storage) Delete(keys ...string) error
func (s *Storage) Reset() error
func (s *Storage) Close() error
func (s *Storage) Conn() redis.UniversalClient
```

## Installation

Install the redis implementation:

```bash
go get github.com/paul-norman/go-fiber-storage/redis
```

## Initialisation

Import the storage package.

```go
import "github.com/paul-norman/go-fiber-storage/redis"
```

You can use the following possibilities to create a storage *(defaults do not need to be included, just shown for illustrative purposes)*:

```go
// Initialise default config
store1 := redis.New()

// Initialise custom config
sessions := redis.New(redis.Config{
	Host:      "127.0.0.1",
	Port:      6379,
	Username:  "",
	Password:  "",
	URL:       "",
	Database:  0,
	Reset:     false,
	TLSConfig: nil,
	PoolSize:  10 * runtime.GOMAXPROCS(0),
	Namespace: "sessions",
})

// Initialise custom config using connection string
objects := redis.New(redis.Config{
	ConnectionURI: "redis://<username>:<password>@<host>:<port>/<database>",
	Reset:         false,
	Namespace:     "objects",
})

// Initialise custom config using existing DB connection
db := redisClient.NewUniversalClient(&redisClient.redis.UniversalOptions{ ... })
names := redis.New(redis.Config{
	DB:        db,
	Namespace: "names",
})
```

## Usage

```go
package main

import (
	"fmt"
	"time"
	"github.com/paul-norman/go-fiber-storage/redis"
)

func main() {
	store := redis.New()
	defer store.Close()

	// Error handling omitted for brevity
	err := store.Set("my_key", "lives forever")
	err  = store.Set("one_second", "lives for 1 second", 1 * time.Second)
	err  = store.Set("complex_type", map[string]int64{ "test": 123 })

	// Fetch the Result struct, check that we found a value and that there wasn't an error
	result := store.Get("my_key")
	if !result.Miss() && result.Err() == nil {
		fmt.Println("Value: " + result.String()) // Convert the interface{} to a string
	}

	// Fetch parsed information (value, error, miss) separately
	str, err, miss := store.Get("one_second").String()
	if !miss && err == nil {
		fmt.Println("Value: " + str) // If we are here the value is a string
	}

	// Sleep for 2 seconds
	time.Sleep(2 * time.Second)

	// This will result in a miss
	str, err, miss = store.Get("one_second").String()
	if miss {
		fmt.Println("The value has gone...") // str and err will be nil
	}

	// Complex types are also simple
	item, err, miss := store.Get("complex_type").Interface()
	if !miss && err == nil {
		fmt.Println(item.(map[string]int64)) // Convert the interface to the desired type
	}

	// Remove the keys - doesn't matter that one has already expired
	err = store.Delete("my_key", "one_second", "complex_type")

	// Or
	err = store.Reset()
}
```

## Config Options
```go
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

	// URL the standard format redis url to parse all other options. If this is set all other config options, Host, Port, Username, Password, Database have no effect.
	//
	// Example: redis://<user>:<pass>@localhost:6379/<db>
	// Optional. Default is ""
	URL string

	// Reset clears any existing keys in existing Collection
	//
	// Optional. Default is false
	Reset bool

	// TLS Config to use. When set TLS will be negotiated.
	//
	// Optional. Default is nil
	TLSConfig *tls.Config

	// Maximum number of socket connections.
	//
	// Optional. Default is 10 connections per every available CPU as reported by runtime.GOMAXPROCS.
	PoolSize int

	// Namespace to allow different types of information in the same database
	//
	// Optional. Default is ""
	Namespace string
}
```

## Default Config

```go
var ConfigDefault = Config{
	DB:            nil,
	ConnectionURI: "",
	Addresses:     []string{},
	Host:          "127.0.0.1",
	Port:          6379,
	Username:      "",
	Password:      "",
	URL:           "",
	Database:      0,
	Reset:         false,
	TLSConfig:     nil,
	PoolSize:      10 * runtime.GOMAXPROCS(0),
	Namespace:     "",
}
```