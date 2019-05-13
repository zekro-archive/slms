package webserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"github.com/zekroTJA/slms/internal/auth"
	"github.com/zekroTJA/slms/internal/logger"
	"github.com/zekroTJA/slms/internal/shortlink"
	"github.com/zekroTJA/slms/internal/static"
	"github.com/zekroTJA/slms/internal/util"
)

// Error Objects
var (
	errNotFound           = errors.New("not found")
	errUpdatedBoth        = errors.New("you can not update short and root link at once")
	errShortAlreadyExists = errors.New("the set short identifyer already exists")
	errInvalidArguments   = errors.New("invalid arguments")
)

// Static File Handlers
var (
	fileHandlerStatic = fasthttp.FS{
		Root:       "./web/dist",
		IndexNames: []string{"index.html"},
		PathRewrite: func(ctx *fasthttp.RequestCtx) []byte {
			return ctx.Path()[7:]
		},
	}
)

var allowedRx = regexp.MustCompile(`[\w_\-]+`)

const reservedWords = "manage"

// --- HELPER FUNCTIONS AND HANDLERS -------------------------------------

// jsonError writes the error message of err and the
// passed status to response context and aborts the
// execution of following registered handlers ONLY IF
// err != nil.
// This function always returns a nil error that the
// default error handler can be bypassed.
func jsonError(ctx *routing.Context, err error, status int) error {
	if err != nil {
		ctx.Response.Header.SetContentType("application/json")
		ctx.SetStatusCode(status)
		ctx.SetBodyString(fmt.Sprintf("{\n  \"code\": %d,\n  \"message\": \"%s\"\n}",
			status, err.Error()))
		ctx.Abort()
	}
	return nil
}

// jsonResponse tries to parse the passed interface v
// to JSON and writes it to the response context body
// as same as the passed status code.
// If the parsing fails, this will result in a jsonError
// output of the error with status 500.
// This function always returns a nil error.
func jsonResponse(ctx *routing.Context, v interface{}, status int) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	ctx.Response.Header.SetContentType("application/json")
	ctx.SetStatusCode(status)
	_, err = ctx.Write(data)

	return jsonError(ctx, err, fasthttp.StatusInternalServerError)
}

// parseJSONBody tries to parse a requests JSON
// body to the passed object pointer. If the
// parsing fails, this will result in a jsonError
// output with status 400.
// This function always returns a nil error.
func parseJSONBody(ctx *routing.Context, v interface{}) error {
	data := ctx.PostBody()
	err := json.Unmarshal(data, v)
	if err != nil {
		jsonError(ctx, err, fasthttp.StatusBadRequest)
	}
	return err
}

// getShortLink tries to get the ID or short from the path
// parameter <id> and attempts to find the corresponding
// short link database entry weather by ID or by short string.
// If the attempt fails, this results in a jsonError response
// with wether status code 404 if no link was found or 500 if
// the search attempt in the database failed.
//
// Parameters:
//   ctx         : the routing.Context of the request
//   onlyByShort : bypasses the seatch by ID and attempts to find
//                 shortlink only by short string
//
// Returns:
//   *ShortLink : the found short link
//   bool       : equals true if a shortlink was found
//                and false, if no link was found or an
//                error occured
func (ws *WebServer) getShortLink(ctx *routing.Context, onlyByShort bool) (*shortlink.ShortLink, bool) {
	var sl *shortlink.ShortLink
	var err error
	id := ctx.Param("id")

	if !onlyByShort {
		sl, err = ws.db.GetShortLink(id, "", "")
		if err != nil {
			jsonError(ctx, err, fasthttp.StatusInternalServerError)
			return nil, false
		}
	}

	if sl == nil {
		sl, err = ws.db.GetShortLink("", "", id)
		if err != nil {
			jsonError(ctx, err, fasthttp.StatusInternalServerError)
			return nil, false
		}
		if sl == nil {
			jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
			return nil, false
		}
	}

	return sl, true
}

