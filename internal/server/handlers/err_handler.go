package handlers

import (
	"errors"
	"karma8/internal/service"
	"karma8/internal/storage"
	"net/http"
)

func (h *Handler) handleErr(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, service.ErrNotEnoughSpace):
		h.replyWithError(r, w, http.StatusOK, service.ErrNotEnoughSpace.Error())
		return
	case errors.Is(err, storage.ErrAlreadyExist):
		h.replyWithError(r, w, http.StatusOK, storage.ErrAlreadyExist.Error())
		return
	case errors.Is(err, storage.ErrNotFound):
		h.replyWithError(r, w, http.StatusNotFound, storage.ErrNotFound.Error())
		return
	case errors.Is(err, service.ErrEmptyFile):
		h.replyWithError(r, w, http.StatusNotFound, service.ErrEmptyFile.Error())
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
