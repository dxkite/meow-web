package application

import (
	"net/http"
)

type ApplicationOption func(app *Application)

func WithRouter(handler http.Handler) ApplicationOption {
	return func(app *Application) {
		app.handler = handler
	}
}

type Application struct {
	handler http.Handler
}

func New(ops ...ApplicationOption) *Application {
	app := &Application{}
	return app
}

func (app *Application) Serve() error {
	return http.ListenAndServe(":2333", app.handler)
}
