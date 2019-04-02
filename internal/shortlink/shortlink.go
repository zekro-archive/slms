package shortlink

import "time"

// A ShortLink contains the ID, root link,
// short string, created date, access count
// and edited date of a short link.
type ShortLink struct {
	ID        int       `json:"id"`
	RootLink  string    `json:"root_link"`
	ShortLink string    `json:"short_link"`
	Created   time.Time `json:"created"`
	Accesses  int       `json:"accesses"`
	Edited    time.Time `json:"edited"`
}
