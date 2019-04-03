package main

import (
	"flag"

	"github.com/zekroTJA/slms/internal/auth"

	"github.com/zekroTJA/slms/internal/database/mysql"
	"github.com/zekroTJA/slms/internal/webserver"

	"github.com/zekroTJA/slms/internal/config"
	"github.com/zekroTJA/slms/internal/logger"
)

var (
	flagConfig = flag.String("c", "./config.yml", "config file location")
	flagLogLvl = flag.Int("l", 4, "log level (see https://github.com/op/go-logging/blob/master/level.go#L20)")
)

func main() {
	flag.Parse()

	////////////
	// LOGGER //
	////////////

	logger.Setup(`%{color}▶  %{level:.4s} %{id:04d}%{color:reset} %{message}`, *flagLogLvl)

	////////////
	// CONFIG //
	////////////

	cfg, isNew, err := config.OpenAndParse(*flagConfig)
	if err != nil {
		logger.Fatal("CONFIG :: failed opening and parsing: %s", err.Error())
	}
	if isNew {
		logger.Info("CONFIG :: New config file was created at '%s'. Please edit and restart.", *flagConfig)
		return
	}
	logger.Info("CONFIG :: initialized")
	logger.Debug("CONFIG :: %+v", cfg)

	//////////////
	// DATABASE //
	//////////////

	db := new(mysql.MySQL)
	if err = db.Open(cfg.Database); err != nil {
		logger.Fatal("DATABASE :: failed connecting: %s", err.Error())
	}
	logger.Info("DATABASE :: initialized")

	////////////////
	// WEB SERVER //
	////////////////

	authProvider := auth.NewTokenAuthProvider(cfg.WebServer.APITokenHash)

	logger.Info("WEBSERVER :: running at address %s", cfg.WebServer.Address)
	if !cfg.WebServer.TLS.Use {
		logger.Warning("WEBSERVER :: ATTENTION! WEB SERVER IS CONFIGURED IN NON TLS MODE")
	}

	ws, err := webserver.NewWebServer(cfg.WebServer, db, authProvider)
	if err != nil {
		logger.Fatal("WEBSERVER :: init failed: %s", err.Error())
	}
	if err = ws.ListenAndServeBlocking(); err != nil {
		logger.Fatal("WEBSERVER :: failed startup: %s", err.Error())
	}
}