package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var dbScheme = `
CREATE TABLE IF NOT EXISTS shortlinks (
	rootlink text NOT NULL,
	shortlink text NOT NULL,
	created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	accesses bigint(20) NOT NULL DEFAULT '0',
	lastaccess timestamp NOT NULL DEFAULT '0000-00-00 00:00:00'
  ) ENGINE=InnoDB DEFAULT CHARSET=latin1;
`

type MySql struct {
	Dsn string
	DB  *sql.DB
}

func (this *MySql) prepareDatabase() error {
	commands := strings.Split(dbScheme, ";")
	for _, cmd := range commands {
		if strings.Trim(cmd, " \n\t") != "" {
			_, err := this.Query(cmd)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func NewMySql(creds *MySqlCreds) (*MySql, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", creds.Username, creds.Password, creds.Address, creds.Database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	mysql := &MySql{dsn, db}
	_, err = mysql.Query("SHOW TABLES;")
	if err != nil {
		return nil, err
	}
	err = mysql.prepareDatabase()
	if err != nil {
		return nil, err
	}
	return mysql, nil
}

func (this *MySql) Close() {
	if this == nil {
		return
	}
	this.DB.Close()
}

func (this *MySql) Query(statement string, values ...interface{}) (*sql.Rows, error) {
	if this == nil {
		return nil, errors.New("nullptr")
	}
	stm, err := this.DB.Prepare(statement)
	if err != nil {
		return nil, err
	}
	defer stm.Close()

	return stm.Query(values...)
}
