package mysql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/zekroTJA/slms/pkg/multierror"

	// MySQL driver import
	_ "github.com/go-sql-driver/mysql"
	"github.com/zekroTJA/slms/internal/database"
	"github.com/zekroTJA/slms/internal/shortlink"
)

const timeFormat = "2006-01-02 15:04:05"

// MySQL maintains the connection
// to a MySQL database.
type MySQL struct {
	db    *sql.DB
	stmts *prepStmts
}

type prepStmts struct {
	getSLByID    *sql.Stmt
	getSLByRoot  *sql.Stmt
	getSLByShort *sql.Stmt
	updateSLByID *sql.Stmt
}

// Config contains the configuration
// for a MySQL database connection.
type Config struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

// Open attempts to stablishes a
// connection to a MySQL database.
func (m *MySQL) Open(cfg interface{}) error {
	var err error

	conf, ok := cfg.(*Config)
	if !ok {
		return errors.New("cfg is not type of mysql.Config")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		conf.Username, conf.Password, conf.Host, conf.Database)

	if m.db, err = sql.Open("mysql", dsn); err != nil {
		return err
	}

	return m.prepStatements()
}

// Close cleanly closes the
// database connection.
func (m *MySQL) Close() {
	m.db.Close()
}

func (m *MySQL) prepStatements() error {
	var err error
	mErr := multierror.New(nil)

	m.stmts = new(prepStmts)

	m.stmts.getSLByID, err = m.db.Prepare(
		"SELECT `id`, `rootlink`, `shortlink`, `created`, `accesses`, `edited` FROM `shortlinks` " +
			"WHERE `deleted` = 0 AND `id` = ?;")
	mErr.Append(err)

	m.stmts.getSLByRoot, err = m.db.Prepare(
		"SELECT `id`, `rootlink`, `shortlink`, `created`, `accesses`, `edited` FROM `shortlinks` " +
			"WHERE `deleted` = 0 AND `rootlink` = ?;")
	mErr.Append(err)

	m.stmts.getSLByShort, err = m.db.Prepare(
		"SELECT `id`, `rootlink`, `shortlink`, `created`, `accesses`, `edited` FROM `shortlinks` " +
			"WHERE `deleted` = 0 AND `shortlink` = ?;")
	mErr.Append(err)

	m.stmts.updateSLByID, err = m.db.Prepare(
		"UPDATE `shortlinks` SET `shortlink` = ?, `rootlink` = ? " +
			"WHERE `id` = ?;")
	mErr.Append(err)

	return mErr.Concat()
}

// GetShortLink gets a short link object from database by
// id, root link or short link, depending on which was passed
// first (in this order).
// If no short link was found, no error will be returned and
// the returned short link object will be nil.
func (m *MySQL) GetShortLink(id, root, short string) (*shortlink.ShortLink, error) {
	switch {
	case id != "":
		return m.getShortLinkWithStrategy(id, m.stmts.getSLByID)
	case root != "":
		return m.getShortLinkWithStrategy(root, m.stmts.getSLByRoot)
	case short != "":
		return m.getShortLinkWithStrategy(short, m.stmts.getSLByShort)
	default:
		return nil, nil
	}
}

// getShortLinkWithStrategy attempts to find a short link object
// in the database by a given ident which will be passed to a
// strategy (SQL prepared statement) defined in the arguments.
func (m *MySQL) getShortLinkWithStrategy(ident string, strategy *sql.Stmt) (*shortlink.ShortLink, error) {
	var created, edited database.Timestamp
	sl := new(shortlink.ShortLink)

	err := strategy.QueryRow(ident).Scan(
		&sl.ID, &sl.RootLink, &sl.ShortLink, &created, &sl.Accesses, &edited)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	mErr := multierror.New(nil)

	sl.Created, err = created.ToTime(timeFormat)
	mErr.Append(err)

	sl.Edited, err = edited.ToTime(timeFormat)
	mErr.Append(err)

	return sl, mErr.Concat()
}

func (m *MySQL) UpdateShortLink(id int, updated *shortlink.ShortLink) error {
	_, err := m.stmts.updateSLByID.Exec(updated.ShortLink, updated.RootLink, id)
	return err
}
