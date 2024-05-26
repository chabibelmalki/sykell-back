package api

import (
	"net/http"
	"sykell-back/api/handler"
	"sykell-back/api/middleware"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	// Création d'un sous-routeur pour les routes nécessitant une authentification
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware)

	// Appliquer le middleware aux routes
	api.HandleFunc("/urlinfo", handler.FetchURLInfoHandler).Methods(http.MethodPost)

	return r
}

func HandleCORS(r *mux.Router) http.Handler {
	c := cors.AllowAll()
	return c.Handler(r)
}
