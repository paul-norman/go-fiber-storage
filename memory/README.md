# Memory

An in-memory storage driver for [Fiber](https://gofiber.io/). No extra support is required for prefixing (namespacing) to allow multiple storage spaces to run independently, just create separate objects.

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
func (s *Storage) Conn() map[string]entry
```

## Installation

Install the memory implementation:

```bash
go get github.com/paul-norman/go-fiber-storage/memory
```

## Initialisation

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
	GCInterval: 5 * time.Second,
})
```

## Usage

```go
package main

import (
	"fmt"
	"time"
	"github.com/paul-norman/go-fiber-storage/memory"
)

func main() {
	store := memory.New()
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
	// Time before deleting expired keys
	//
	// Default is 10 * time.Second
	GCInterval time.Duration
}
```

## Default Config

```go
var ConfigDefault = Config{
	GCInterval: 10 * time.Second,
}
```