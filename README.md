# Fiber Storage Drivers

Premade storage drivers that implement a common `Storage` interface, designed to be used with various [Fiber middlewares](https://github.com/gofiber/fiber/tree/master/middleware).

These differ from the standard Fiber versions in that they allow any data type to be entered and retrieved as an `interface{}`.

```go
// Storage interface for communicating with different database/key-value
// providers. Visit https://github.com/gofiber/storage for more info.
type Storage interface {
	// Get gets the value for the given key.
	// `nil, nil` is returned when the key does not exist
	Get(key string) (any, error)

	// Set stores the given value for the given key along
	// with an expiration value, 0 means no expiration.
	// Empty key or value will be ignored without an error.
	Set(key string, val any, exp time.Duration) error

	// Delete deletes the value for the given key.
	// It returns no error if the storage does not contain the key,
	Delete(key string) error

	// Reset resets the storage and delete all keys.
	Reset() error

	// Close closes the storage and will stop any running garbage
	// collectors and open connections.
	Close() error
}
```

## 📑 Storage Implementations

- [Memcache](./memcache/README.md)
- [Memory](./memory/README.md)
- [MySQL](./mysql/README.md)
- [Postgres](./postgres/README.md)
- [Redis](./redis/README.md)
- [SQLite3](./sqlite3/README.md)
