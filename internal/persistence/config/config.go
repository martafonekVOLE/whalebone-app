package config

import "fmt"

// Config is a configuration struct for database.
type Config struct {
	Host               string
	Port               string
	User               string
	Password           string
	DBName             string
	MaxIdleConnections *int
	MaxOpenConnections *int
	MaxLifeTimeSeconds *int
}

// GetDSN returns a database connection string.
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Prague",
		c.Host, c.User, c.Password, c.DBName, c.Port,
	)
}
