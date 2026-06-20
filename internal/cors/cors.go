package cors

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func EnableCORS(router *chi.Mux, trustedOrigins []string) {
	router.Use(cors.Handler(cors.Options{

		AllowedOrigins:   trustedOrigins,
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           7200,
		Debug:            true,
	}))

}
