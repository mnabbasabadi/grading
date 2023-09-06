// Package db is the implementation of the storage layer using PostgreSQL.
package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	// PostgreSQL driver
	_ "github.com/lib/pq"
	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

// SSLMode ...
type SSLMode string

const (
	// Disable ...
	Disable SSLMode = "disable"
	// Require ...
	Require SSLMode = "require"
	// VerifyCA ...
	VerifyCA SSLMode = "verify-ca"
	// VerifyFull ...
	VerifyFull SSLMode = "verify-full"
)

// Config ...
type Config struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
	Sslmode  SSLMode
}

// Option ...
type Option func(*Config)

// WithUser ...
func WithUser(user string) Option {
	return func(c *Config) {
		c.User = user
	}
}

// WithPassword ...
func WithPassword(password string) Option {
	return func(c *Config) {
		c.Password = password
	}
}

// WithHost ...
func WithHost(host string) Option {
	return func(c *Config) {
		c.Host = host
	}
}

// WithPort ...
func WithPort(port string) Option {
	return func(c *Config) {
		c.Port = port
	}
}

// WithDatabase ...
func WithDatabase(database string) Option {
	return func(c *Config) {
		c.Database = database
	}
}

// WithSSlMode ...
func WithSSlMode(sslmode SSLMode) Option {
	return func(c *Config) {
		c.Sslmode = sslmode
	}
}

// ConnectToPostgres ...
func ConnectToPostgres(options ...Option) (*sqlx.DB, error) {
	config := &Config{}

	// Apply options
	for _, opt := range options {
		opt(config)
	}

	connectionString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		config.User, config.Password, config.Host, config.Port, config.Database, config.Sslmode)

	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")

	return db, nil
}

// ConnectToMySQL ...
func ConnectToMySQL(option ...Option) (*sqlx.DB, error) {
	config := &Config{}

	// Apply options
	for _, opt := range option {
		opt(config)
	}

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.User, config.Password, config.Host, config.Port, config.Database)

	db, err := sqlx.Connect("mysql", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %v", err)
	}

	fmt.Println("Successfully connected to MySQL!")

	return db, nil
}
