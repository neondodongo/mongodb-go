package mongodb

import (
	"errors"
	"strings"
)

const (
	_maxTimeoutMS = 60000
	_minTimeoutMS = 10000
)

// Config holds the connection info required to interface with MongoDB.
type Config struct {
	// Database represents the Mongo database to connect to; required.
	Database string `json:"database"`

	// Password represents the database user's password.
	Password string `json:"password"`

	// TimeoutMS represents a contextual timeout for operations in milliseconds.
	TimeoutMS int `json:"timeout_ms"`

	// URI represents the database URI used to establish a connection.
	URI string `json:"uri"`

	// Username represents the database username.
	Username string `json:"username"`
}

func (c *Config) sanitizeAndValidate() error {
	if c.Database = strings.TrimSpace(c.Database); c.Database == "" {
		return errors.New("cannot initialize mongo integration with empty or whitespace database")
	}

	if c.URI = strings.TrimSpace(c.URI); c.URI == "" {
		return errors.New("cannot initialize mongo integration with empty or whitespace uri")
	}

	c.Password = strings.TrimSpace(c.Password)
	c.Username = strings.TrimSpace(c.Username)

	if c.TimeoutMS < _minTimeoutMS {
		c.TimeoutMS = _minTimeoutMS
	} else if c.TimeoutMS > _maxTimeoutMS {
		c.TimeoutMS = _maxTimeoutMS
	}

	return nil
}
