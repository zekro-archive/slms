package database

import (
	"strings"
	"time"

	"github.com/zekroTJA/slms/internal/shortlink"
)

// Timestamp is a uint8 array typed
// timetsmap which can be transformed
// to time.Time
type Timestamp []uint8

// ToTime parses the timestamp to a time object
func (t Timestamp) ToTime(format string) (time.Time, error) {
	tobj, err := time.Parse(format, string(t))
	if err != nil && strings.Contains(err.Error(), "out of range") {
		return time.Time{}, nil
	}
	return tobj, err
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

	// GetShortLinkCount returns the number of short
	// link entries in the database.
	GetShortLinkCount() (int, error)
	// GetShortLink gets a shortlink entry from
	// database wether by id, root or short link
	// (excatly in this order).
	GetShortLink(id, root, short string) (*shortlink.ShortLink, error)
	// GetShortLinks returns a list of short links which
	// is ordered by created date descending between
	// from index and limit ammount.
	GetShortLinks(from, limit int) ([]*shortlink.ShortLink, error)
	// UpdateShortLink updates a short link by
	// all values contained in updated.
	UpdateShortLink(id int, updated *shortlink.ShortLink) error
	// CreateShortLink creates a new shortlink
	// entry in the database and returnes the
	// new shortlink object whis was created.
	CreateShortLink(sl *shortlink.ShortLink) (*shortlink.ShortLink, error)
	// Deletes a shortlink from the database
	// or marks it at least as unavailable.
	DeleteShortLink(id int) error
}
