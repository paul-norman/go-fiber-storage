# Memory

An in-memory storage driver for [Fiber](https://gofiber.io/). No extra support is required for prefixing (namespacing) to allow multiple storage spaces to run independently, just create separate objects.

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
func (s *Storage) Conn() map[string]entry
```

### Installation

Install the memory implementation:

```bash
go get github.com/paul-norman/go-fiber-storage/memory
```

### Examples

Import the storage package.

```go
import "github.com/paul-norman/go-fiber-storage/memory"
```

You can use the following possibilities to create a storage:

```go
// Initialise default config
store1 := memory.New()

// Initialise custom config
sessions := memory.New(memory.Config{
	GCInterval: 10 * time.Second,
})
```

### Config

```go
type Config struct {
	// Time before deleting expired keys
	//
	// Default is 10 * time.Second
	GCInterval time.Duration
}
```

### Default Config

```go
var ConfigDefault = Config{
	GCInterval: 10 * time.Second,
}
```