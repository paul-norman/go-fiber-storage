# Fiber Storage Drivers

Storage drivers that implement a common `Storage` interface, designed to be used with [Fiber](https://gofiber.io/). These **ARE NOT** directly compatible with the standard storage drivers for Fiber, but only require minimal code tweaks for them to work with code that uses those.

These differ from the standard Fiber versions in that they allow any data type to be entered and retrieved, and allow them to be recalled either as an `interface{}` or as their original type.

```go
// Storage interface for communicating with different database/key-value providers
type Storage interface {
	// Get the value for the given key.
	// A Result struct is returned allowing various extraction options.
	Get(key string) Result

	// Set the value for the given key along with an optional expiration value, 0 means no expiration.
	// An empty key will flag an error, but empty values are allowed.
	Set(key string, val any, expiry ...time.Duration) error

	// Deletes the values for the given keys.
	// It returns no error if the storage does not contain the keys.
	Delete(keys ...string) error

	// Removes all keys for the specified namespace.
	Reset() error

	// Close the storage, stop any running garbage collectors and closes open connections.
	Close() error
}
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

## Storage Implementations

- [Memcache](./memcache/README.md)
- [Memory](./memory/README.md)
- [MySQL](./mysql/README.md)
- [Postgres](./postgres/README.md)
- [Redis](./redis/README.md)
- [SQLite3](./sqlite3/README.md)
