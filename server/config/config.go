package config

import (
	"database/sql"
)

// ServerConfig holds all configuration for the server
type ServerConfig struct {
	Port string
	DB   *sql.DB
}

// NewServerConfig creates a new server configuration with the given parameters
func NewServerConfig(port string, db *sql.DB) *ServerConfig {
	return &ServerConfig{
		Port: port,
		DB:   db,
	}
}
