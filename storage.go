package storage

import "time"

// Storage interface for communicating with different database/key-value providers
type Storage interface {
	// Get gets the value for the given key.
	// `nil, nil` is returned when the key does not exist.
	Get(key string) (any, error)

	// Set stores the given value for the given key along with an expiration value, 0 means no expiration.
	// An empty key will flag an error, but empty values are allowed.
	Set(key string, val any, exp time.Duration) error

	// Delete deletes the value for the given key.
	// It returns no error if the storage does not contain the key.
	Delete(key string) error

	// Reset deletes all keys for the specified namespace.
	Reset() error

	// Close closes the storage and will stop any running garbage collectors and open connections.
	Close() error
}