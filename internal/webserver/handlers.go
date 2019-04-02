package webserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/zekroTJA/slms/internal/static"
	"github.com/zekroTJA/slms/internal/util"

	"github.com/zekroTJA/slms/internal/shortlink"

	"github.com/qiangxue/fasthttp-routing"
)

var (
	errNotFound           = errors.New("not found")
	errUpdatedBoth        = errors.New("you can not update short and root link at once")
	errShortAlreadyExists = errors.New("the set short identifyer already exists")
	errInvalidArguments   = errors.New("invalid arguments")
)

const (
	statusOK                  = 200
	statusMovedpermanently    = 301
	statusBadRequest          = 400
	statusUnauthorized        = 401
	statusNotFound            = 404
	statusInternalServerError = 500
)

// --- HELPER FUNCTIONS AND HANDLERS -------------------------------------

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

func jsonResponse(ctx *routing.Context, v interface{}, status int) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return jsonError(ctx, err, statusInternalServerError)
	}

	ctx.Response.Header.SetContentType("application/json")
	ctx.SetStatusCode(status)
	_, err = ctx.Write(data)

	return jsonError(ctx, err, statusInternalServerError)
}

func parseJSONBody(ctx *routing.Context, v interface{}) error {
	data := ctx.PostBody()
	err := json.Unmarshal(data, v)
	return jsonError(ctx, err, statusBadRequest)
}

func (ws *WebServer) getShortLink(ctx *routing.Context, onlyByShort bool) (*shortlink.ShortLink, bool) {
	var sl *shortlink.ShortLink
	var err error
	id := ctx.Param("id")

	if !onlyByShort {
		sl, err = ws.db.GetShortLink(id, "", "")
		if err != nil {
			jsonError(ctx, err, statusInternalServerError)
			return nil, false
		}
	}

	if sl == nil {
		sl, err = ws.db.GetShortLink("", "", id)
		if err != nil {
			jsonError(ctx, err, statusInternalServerError)
			return nil, false
		}
		if sl == nil {
			jsonError(ctx, errNotFound, statusNotFound)
			return nil, false
		}
	}

	return sl, true
}

// --- GENERAL HANDLERS --------------------------------------------------

// Cahnges response "Server" header value
func (ws *WebServer) handlerHeaderServer(ctx *routing.Context) error {
	ctx.Response.Header.SetServer(
		fmt.Sprintf("slms v.%s (%s)", static.AppVersion, static.AppCommit))
	return nil
}

// General Authorization handler
func (ws *WebServer) handlerAuth(ctx *routing.Context) error {
	_, err := ws.auth.Authenticate(ctx)
	return jsonError(ctx, err, statusUnauthorized)
}

// Actual short link handler
func (ws *WebServer) handlerShort(ctx *routing.Context) error {
	short := ctx.Param("short")
	ctx.Response.Header.SetContentType("text/html")

	sl, err := ws.db.GetShortLink("", "", short)
	if err != nil {
		ctx.SetStatusCode(statusInternalServerError)
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
		ctx.SetStatusCode(statusNotFound)
		ctx.SetBodyString(
			"<html>" +
				"<body>" +
				"<h1>404 - Not Found</h1><br/>" +
				"<p>No redirection could be found for this short link.</p>" +
				"</body>" +
				"</html>")
		ctx.Abort()
		return nil
	}

	ctx.SetStatusCode(statusMovedpermanently)
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
	return nil
}

// --- REST API HANDLERS -------------------------------------------------

