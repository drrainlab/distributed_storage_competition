package handlers

import (
	"encoding/json"
	"io"
	"karma8/internal/service"
	"log"
	"net/http"
)

type Handler struct {
	service service.Service
	done    chan struct{}
}

func NewHandler(service service.Service) *Handler {
	return &Handler{
		service: service,
		done:    make(chan struct{}),
	}
}

func (h *Handler) Store(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("key")
	if filename == "" || r.ContentLength <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.service.Store(r.Context(), filename, uint64(r.ContentLength), r.Body)

	if err != nil {
		h.handleErr(w, r, err)
	}

	h.replyJSON(r, w, http.StatusOK, nil)
}

func (h *Handler) Download(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("key")
	if filename == "" {
		h.replyWithError(r, w, http.StatusBadRequest, "missing key")
		return
	}

	reader, err := h.service.Load(r.Context(), filename)
	if err != nil {
		h.handleErr(w, r, err)
	}

	_, err = io.Copy(w, reader)
	if err != nil {
		log.Printf("transmitting error: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (h *Handler) replyBytes(r *http.Request, w http.ResponseWriter, status int, data []byte) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err := w.Write(data)
	if err != nil {
		log.Printf("error writing response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) replyJSON(r *http.Request, w http.ResponseWriter, status int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("can't marshal response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.replyBytes(r, w, status, data)
}

func (h *Handler) replyWithError(r *http.Request, w http.ResponseWriter, status int, reason string) {
	h.replyJSON(r, w, status, map[string]string{"error_reason": reason})
}

func (h *Handler) Shutdown() {
	close(h.done)
}
