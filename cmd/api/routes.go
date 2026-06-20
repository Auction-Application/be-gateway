package main

import (
	"net/http"

	"github.com/Auction-Application/be-gateway/internal/cors"
	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	cors.EnableCORS(router, app.config.cors.trustedOrigins)

	router.Post("/v1/signup", app.signUpHandler)
	router.Post("/v1/signup/confirm", app.confirmHandler)
	router.Post("/v1/signin", app.signInHandler)

	// return cors.EnableCORS(app.config.cors.trustedOrigins, router)

	return router
}