// GET /api/shortlinks
func (ws *WebServer) handlerGetShortLinks(ctx *routing.Context) error {
	var err error
	limit, from := 50, 0

	query := ctx.QueryArgs()

	if query.Has("from") {
		from, err = strconv.Atoi(string(query.Peek("from")))
		if err != nil {
			return jsonError(ctx, err, statusBadRequest)
		}
		if from < 1 {
			return jsonError(ctx, errors.New("from must be larger than 0"), statusBadRequest)
		}
	}

	if query.Has("limit") {
		limit, err = strconv.Atoi(string(query.Peek("limit")))
		if err != nil {
			return jsonError(ctx, err, statusBadRequest)
		}
		if limit < 1 || limit > 1000 {
			return jsonError(ctx, errors.New("limit must be in range (0, 1000]"), statusBadRequest)
		}
	}

	sls, err := ws.db.GetShortLinks(from, limit)
	if err != nil {
		return jsonError(ctx, err, statusInternalServerError)
	}

	return jsonResponse(ctx, map[string]interface{}{
		"n":       len(sls),
		"results": sls,
	}, statusOK)
}

// POST /api/shortlinks
func (ws *WebServer) handlerCreateShortLink(ctx *routing.Context) error {
	newSl := new(shortlink.ShortLink)
	err := parseJSONBody(ctx, newSl)
	if err != nil {
		return jsonError(ctx, err, statusBadRequest)
	}

	if newSl.RootLink == "" {
		return jsonError(ctx, errInvalidArguments, statusBadRequest)
	}

	if newSl.ShortLink == "" {
		newSl.ShortLink = util.GetRandString(static.RandShortLen)
	}

	exSl, err := ws.db.GetShortLink("", "", newSl.ShortLink)
	if err != nil {
		return jsonError(ctx, err, statusInternalServerError)
	}
	if exSl != nil {
		return jsonError(ctx, errShortAlreadyExists, statusBadRequest)
	}

	resSl, err := ws.db.CreateShortLink(newSl)
	if err != nil {
		return jsonError(ctx, err, statusInternalServerError)
	}

	return jsonResponse(ctx, resSl, statusOK)
}

// GET /api/shortlinks/:ID
func (ws *WebServer) handlerGetShortLink(ctx *routing.Context) error {
	sl, ok := ws.getShortLink(ctx, false)
	if !ok {
		return nil
	}

	return jsonResponse(ctx, sl, statusOK)
}

// POST /api/shortlink/:ID
func (ws *WebServer) handlerEditShortLink(ctx *routing.Context) error {
	slUpdated := new(shortlink.ShortLink)
	if err := parseJSONBody(ctx, slUpdated); err != nil {
		return err
	}

	sl, ok := ws.getShortLink(ctx, false)
	if !ok {
		return nil
	}

	shortLinkUpdated := slUpdated.ShortLink != "" && sl.ShortLink != slUpdated.ShortLink
	rootLinkUpdated := slUpdated.RootLink != "" && sl.RootLink != slUpdated.RootLink

	if shortLinkUpdated && rootLinkUpdated {
		return jsonError(ctx, errUpdatedBoth, statusBadRequest)
	}

	if shortLinkUpdated {
		if dsl, err := ws.db.GetShortLink("", "", slUpdated.ShortLink); err != nil {
			return jsonError(ctx, err, statusInternalServerError)
		} else if dsl != nil {
			return jsonError(ctx, errShortAlreadyExists, statusBadRequest)
		}
	}

	if shortLinkUpdated {
		sl.ShortLink = slUpdated.ShortLink
	}

	if rootLinkUpdated {
		sl.RootLink = slUpdated.RootLink
	}

	if err := ws.db.UpdateShortLink(sl.ID, sl); err != nil {
		return jsonError(ctx, err, statusInternalServerError)
	}

	return jsonResponse(ctx, sl, statusOK)
}

// DELETE /api/shortlink/:ID
func (ws *WebServer) handlerDeleteShortLink(ctx *routing.Context) error {
	sl, ok := ws.getShortLink(ctx, false)
	if !ok {
		return nil
	}

	err := ws.db.DeleteShortLink(sl.ID)
	if err != nil {
		return jsonError(ctx, err, statusInternalServerError)
	}

	ctx.SetStatusCode(statusOK)
	return nil
}
