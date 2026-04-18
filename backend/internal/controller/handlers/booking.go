package handlers

import (
	"backend/internal/usecase"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type BookingHandler struct {
	svc usecase.BookingService
	log *slog.Logger
}

func NewBookingHandler(r *chi.Mux, svc usecase.BookingService, log *slog.Logger, _ interface{}) *BookingHandler {
	h := &BookingHandler{svc: svc, log: log}
	r.Post("/bookings", h.CreateBooking)
	r.Get("/bookings/my", h.GetMine)
	return h
}

func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	// 💡 ВРЕМЕННО: фиксированный ID демо-пользователя
	userID := 1

	var req struct {
		PlaceID int `json:"place_id"`
		SlotID  int `json:"slot_id"`
	}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		h.log.Error("decode booking: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.svc.CreateBooking(r.Context(), userID, req.PlaceID, req.SlotID); err != nil {
		h.log.Error("create booking: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, map[string]bool{"success": true})
}

func (h *BookingHandler) GetMine(w http.ResponseWriter, r *http.Request) {
	// 💡 ВРЕМЕННО: фиксированный ID демо-пользователя
	userID := 1

	bs, err := h.svc.ListMyBookings(r.Context(), userID)
	if err != nil {
		h.log.Error("list mine: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, bs)
}
