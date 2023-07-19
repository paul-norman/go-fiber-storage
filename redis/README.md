# Redis

A Redis storage driver using [go-redis/redis](https://github.com/go-redis/redis) which also supports namespacing to allow multiple storage spaces to run independently using the same connection pool.

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
func (s *Storage) Conn() redis.UniversalClient
```

### Installation

Install the redis implementation:

```bash
go get github.com/paul-norman/go-fiber-storage/redis
```

### Examples

Import the storage package.

```go
import "github.com/paul-norman/go-fiber-storage/redis"
```

You can use the following possibilities to create a storage:

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

// or just the url with all information
names := redis.New(redis.Config{
    URL:       "redis://<user>:<pass>@127.0.0.1:6379/<db>",
    Reset:     false,
	Namespace: "names",
})
```

### Config
```go
type Config struct {
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

### Default Config

```go
var ConfigDefault = Config{
	Host:      "127.0.0.1",
	Port:      6379,
	Username:  "",
	Password:  "",
	URL:       "",
	Database:  0,
	Reset:     false,
	TLSConfig: nil,
	PoolSize:  10 * runtime.GOMAXPROCS(0),
	Namespace: "",
}
```