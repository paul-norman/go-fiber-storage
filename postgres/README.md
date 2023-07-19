# Postgres

A Postgres storage driver using [jackc/pgx](https://github.com/jackc/pgx) which also supports prefixing (namespacing) to allow multiple storage spaces to run independently using the same table / connection pool.

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
func (s *Storage) Conn() *pgxpool.Pool
```

### Installation

Install the postgres implementation:

```bash
go get github.com/paul-norman/go-fiber-storage/postgres
```

### Examples

Import the storage package:

```go
import "github.com/paul-norman/go-fiber-storage/postgres"
```

You can use the following possibilities to create a storage:

```go
// Initialise default config
store1 := postgres.New()

// Initialise custom config
sessions := postgres.New(postgres.Config{
	Db:          dbPool,
	Table:       "general_store",
	Reset:       false,
	GCInterval:  10 * time.Second,
	Prefix:      "session",
})
```

### Config

```go
// Config defines the config for storage.
type Config struct {
	// DB pgxpool.Pool object will override connection uri and other connection fields
	//
	// Optional. Default is nil
	DB *pgxpool.Pool

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
	// Optional. Default is 5432
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

	// The SSL mode for the connection
	//
	// Optional. Default is "disable"
	SSLMode string

	// Reset clears any existing keys in existing Table
	//
	// Optional. Default is false
	Reset bool

	// Time before deleting expired keys
	//
	// Optional. Default is 10 * time.Second
	GCInterval time.Duration

	// Prefix to allow different types of information in the same table (namespace)
	//
	// Optional. Default is ""
	Prefix string
}
```

### Default Config
```go
// ConfigDefault is the default config
var ConfigDefault = Config{
	ConnectionURI: "",
	Host:          "127.0.0.1",
	Port:          5432,
	Database:      "fiber",
	Table:         "fiber_storage",
	SSLMode:       "disable",
	Reset:         false,
	GCInterval:    10 * time.Second,
	Prefix:        "",
}
```
