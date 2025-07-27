package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"sync"
	"time"
)

type Config struct {
	AppInfo     string // if you're not using a custom type
	LogLevel    string
	Username    string
	Password    string
	Database    string
	Host        string
	SSLMode     string
	Port        int
	ConnMaxOpen int
	ConnMaxIdle int // to be used with NewDB
	ConnMinIdle int // to be used with NewDBWithPgxPool
	Logging     bool
	Tracing     bool
}

var (
	dbInstance *sql.DB
	once       sync.Once
	dbErr      error
)

func (c Config) DSN() string {
	dsn := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		c.Host, c.Port, c.Database, c.Username, c.Password, c.SSLMode,
	)
	return dsn
}

func NewDB(cfg Config) (*sql.DB, error) {
	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s application_name=%s",
			cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database, cfg.SSLMode, cfg.AppInfo,
		)

		dbInstance, dbErr = sql.Open("postgres", dsn)
		if dbErr != nil {
			return
		}

		// Connection pooling
		dbInstance.SetMaxOpenConns(cfg.ConnMaxOpen)
		dbInstance.SetMaxIdleConns(cfg.ConnMaxIdle)
		dbInstance.SetConnMaxLifetime(30 * time.Minute)

		// Check connectivity
		dbErr = dbInstance.Ping()
		if dbErr == nil && cfg.Logging {
			log.Println("[DB] Connected to PostgreSQL with singleton")
		}
	})

	return dbInstance, dbErr
}
