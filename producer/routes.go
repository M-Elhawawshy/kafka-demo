package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /test", app.logRequest(http.HandlerFunc(app.testHandler)))

	return mux
}
