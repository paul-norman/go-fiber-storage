---
id: memcache
title: Memcache
---

# Memcache

A Memcache storage driver using [`bradfitz/gomemcache`](https://github.com/bradfitz/gomemcache). No extra support is required for namespacing to allow multiple storage spaces to run independently, just create separate objects.

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
func (s *Storage) Conn() *mc.Client
```

### Installation

Install the memory implementation:

```bash
go get github.com/paul-norman/go-fiber-storage/memcache
```

### Examples

Import the storage package:

```go
import "github.com/paul-norman/go-fiber-storage/memcache"
```

You can use the following possibilities to create a storage:

```go
// Initialise default config
store1 := memcache.New()

// Initialise custom config
sessions := memcache.New(memcache.Config{
	Servers: "localhost:11211",
})
```

### Config

```go
type Config struct {
	// Server list divided by ,
	// i.e. server1:11211, server2:11212
	//
	// Optional. Default is "127.0.0.1:11211"
	Servers string

	// Reset clears any existing keys in existing Table
	//
	// Optional. Default is false
	Reset bool
}
```

### Default Config

```go
var ConfigDefault = Config{
	Servers: "127.0.0.1:11211",
}
```