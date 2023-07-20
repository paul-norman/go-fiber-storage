# MySQL

A MySQL storage driver for [Fiber](https://gofiber.io/) using [jmoiron/sqlx](https://jmoiron.github.io/sqlx/) + [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) which also supports prefixing (namespacing) to allow multiple storage spaces to run independently using the same table / connection pool.

This **IS NOT** directly compatible with the standard storage drivers for Fiber, but only requires minimal code tweaks for it to work with code that uses those. This differs from the standard Fiber versions in that it allows any data type to be entered and retrieved, and allows values to be recalled either as an `interface{}` or as its original type.

## Table of Contents

- [Signatures](#signatures)
- [Installation](#installation)
- [Examples](#examples)
- [Config](#config)
- [Default Config](#default-config)

### Signatures

```go
func New(config ...Config) Storage
func (s *Storage) Get(key string) *Result
func (s *Storage) Set(key string, value any, expiry ...time.Duration) error
func (s *Storage) Delete(keys ...string) error
func (s *Storage) Reset() error
func (s *Storage) Close() error
func (s *Storage) Conn() *sql.DB
```

## Installation

Install the mysql implementation:

```bash
go get github.com/paul-norman/go-fiber-storage/mysql
```

## Initialisation

Import the storage package:

```go
import "github.com/paul-norman/go-fiber-storage/mysql"
```

You can use the following possibilities to create a store *(defaults do not need to be included, just shown for illustrative purposes)*:

```go
// Initialise default config
store1 := mysql.New()

// Initialise custom config
sessions := mysql.New(mysql.Config{
	Host:       "127.0.0.1",
	Port:       3306,
	Database:   "general",
	Table:      "general_store",
	Reset:      false,
	GCInterval: 10 * time.Second,
	Namespace:  "sessions",
})

// Initialise custom config using connection string
objects := mysql.New(mysql.Config{
	ConnectionURI: "<username>:<password>@tcp(<host>:<port>)/<database>"
	Table:         "general_store",
	Namespace:     "objects",
})

// Initialise custom config using existing DB connection
db, _ := sqlx.Open("mysql", "<username>:<password>@tcp(<host>:<port>)/<database>")
names := mysql.New(mysql.Config{
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
	"github.com/paul-norman/go-fiber-storage/mysql"
)

func main() {
	store := mysql.New()
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
```

## Default Config

```go
var ConfigDefault = Config{
	ConnectionURI:   "",
	Host:            "127.0.0.1",
	Port:            3306,
	Database:        "fiber",
	Table:           "fiber_storage",
	Reset:           false,
	GCInterval:      10 * time.Second,
	Prefix:          "",
}
```