// checkRequestAuth first checks for a Basic auth
// token as Authorization header. If the header
// has no value or the value does not match with
// the defined token hash, the function attempts to
// decode a session from the passed cookie header.
// If both fails, the auhtorization fails and false
// will be returned.
func (ws *WebServer) checkRequestAuth(ctx *routing.Context) bool {
	_, err := ws.auth.Authenticate(ctx)
	if err != nil {
		s, err := ws.sessions.Get(ctx.RequestCtx, "session")
		if err != nil {
			logger.Debug("WEBSERVER :: AUTH :: %s", err.Error())
			return false
		}
		if s.IsNew {
			logger.Debug("WEBSERVER :: AUTH :: is new")
			return false
		}
		return true
	}
	return err == nil
}

// --- GENERAL HANDLERS --------------------------------------------------

// handlerHeaderServer changes response "Server" header value.
func (ws *WebServer) handlerHeaderServer(ctx *routing.Context) error {
	ctx.Response.Header.SetServer(
		fmt.Sprintf("slms v.%s (%s)", static.AppVersion, static.AppCommit))

	if static.Release != "TRUE" {
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "http://localhost:8081")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "authorization, content-type")
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST, GET, DELETE")
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
	}

	return nil
}

// handlerFileServer checks if the request path points to a
// defined static resource location and then uses the specified
// file server for responding to the request.
// If no valid authentication was provided, an attempt to access
// static page files will always serve the login.html page.
func (ws *WebServer) handlerFileServer(ctx *routing.Context) error {
	if string(ctx.Path()[:7]) == "/manage" {
		fileHandlerStatic.NewRequestHandler()(ctx.RequestCtx)
		ctx.Abort()
	}
	return nil
}

// handlerAuth manages general authorization for
// API endpoints resulting in a jsonError on
// unauthorized request.
func (ws *WebServer) handlerAuth(ctx *routing.Context) error {
	if !ws.checkRequestAuth(ctx) {
		return jsonError(ctx, auth.ErrUnauthorized, fasthttp.StatusUnauthorized)
	}
	return nil
}

// handlerShort handles short link redirect
// requests.
func (ws *WebServer) handlerShort(ctx *routing.Context) error {
	short := ctx.Param("short")
	ctx.Response.Header.SetContentType("text/html")

	sl, err := ws.db.GetShortLink("", "", short)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(
			"<html>" +
				"<body>" +
				"<h1>500 - Internal Error</h1><br/>" +
				"<p>Something went wrong getting the short link data:</p><br/>" +
				"<code>" + err.Error() + "</code>" +
				"</body>" +
				"</html>")
		ctx.Abort()
		return nil
	}

	if sl == nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.SendFile("./web/pages/invalid.html")
		ctx.Abort()
		return nil
	}

	ctx.SetStatusCode(fasthttp.StatusMovedPermanently)
	ctx.Response.Header.Set("Location", sl.RootLink)
	ctx.SetBodyString(
		"<html>" +
			"<head>" +
			"<title>Short Link Management System</title>" +
			"</head>" +
			"<body>" +
			"</body>" +
			"<a href=\"" + sl.RootLink + "\">moved here</a>" +
			"</html>")

	go func() {
		sl.Accesses++
		ws.db.UpdateShortLink(sl.ID, sl)
	}()

	return nil
}

// --- REST API HANDLERS -------------------------------------------------

// POST /api/login
func (ws *WebServer) handlerLogin(ctx *routing.Context) error {
	s, err := ws.sessions.Get(ctx.RequestCtx, "session")
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}
	err = s.Save(ctx.RequestCtx)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}
	return jsonResponse(ctx, struct{}{}, fasthttp.StatusOK)
}

