package config

import (
	"io/ioutil"
	"os"

	"github.com/zekroTJA/slms/internal/util"

	"github.com/ghodss/yaml"
	"github.com/zekroTJA/slms/internal/database/mysql"
	"github.com/zekroTJA/slms/internal/webserver"
)

// Defining marshal and unmarshal functions
// for config here if they are wanted to be
// swapped for example with json.Marhsal /
// json.Unmarshal
var (
	unmarshal = yaml.Unmarshal
	marshal   = yaml.Marshal
)

// Default config for new config files.
var defConf = &Main{
	WebServer: &webserver.Config{
		Address:         ":443",
		APITokenHash:    "",
		SessionStoreKey: util.GetRandString(64),
		TLS: &webserver.ConfigTLS{
			Use:      true,
			CertFile: "/var/cert/example.com.cer",
			KeyFile:  "/var/cert/example.com.key",
		},
	},
	Database: &mysql.Config{
		Host:     "localhost",
		Username: "slms",
		Database: "slms",
	},
}

// Main contains the main configuration
// for this application.
type Main struct {
	WebServer *webserver.Config `json:"web_server"`
	Database  *mysql.Config     `json:"database"`
}

// OpenAndParse tries to open a confic file from
// passed loc and returns the parsed object.
// If there is no config file existent at the
// give location, it will attempt to create a
// new config file there.
// Attention: conf MAY NOT be nil when the
// function fails!
//
// Params:
//   loc : the config file location
//
// Returns:
//   *Main : the parsed Main config object
//   bool  : equals true if a new config file was created
//           and the default conf was successfully written
//   error : error object when somehting fails
func OpenAndParse(loc string) (*Main, bool, error) {
	data, err := ioutil.ReadFile(loc)
	if os.IsNotExist(err) {
		return createNew(loc)
	}
	if err != nil {
		return nil, false, err
	}

	conf := new(Main)
	err = unmarshal(data, conf)

	return conf, false, err
}

// createFile attempts to create a new config file
// at the defined loc. If this succeeds, this
// function will return true as second return value.
//
// Params:
//   loc : file location where the new config will
//         be created
//
// Returns:
//   *Main : this will be always nil, just returned
//           to fit return values of OpenAndParse
//   bool  : will be true if the config file was
//           successfully created and written to
//   error : error if something fails
func createNew(loc string) (*Main, bool, error) {
	f, err := os.Create(loc)
	if err != nil {
		return nil, false, err
	}
	defer f.Close()

	data, err := marshal(defConf)
	if err != nil {
		return nil, false, err
	}

	_, err = f.Write(data)

	return nil, err == nil, err
}
