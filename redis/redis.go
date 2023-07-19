package redis

import (
	"context"
	"errors"
	"strings"
	"time"

	redis "github.com/redis/go-redis/v9"
)

// Storage interface that is implemented by storage providers
type Storage struct {
	db redis.UniversalClient
	namespace string
}

// New creates a new redis storage
func New(config ...Config) *Storage {
	// Set default config
	cfg := configDefault(config...)

	// Create new redis universal client
	db := cfg.DB
	if db == nil {
		db = redis.NewUniversalClient(cfg.getUniversalOptions())		
	}

	// Test connection
	if err := db.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	// Empty collection if true
	if cfg.Reset {
		if err := db.FlushDB(context.Background()).Err(); err != nil {
			panic(err)
		}
	}
	
	if len(cfg.Namespace) > 0 {
		cfg.Namespace = strings.TrimRight(cfg.Namespace, "/\\ :_-") + ":"
	}

	// Create new store
	return &Storage{
		db: db,
		namespace: cfg.Namespace,
	}
}

// Get value by key
func (s *Storage) Get(key string) (any, error) {
	if len(key) <= 0 {
		return nil, errors.New("storage keys cannot be zero length")
	}

	key = s.namespace + key

	val, err := s.db.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return nil, nil
	}

	return val, err
}

// Set key with value
func (s *Storage) Set(key string, val any, exp time.Duration) error {
	if len(key) <= 0 {
		return errors.New("storage keys cannot be zero length")
	}

	key = s.namespace + key

	return s.db.Set(context.Background(), key, val, exp).Err()
}

// Delete entry by key
func (s *Storage) Delete(key string) error {
	if len(key) <= 0 {
		return errors.New("storage keys cannot be zero length")
	}

	key = s.namespace + key
	
	return s.db.Del(context.Background(), key).Err()
}

// Reset all entries in the namespace
func (s *Storage) Reset() error {
	if s.namespace == "" {
		return s.db.FlushDB(context.Background()).Err()
	} else {
		iter := s.db.Scan(context.Background(), 0, s.namespace + "*", 0).Iterator()
		for iter.Next(context.Background()) {
			if err := s.db.Del(context.Background(), iter.Val()).Err(); err != nil {
				panic(err)
			}
		}
		if err := iter.Err(); err != nil {
			panic(err)
		}
	}

	return nil
}

// Close the database
func (s *Storage) Close() error {
	return s.db.Close()
}

// Return database client
func (s *Storage) Conn() redis.UniversalClient {
	return s.db
}