// GET /api/shortlinks
func (ws *WebServer) handlerGetShortLinks(ctx *routing.Context) error {
	var err error
	limit, from := 50, 0

	query := ctx.QueryArgs()

	if query.Has("from") {
		from, err = strconv.Atoi(string(query.Peek("from")))
		if err != nil {
			return jsonError(ctx, err, fasthttp.StatusBadRequest)
		}
		if from < 0 {
			return jsonError(ctx, errors.New("from must be at leats 0"), fasthttp.StatusBadRequest)
		}
	}

	if query.Has("limit") {
		limit, err = strconv.Atoi(string(query.Peek("limit")))
		if err != nil {
			return jsonError(ctx, err, fasthttp.StatusBadRequest)
		}
		if limit < 1 || limit > 1000 {
			return jsonError(ctx, errors.New("limit must be in range (0, 1000]"), fasthttp.StatusBadRequest)
		}
	}

	sls, err := ws.db.GetShortLinks(from, limit)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, map[string]interface{}{
		"n":       len(sls),
		"results": sls,
	}, fasthttp.StatusOK)
}

// POST /api/shortlinks
func (ws *WebServer) handlerCreateShortLink(ctx *routing.Context) error {
	newSl := new(shortlink.ShortLink)
	err := parseJSONBody(ctx, newSl)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	if newSl.RootLink == "" {
		return jsonError(ctx, errInvalidArguments, fasthttp.StatusBadRequest)
	}

	if newSl.ShortLink == "" {
		newSl.ShortLink = util.GetRandString(static.RandShortLen)
	}

	if err = util.CheckIfValidLink(newSl.RootLink, ws.config.OnlyHTTPSRootLink); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	if err = util.CheckIfValidShort(newSl.ShortLink, reservedWords, allowedRx); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	exSl, err := ws.db.GetShortLink("", "", newSl.ShortLink)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}
	if exSl != nil {
		return jsonError(ctx, errShortAlreadyExists, fasthttp.StatusBadRequest)
	}

	resSl, err := ws.db.CreateShortLink(newSl)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, resSl, fasthttp.StatusOK)
}

// GET /api/shortlinks/:ID
func (ws *WebServer) handlerGetShortLink(ctx *routing.Context) error {
	sl, ok := ws.getShortLink(ctx, false)
	if !ok {
		return nil
	}

	return jsonResponse(ctx, sl, fasthttp.StatusOK)
}

// POST /api/shortlinks/:ID
func (ws *WebServer) handlerEditShortLink(ctx *routing.Context) error {
	slUpdated := new(shortlink.ShortLink)
	err := parseJSONBody(ctx, slUpdated)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	sl, ok := ws.getShortLink(ctx, false)
	if !ok {
		return nil
	}

	shortLinkUpdated := slUpdated.ShortLink != "" && sl.ShortLink != slUpdated.ShortLink
	rootLinkUpdated := slUpdated.RootLink != "" && sl.RootLink != slUpdated.RootLink

	if shortLinkUpdated && rootLinkUpdated {
		return jsonError(ctx, errUpdatedBoth, fasthttp.StatusBadRequest)
	}

	if shortLinkUpdated {
		if err := util.CheckIfValidShort(slUpdated.ShortLink, reservedWords, allowedRx); err != nil {
			return jsonError(ctx, err, fasthttp.StatusBadRequest)
		}
		if dsl, err := ws.db.GetShortLink("", "", slUpdated.ShortLink); err != nil {
			return jsonError(ctx, err, fasthttp.StatusInternalServerError)
		} else if dsl != nil {
			return jsonError(ctx, errShortAlreadyExists, fasthttp.StatusBadRequest)
		}
	}

	if shortLinkUpdated {
		sl.ShortLink = slUpdated.ShortLink
	}

	if rootLinkUpdated {
		if err := util.CheckIfValidLink(slUpdated.RootLink, ws.config.OnlyHTTPSRootLink); err != nil {
			return jsonError(ctx, err, fasthttp.StatusBadRequest)
		}
		sl.RootLink = slUpdated.RootLink
	}

	if err := ws.db.UpdateShortLink(sl.ID, sl); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, sl, fasthttp.StatusOK)
}

// DELETE /api/shortlink/:ID
func (ws *WebServer) handlerDeleteShortLink(ctx *routing.Context) error {
	sl, ok := ws.getShortLink(ctx, false)
	if !ok {
		return nil
	}

	err := ws.db.DeleteShortLink(sl.ID)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	return nil
}
