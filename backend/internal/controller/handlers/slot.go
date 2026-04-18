package handlers

import (
    "backend/internal/usecase"
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/render"
    "log/slog"
)

type SlotHandler struct {
    svc usecase.PlaceService
    log *slog.Logger
}

func NewSlotHandler(r *chi.Mux, svc usecase.PlaceService, log *slog.Logger) *SlotHandler {
    h := &SlotHandler{svc: svc, log: log}
    r.Get("/places/{id}/slots", h.ListSlots)
    return h
}

func (h *SlotHandler) ListSlots(w http.ResponseWriter, r *http.Request) {
    placeID := chi.URLParam(r, "id")
    slots, err := h.svc.ListSlots(r.Context(), placeID)
    if err != nil {
        h.log.Error("list slots: " + err.Error())
        w.WriteHeader(500)
        return
    }
    render.JSON(w, r, slots)
}
