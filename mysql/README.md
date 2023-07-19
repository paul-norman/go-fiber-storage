# MySQL

A MySQL storage driver using `database/sql` and [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) which also supports prefixing (namespacing) to allow multiple storage spaces to run independently using the same table / connection pool.

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

Install the mysql implementation:

```bash
go get github.com/paul-norman/go-fiber-storage/mysql
```

### Examples

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

### Config

```go
type Config struct {
	// DB Will override ConnectionURI and all other authentication values if used
	//
	// Optional. Default is nil
	DB *sql.DB
	
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

### Default Config

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
