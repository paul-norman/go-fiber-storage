package mysql

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
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
	dropQuery = "DROP TABLE IF EXISTS %s;"
	initQuery = []string{
		`CREATE TABLE IF NOT EXISTS %s ( 
			key			VARCHAR(64) NOT NULL DEFAULT '', 
			namespace	VARCHAR(64) NOT NULL DEFAULT '', 
			value		BLOB NOT NULL, 
			expiry		BIGINT NOT NULL DEFAULT '0', 
			PRIMARY KEY (namespace, key)
			INDEX namespace (namespace),
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
	}
	checkSchemaQuery = `SELECT DATA_TYPE FROM INFORMATION_SCHEMA.COLUMNS
		WHERE table_name = '%s' AND COLUMN_NAME = 'value';`
)

// New creates a new storage
func New(config ...Config) *Storage {
	// Set default config
	cfg := configDefault(config...)

	db := cfg.DB
	if db == nil {
		var err error
		db, err = sqlx.Open("mysql", cfg.getDSN())
		if err != nil {
			panic(err)
		}

		// Set options
		db.SetMaxOpenConns(cfg.MaxOpenConns)
		db.SetMaxIdleConns(cfg.MaxIdleConns)
		db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}

	// Ping database to ensure a connection has been made
	if err := db.Ping(); err != nil {
		panic(err)
	}

	// Drop table if set to true
	if cfg.Reset {
		query := fmt.Sprintf(dropQuery, cfg.Table)
		if _, err = db.Exec(query); err != nil {
			_ = db.Close()
			panic(err)
		}
	}

	// Init database queries
	for _, query := range initQuery {
		query = fmt.Sprintf(query, cfg.Table)
		if _, err := db.Exec(query); err != nil {
			_ = db.Close()
			panic(err)
		}
	}

	// Create storage
	store := &Storage{
		gcInterval:	cfg.GCInterval,
		db:			db,
		done:		make(chan struct{}),
		sqlSelect:	fmt.Sprintf("SELECT value, expiry FROM %s WHERE key = ? AND namespace = ?;", cfg.Table),
		sqlInsert:	fmt.Sprintf("INSERT INTO %s (key, value, expiry, namespace) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE value = ?, expiry = ?;", cfg.Table),
		sqlDelete:	fmt.Sprintf("DELETE FROM %s WHERE key = ? AND namespace = ?;", cfg.Table),
		sqlReset:	fmt.Sprintf("DELETE FROM %s WHERE namespace = ?;", cfg.Table),
		sqlGC:		fmt.Sprintf("DELETE FROM %s WHERE namespace = ? AND expiry <= ? AND expiry != 0;", cfg.Table),
		namespace:	cfg.Namespace,
	}

	store.checkSchema(cfg.Table)

	// Start garbage collector
	go store.gcTicker()

	return store
}

var noRows = "sql: no rows in result set"

// Get value by key
func (s *Storage) Get(key string) (any, error) {
	if len(key) <= 0 {
		return nil, errors.New("storage keys cannot be zero length")
	}

	var store Store
	if err := s.db.Get(&store, s.sqlSelect, key, s.namespace); err != nil {
		if len(store.Key) == 0 {
			return nil, nil
		}
		return nil, err
	}

	// If the expiration time has already passed, then return nil
	if store.Expiry != 0 && store.Expiry <= time.Now().Unix() {
		return nil, nil
	}

	var decoded interface{}
	err := json.Unmarshal(store.Value, &decoded)

	return decoded, err
}

// Set key with value
func (s *Storage) Set(key string, val any, exp time.Duration) error {
	if len(key) <= 0 {
		return errors.New("storage keys cannot be zero length")
	}

	var expSeconds int64
	if exp != 0 {
		expSeconds = time.Now().Add(exp).Unix()
	}

	value, err := json.Marshal(val)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(s.sqlInsert, key, value, expSeconds, s.namespace, val, expSeconds)

	return err
}

// Delete entry by key
func (s *Storage) Delete(key string) error {
	if len(key) <= 0 {
		return errors.New("storage keys cannot be zero length")
	}

	_, err := s.db.Exec(s.sqlDelete, key, s.namespace)

	return err
}

// Reset all keys in the namespace
func (s *Storage) Reset() error {
	_, err := s.db.Exec(s.sqlReset, s.namespace)

	return err
}

// Close the database
func (s *Storage) Close() error {
	s.done <- struct{}{}

	return s.db.Close()
}

// Return database client
func (s *Storage) Conn() *sqlx.DB {
	return s.db
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

func (s *Storage) checkSchema(tableName string) {
	var data []byte

	row := s.db.QueryRow(fmt.Sprintf(checkSchemaQuery, tableName))
	if err := row.Scan(&data); err != nil {
		panic(err)
	}

	if strings.ToLower(string(data)) != "blob" {
		fmt.Printf(checkSchemaMsg, string(data))
	}
}