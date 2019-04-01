package webserver

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/zekroTJA/slms/internal/shortlink"

	"github.com/qiangxue/fasthttp-routing"
)

var (
	errNotFound           = errors.New("not found")
	errUpdatedBoth        = errors.New("you can not update short and root link at once")
	errShortAlreadyExists = errors.New("the set short identifyer already exists")
)

const (
	statusOK                  = 200
	statusBadRequest          = 400
	statusUnauthorized        = 401
	statusNotFound            = 404
	statusInternalServerError = 500
)

// --- HELPER FUNCTIONS AND HANDLERS -------------------------------------

func jsonError(ctx *routing.Context, err error, status int) error {
	if err != nil {
		ctx.SetStatusCode(status)
		return fmt.Errorf("{\n  \"code\": %d,\n  \"message\": \"%s\"\n}",
			status, err.Error())
	}
	return nil
}

func jsonResponse(ctx *routing.Context, v interface{}, status int) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return jsonError(ctx, err, statusInternalServerError)
	}

	ctx.SetStatusCode(status)
	_, err = ctx.Write(data)

	return jsonError(ctx, err, statusInternalServerError)
}

func parseJSONBody(ctx *routing.Context, v interface{}) error {
	data := ctx.PostBody()
	err := json.Unmarshal(data, v)
	return jsonError(ctx, err, statusBadRequest)
}

func (ws *WebServer) handlerAuth(ctx *routing.Context) error {
	_, err := ws.auth.Authenticate(ctx)
	return jsonError(ctx, err, statusUnauthorized)
}

func (ws *WebServer) getShortLink(ctx *routing.Context) (*shortlink.ShortLink, error) {
	id := ctx.Param("id")

	sl, err := ws.db.GetShortLink(id, "", "")
	if err != nil {
		return nil, jsonError(ctx, err, statusInternalServerError)
	}

	if sl == nil {
		sl, err = ws.db.GetShortLink("", "", id)
		if err != nil {
			return nil, jsonError(ctx, err, statusInternalServerError)
		}
		if sl == nil {
			return nil, jsonError(ctx, errNotFound, statusNotFound)
		}
	}

	return sl, nil
}

// --- API HANDLERS ------------------------------------------------------

// GET /api/shortlinks/:ID
func (ws *WebServer) handlerGetShortLink(ctx *routing.Context) error {
	sl, err := ws.getShortLink(ctx)
	if err != nil {
		return err
	}

	return jsonResponse(ctx, sl, statusOK)
}

func (ws *WebServer) handlerEditShortLink(ctx *routing.Context) error {
	slUpdated := new(shortlink.ShortLink)
	if err := parseJSONBody(ctx, slUpdated); err != nil {
		return err
	}

	sl, err := ws.getShortLink(ctx)
	if err != nil {
		return err
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

	if err = ws.db.UpdateShortLink(sl.ID, sl); err != nil {
		return jsonError(ctx, err, statusInternalServerError)
	}

	return jsonResponse(ctx, sl, statusOK)
}
