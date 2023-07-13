package router

import (
	"karma8/internal/server/handlers"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gorilla/mux"
)

func Create(h *handlers.Handler) *mux.Router {
	router := mux.NewRouter()

	api := router.PathPrefix("/api/v1").Subrouter()
	api.Methods("POST").Path("/store").HandlerFunc(h.Store)
	api.Methods("GET").Path("/download").HandlerFunc(h.Download)

	router.Use(recoverMiddleware)

	return router
}

func recoverMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic occured:%s: \n\nstack: \n%s", err, string(debug.Stack()))
			}
		}()
		h.ServeHTTP(resp, request)
	})
}
