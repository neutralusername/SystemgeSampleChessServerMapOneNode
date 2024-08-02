package app

import (
	"net/http"

	"github.com/neutralusername/Systemge/HTTP"
)

func (app *App) GetHTTPMessageHandlers() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"/": HTTP.SendDirectory("../frontend"),
	}
}
