package webserver

import (
	"errors"

	"github.com/go-gem/sessions"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"github.com/zekroTJA/slms/internal/auth"
	"github.com/zekroTJA/slms/internal/database"
)

// A WebServer handles the REST API
// connections.
type WebServer struct {
	db       database.Middleware
	auth     auth.Provider
	sessions sessions.Store
	config   *Config
	server   *fasthttp.Server
	router   *routing.Router
}

// Config contains the configuration
// values for the WebServer.
type Config struct {
	Address           string     `json:"address"`
	OnlyHTTPSRootLink bool       `json:"only_https_rootlink"`
	APITokenHash      string     `json:"api_token_hash"`
	SessionStoreKey   string     `json:"session_store_key"`
	TLS               *ConfigTLS `json:"tls"`
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
	if len(conf.APITokenHash) < 8 {
		return nil, errors.New("api_token must have at least 8 characters")
	}

	router := routing.New()

	cookieStore := sessions.NewCookieStore([]byte(conf.SessionStoreKey))
	cookieStore.MaxAge(600)

	ws := &WebServer{
		auth:     authProvider,
		sessions: cookieStore,
		db:       db,
		config:   conf,
		router:   router,
		server: &fasthttp.Server{
			Handler: sessions.ClearHandler(router.HandleRequest),
		},
	}

	ws.registerHandlers()

	return ws, nil
}

func (ws *WebServer) registerHandlers() {
	ws.router.Use(ws.handlerHeaderServer, ws.handlerFileServer)

	// /:SHORT
	ws.router.Get("/<short>", ws.handlerShort)

	// /api
	api := ws.router.Group("/api")
	api.Use(ws.handlerAuth)

	// /api/login
	api.Post("/login", ws.handlerLogin)

	// /api/shortlinks
	shortLinks := api.Get("/shortlinks", ws.handlerGetShortLinks)
	shortLinks.Post(ws.handlerCreateShortLink)

	// /api/shortlinks/:ID
	shortLinksID := api.Get("/shortlinks/<id>", ws.handlerGetShortLink)
	shortLinksID.Post(ws.handlerEditShortLink)
	shortLinksID.Delete(ws.handlerDeleteShortLink)
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
