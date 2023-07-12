package handlers

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gorilla/mux"
)

const apiPrefix = "/api/v1"

func CreateRouter(s *StorageService) *mux.Router {
	router := mux.NewRouter()

	api := router.PathPrefix(apiPrefix).Subrouter()
	api.Methods(http.MethodPost).Path(filePath).HandlerFunc(s.UploadFile)
	api.Methods(http.MethodGet).Path(filePath).HandlerFunc(s.GetFile)

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
