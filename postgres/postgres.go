package postgres

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Storage interface that is implemented by storage providers
type Storage struct {
	db			*pgxpool.Pool
	gcInterval	time.Duration
	done		chan struct{}
	namespace	string

	sqlSelect	string
	sqlInsert	string
	sqlDelete	string
	sqlReset	string
	sqlGC		string
}

var (
	checkSchemaMsg = "The `value` row has an incorrect data type. " +
		"It should be BYTEA but is instead %s. This will cause encoding-related panics if the DB is not migrated (see https://github.com/gofiber/storage/blob/main/MIGRATE.md)."
	dropQuery = `DROP TABLE IF EXISTS %s;`
	initQuery = []string{
		`CREATE TABLE IF NOT EXISTS %s (
			key		VARCHAR(64) PRIMARY KEY NOT NULL DEFAULT '',
			prefix	VARCHAR(64) NOT NULL DEFAULT '',
			value	BYTEA NOT NULL,
			expiry	BIGINT NOT NULL DEFAULT '0'
		);`,
		`CREATE INDEX IF NOT EXISTS prefix ON %s (prefix);`,
		`CREATE INDEX IF NOT EXISTS expiry ON %s (expiry);`,
	}
	checkSchemaQuery = `SELECT DATA_TYPE FROM INFORMATION_SCHEMA.COLUMNS
		WHERE table_name = '%s' AND COLUMN_NAME = 'value';`
)

// New creates a new storage
func New(config ...Config) *Storage {
	// Set default config
	cfg := configDefault(config...)

	// Select db connection
	var err error
	db := cfg.DB
	if db == nil {
		db, err = pgxpool.New(context.Background(), cfg.getDSN())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		}
	}

	// Ping database
	if err := db.Ping(context.Background()); err != nil {
		panic(err)
	}

	// Drop table if set to true
	if cfg.Reset {
		if _, err := db.Exec(context.Background(), fmt.Sprintf(dropQuery, cfg.Table)); err != nil {
			db.Close()
			panic(err)
		}
	}

	// Init database queries
	for _, query := range initQuery {
		if _, err := db.Exec(context.Background(), fmt.Sprintf(query, cfg.Table)); err != nil {
			db.Close()
			panic(err)
		}
	}

	// Create storage
	store := &Storage{
		db:			db,
		gcInterval:	cfg.GCInterval,
		done:		make(chan struct{}),
		namespace:	cfg.Namespace,
		sqlSelect:	fmt.Sprintf(`SELECT value, expiry FROM %s WHERE key = $1 AND prefix = $2;`, cfg.Table),
		sqlInsert:	fmt.Sprintf("INSERT INTO %s (key, value, expiry, prefix) VALUES ($1, $2, $3, $4) ON CONFLICT (key) DO UPDATE SET value = $5, expiry = $6;", cfg.Table),
		sqlDelete:	fmt.Sprintf("DELETE FROM %s WHERE key = $1 AND prefix = $2;", cfg.Table),
		sqlReset:	fmt.Sprintf("DELETE FROM %s WHERE prefix = $1;", cfg.Table),
		sqlGC:		fmt.Sprintf("DELETE FROM %s WHERE prefix = $1 AND expiry <= $2 AND expiry != 0;", cfg.Table),
	}

	store.checkSchema(cfg.Table)

	// Start garbage collector
	go store.gcTicker()

	return store
}

// Get value by key
func (s *Storage) Get(key string) (any, error) {
	if len(key) <= 0 {
		return nil, errors.New("storage keys cannot be zero length")
	}

	row := s.db.QueryRow(context.Background(), s.sqlSelect, key, s.namespace)

	data	:= []byte{}
	exp		:= int64(0)
	if err := row.Scan(&data, &exp); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// If the expiration time has already passed, then return nil
	if exp != 0 && exp <= time.Now().Unix() {
		return nil, nil
	}

	var decoded interface{}
	err := json.Unmarshal(data, &decoded)

	return decoded, err
}

// Set key with value
func (s *Storage) Set(key string, val any, exp time.Duration) error {
	if len(key) <= 0 {
		return errors.New("storage keys cannot be zero length")
	}

	expSeconds := int64(0)
	if exp != 0 {
		expSeconds = time.Now().Add(exp).Unix()
	}

	value, err := json.Marshal(val)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(context.Background(), s.sqlInsert, key, value, expSeconds, s.namespace, val, expSeconds)

	return err
}

// Delete entry by key
func (s *Storage) Delete(key string) error {
	if len(key) <= 0 {
		return errors.New("storage keys cannot be zero length")
	}

	_, err := s.db.Exec(context.Background(), s.sqlDelete, key, s.namespace)

	return err
}

// Reset all entries in the namespace
func (s *Storage) Reset() error {
	_, err := s.db.Exec(context.Background(), s.sqlReset, s.namespace)

	return err
}

// Close the database
func (s *Storage) Close() error {
	s.done <- struct{}{}
	s.db.Stat()
	s.db.Close()

	return nil
}

// Return database client
func (s *Storage) Conn() *pgxpool.Pool {
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
	_, _ = s.db.Exec(context.Background(), s.sqlGC, s.namespace, t.Unix())
}

func (s *Storage) checkSchema(tableName string) {
	var data []byte

	row := s.db.QueryRow(context.Background(), fmt.Sprintf(checkSchemaQuery, tableName))
	if err := row.Scan(&data); err != nil {
		panic(err)
	}

	if strings.ToLower(string(data)) != "bytea" {
		fmt.Printf(checkSchemaMsg, string(data))
	}
}