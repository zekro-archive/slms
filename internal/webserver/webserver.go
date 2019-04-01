package webserver

import (
	"errors"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"github.com/zekroTJA/slms/internal/auth"
	"github.com/zekroTJA/slms/internal/database"
)

// A WebServer handles the REST API
// connections.
type WebServer struct {
	db     database.Middleware
	auth   auth.Provider
	config *Config
	server *fasthttp.Server
	router *routing.Router
}

// Config contains the configuration
// values for the WebServer.
type Config struct {
	Address  string     `json:"address"`
	APIToken string     `json:"api_token"`
	TLS      *ConfigTLS `json:"tls"`
}

// ConfigTLS contains the configuration
// values for TLS encryption for the
// WebServer.
type ConfigTLS struct {
	Use      bool   `json:"use"`
	CertFile string `json:"cert_file"`
	KeyFile  string `json:"key_file"`
}

// NewWebServer creates a new instance
// of WebServer and registers all set
// request handlers.
func NewWebServer(conf *Config, db database.Middleware, authProvider auth.Provider) (*WebServer, error) {
	if len(conf.APIToken) < 8 {
		return nil, errors.New("api_token must have at least 8 characters")
	}

	router := routing.New()

	ws := &WebServer{
		auth:   authProvider,
		db:     db,
		config: conf,
		router: router,
		server: &fasthttp.Server{
			Handler: router.HandleRequest,
		},
	}

	ws.registerHandlers()

	return ws, nil
}

func (ws *WebServer) registerHandlers() {
	api := ws.router.Group("/api")

	api.Use(ws.handlerAuth)
	shortLinksID := api.Get("/shortlinks/<id>", ws.handlerGetShortLink)
	shortLinksID.Post(ws.handlerEditShortLink)
}

// ListenAndServeBlocking starts listening for HTTP requests
// and serving responses to the specified address in the config.
// The server will run in TLS mode when set in the config.
// The startet event loop will block the current go routine.
func (ws *WebServer) ListenAndServeBlocking() error {
	if ws.config.TLS != nil && ws.config.TLS.Use {
		if ws.config.TLS.CertFile == "" || ws.config.TLS.KeyFile == "" {
			return errors.New("cert file and key file must be specified")
		}
		return ws.server.ListenAndServeTLS(
			ws.config.Address, ws.config.TLS.CertFile, ws.config.TLS.KeyFile)
	}

	return ws.server.ListenAndServe(ws.config.Address)
}
