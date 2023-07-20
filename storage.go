package storage

import "time"

// Storage interface for communicating with different database/key-value providers
type Storage interface {
	// Get the value for the given key.
	// A Result struct is returned allowing various extraction options.
	Get(key string) *Result

	// Set the value for the given key along with an optional expiration value, 0 means no expiration.
	// An empty key will flag an error, but empty values are allowed.
	Set(key string, val any, expiry ...time.Duration) error

	// Deletes the values for the given keys.
	// It returns no error if the storage does not contain the keys.
	Delete(key ...string) error

	// Removes all keys for the specified namespace.
	Reset() error

	// Close the storage, stop any running garbage collectors and closes open connections.
	Close() error
}