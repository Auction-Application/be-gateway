package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application )routes() http.Handler{
	  router := chi.NewRouter()

	  router.Post("/signup",app.signUpHandler)
	  router.Post("/confirm",app.confirmHandler)
	  router.Post("/signin",app.signInHandler)

	  return router;
}