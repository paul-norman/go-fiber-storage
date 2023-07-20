package sqlite3

import (
	"errors"
	"fmt"
	"time"

	"github.com/paul-norman/go-fiber-storage"
	"github.com/goccy/go-json"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// Storage interface that is implemented by storage providers
type Storage struct {
	db			*sqlx.DB
	gcInterval	time.Duration
	done		chan struct{}
	namespace	string

	sqlSelect	string
	sqlInsert	string
	sqlDelete	string
	sqlReset	string
	sqlGC		string
}

type Store struct {
	Key			string	`db:"key"`
	Value		[]byte	`db:"value"`
	Namespace	string	`db:"namespace"`
	Expiry		int64	`db:"expiry"`
}

var (
	checkSchemaMsg = "The `value` row has an incorrect data type. " +
		"It should be BLOB but is instead %s. This will cause encoding-related panics if the DB is not migrated (see https://github.com/gofiber/storage/blob/main/MIGRATE.md)."
	dropQuery = `DROP TABLE IF EXISTS %s;`
	initQuery = []string{
		`CREATE TABLE IF NOT EXISTS %s (
			key			VARCHAR(64) PRIMARY KEY NOT NULL DEFAULT '',
			namespace	VARCHAR(64) NOT NULL DEFAULT '',
			value		BLOB NOT NULL,
			expiry		BIGINT NOT NULL DEFAULT '0'
		);`,
		`CREATE INDEX IF NOT EXISTS namespace ON %s (namespace);`,
		`CREATE INDEX IF NOT EXISTS expiry ON %s (expiry);`,
	}
)

// New creates a new storage
func New(config ...Config) *Storage {
	// Set default config
	cfg := configDefault(config...)

	db := cfg.DB
	if db == nil {
		var err error
		db, err = sqlx.Open("sqlite3", cfg.Database)
		if err != nil {
			panic(err)
		}

		// Set database options
		db.SetMaxOpenConns(cfg.MaxOpenConns)
		db.SetMaxIdleConns(cfg.MaxIdleConns)
		db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}

	// Ping database
	if err := db.Ping(); err != nil {
		panic(err)
	}

	// Drop table if set to true
	if cfg.Reset {
		if _, err := db.Exec(fmt.Sprintf(dropQuery, cfg.Table)); err != nil {
			_ = db.Close()
			panic(err)
		}
	}

	// Init database queries
	for _, query := range initQuery {
		if _, err := db.Exec(fmt.Sprintf(query, cfg.Table)); err != nil {
			_ = db.Close()
			panic(err)
		}
	}

	// Create storage
	store := &Storage{
		db:			db,
		gcInterval:	cfg.GCInterval,
		namespace:	cfg.Namespace,
		done:		make(chan struct{}),
		sqlSelect:	fmt.Sprintf(`SELECT key, value, expiry FROM %s WHERE key = ? AND namespace = ?`, cfg.Table),
		sqlInsert:	fmt.Sprintf("INSERT INTO %s (key, value, expiry, namespace) VALUES (?, ?, ?, ?)", cfg.Table),
		sqlDelete:	fmt.Sprintf("DELETE FROM %s WHERE namespace = ? AND key IN (?)", cfg.Table),
		sqlReset:	fmt.Sprintf("DELETE FROM %s WHERE namespace = ?", cfg.Table),
		sqlGC:		fmt.Sprintf("DELETE FROM %s WHERE namespace = ? AND expiry <= ? AND expiry != 0", cfg.Table),
	}

	// Start garbage collector
	go store.gcTicker()

	return store
}

// Get value by key
func (s *Storage) Get(key string) *storage.Result {
	if len(key) <= 0 {
		return &storage.Result{ Value: nil, Error: errors.New("storage keys cannot be zero length"), Missed: false }
	}

	var store Store
	if err := s.db.Get(&store, s.sqlSelect, key, s.namespace); err != nil {
		return &storage.Result{ Value: nil, Error: err, Missed: false }
	}
	if len(store.Key) == 0 || (store.Expiry != 0 && store.Expiry <= time.Now().Unix()) {
		return &storage.Result{ Value: nil, Error: nil, Missed: true }
	}

	var decoded interface{}
	err := json.Unmarshal(store.Value, &decoded)

	return &storage.Result{ Value: decoded, Error: err, Missed: false }
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

	expSeconds := int64(0)
	if exp != 0 {
		expSeconds = time.Now().Add(exp).Unix()
	}

	val, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(s.sqlInsert, key, val, expSeconds, s.namespace)

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

	query, args, err := sqlx.In(s.sqlDelete, s.namespace, keys)
	if err != nil {
		return errors.New("storage keys could not be bound (no keys deleted)")
	}

	query = s.db.Rebind(query)
	_, err = s.db.Exec(query, args...)

	return err
}

// Reset all entries in the namespace
func (s *Storage) Reset() error {
	_, err := s.db.Exec(s.sqlReset, s.namespace)

	return err
}

// Close the database
func (s *Storage) Close() error {
	s.done <- struct{}{}

	return s.db.Close()
}

// gcTicker starts the gc ticker
func (s *Storage) gcTicker() {
	ticker := time.NewTicker(s.gcInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.done:
			return
		case t := <-ticker.C:
			s.gc(t)
		}
	}
}

// gc deletes all expired entries
func (s *Storage) gc(t time.Time) {
	_, _ = s.db.Exec(s.sqlGC, s.namespace, t.Unix())
}

// Return database client
func (s *Storage) Conn() *sqlx.DB {
	return s.db
}