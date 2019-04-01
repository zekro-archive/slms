package database

import (
	"time"

	"github.com/zekroTJA/slms/internal/shortlink"
)

// Timestamp is a uint8 array typed
// timetsmap which can be transformed
// to time.Time
type Timestamp []uint8

// ToTime parses the timestamp to a time object
func (t Timestamp) ToTime(format string) (time.Time, error) {
	return time.Parse(format, string(t))
}

// The Middleware interface describes
// the functions a database middleware
// must provide.
type Middleware interface {
	// Open initializes the database
	// connection with the passed
	// parameters.
	Open(cfg interface{}) error
	// Close closes an existing
	// connection.
	Close()

	// GetShortLink gets a shortlink entry from
	// database wether by id, root or short link
	// (excatly in this order).
	GetShortLink(id, root, short string) (*shortlink.ShortLink, error)

	// UpdateShortLink updates a short link by
	// all values contained in updated.
	UpdateShortLink(id int, updated *shortlink.ShortLink) error
}
