package memcache

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/paul-norman/go-fiber-storage"
	mc "github.com/bradfitz/gomemcache/memcache"
	"github.com/gofiber/utils"
)

// Storage interface that is implemented by storage providers
type Storage struct {
	db    *mc.Client
	items *sync.Pool
}

// New creates a new storage
func New(config ...Config) *Storage {
	// Set default config
	cfg := configDefault(config...)

	// Split comma separated servers into slice
	serverList := strings.Split(utils.Trim(cfg.Servers, ' '), ",")

	// Create db
	db := mc.New(serverList...)

	// Set options
	db.Timeout		= cfg.Timeout
	db.MaxIdleConns = cfg.MaxIdleConns

	// Ping database to ensure a connection has been made
	if err := db.Ping(); err != nil {
		panic(err)
	}

	if cfg.Reset {
		if err := db.DeleteAll(); err != nil {
			panic(err)
		}
	}

	// Create storage
	store := &Storage{
		db: db,
		items: &sync.Pool{
			New: func() interface{} {
				return new(mc.Item)
			},
		},
	}

	return store
}

// Get value by key
func (s *Storage) Get(key string) *storage.Result {
	if len(key) <= 0 {
		return &storage.Result{ Value: nil, Error: errors.New("storage keys cannot be zero length"), Missed: false }
	}

	item, err := s.db.Get(key)

	if err == mc.ErrCacheMiss {
		return &storage.Result{ Value: nil, Error: nil, Missed: true }
	} else if err != nil {
		return &storage.Result{ Value: nil, Error: err, Missed: false }
	}

	return &storage.Result{ Value: item.Value, Error: nil, Missed: false }
}

// Set key with value
func (s *Storage) Set(key string, value any, expiry ...time.Duration) error {
	if len(key) <= 0 {
		return errors.New("storage keys cannot be zero length")
	}

	var exp time.Duration = 0
	if len(expiry) > 0 {
		exp = expiry[0]
	}

	item := s.acquireItem()
	item.Key		= key
	item.Value		= value
	item.Expiration = int32(exp.Seconds())

	err := s.db.Set(item)

	s.releaseItem(item)

	return err
}

// Delete entries by key
func (s *Storage) Delete(keys ...string) error {
	if len(keys) <= 0 {
		return errors.New("at least one key is required for Delete")
	}
	
	for _, v := range keys {
		if len(v) == 0 {
			return errors.New("storage keys cannot be zero length (no keys deleted)")
		}
	}

	for _, v := range keys {
		s.db.Delete(v)
	}

	return nil // TODO - This is bad, should combine any errors
}

// Reset all keys
func (s *Storage) Reset() error {
	return s.db.DeleteAll()
}

// Close the database
func (s *Storage) Close() error {
	return nil
}

// Acquire item from pool
func (s *Storage) acquireItem() *mc.Item {
	return s.items.Get().(*mc.Item)
}

// Release item from pool
func (s *Storage) releaseItem(item *mc.Item) {
	if item != nil {
		item.Key = ""
		item.Value = nil
		item.Expiration = 0

		s.items.Put(item)
	}
}

// Return database client
func (s *Storage) Conn() *mc.Client {
	return s.db
}