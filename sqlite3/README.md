# SQLite

An SQLite3 storage driver for [Fiber](https://gofiber.io/) using [jmoiron/sqlx](https://jmoiron.github.io/sqlx/) + [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) which also supports namespacing to allow multiple storage spaces to run independently using the same table / connection pool.

### Table of Contents

- [Signatures](#signatures)
- [Installation](#installation)
- [Examples](#examples)
- [Config](#config)
- [Default Config](#default-config)

### Signatures

```go
func New(config ...Config) Storage
func (s *Storage) Get(key string) ([]byte, error)
func (s *Storage) Set(key string, val []byte, exp time.Duration) error
func (s *Storage) Delete(key string) error
func (s *Storage) Reset() error
func (s *Storage) Close() error
func (s *Storage) Conn() *sql.DB
```

### Installation

Install the sqlite3 implementation:

```bash
go get github.com/paul-norman/go-fiber-storage/sqlite3
```

### Examples

Import the storage package:

```go
import "github.com/paul-norman/go-fiber-storage/sqlite3"
```

You can use the following possibilities to create a store *(defaults do not need to be included, just shown for illustrative purposes)*:

```go
// Initialise default config
store1 := sqlite3.New()

// Initialise custom config
sessions := sqlite3.New(sqlite3.Config{
	Database:        "./general.sqlite3",
	Table:           "general_store",
	Reset:           false,
	GCInterval:      10 * time.Second,
	MaxOpenConns:    100,
	MaxIdleConns:    100,
	ConnMaxLifetime: 1 * time.Second,
	Namespace:       "sessions",
})

// Initialise custom config using existing DB connection
db, _ := sqlx.Open("sqlite3", YourConfigStruct)
names := sqlite3.New(sqlite3.Config{
	DB:        db,
	Namespace: "names",
})
```

### Config

```go
type Config struct {
	// DB Will override ConnectionURI and all other authentication values if used
	//
	// Optional. Default is nil
	DB *sqlx.DB
	
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

### Default Config

```go
var ConfigDefault = Config{
	Database:        "./fiber.sqlite3",
	Table:           "fiber_storage",
	Reset:           false,
	GCInterval:      10 * time.Second,
	MaxOpenConns:    100,
	MaxIdleConns:    100,
	ConnMaxLifetime: 1 * time.Second,
	Prefix:          "",
}
```