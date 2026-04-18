package handlers

import (
	"backend/internal/usecase"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type PlaceHandler struct {
	svc usecase.PlaceService
	log *slog.Logger
}

func NewPlaceHandler(r *chi.Mux, svc usecase.PlaceService, log *slog.Logger) *PlaceHandler {
	h := &PlaceHandler{svc: svc, log: log}
	r.Get("/places", h.ListPlaces)
	r.Get("/places/{id}", h.GetPlace)
	return h
}

func (h *PlaceHandler) ListPlaces(w http.ResponseWriter, r *http.Request) {
	places, err := h.svc.ListPlaces(r.Context())
	if err != nil {
		h.log.Error("list places: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, places)
}

func (h *PlaceHandler) GetPlace(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	place, err := h.svc.GetPlaceByID(r.Context(), id)
	if err != nil {
		h.log.Error("get place: " + err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	render.JSON(w, r, place)
}
