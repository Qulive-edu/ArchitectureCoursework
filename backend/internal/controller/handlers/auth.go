package handlers

import (
	"backend/internal/usecase"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
)

type AuthHandler struct {
	svc usecase.UserService
	log *slog.Logger
	jwt *jwtauth.JWTAuth
}

func NewAuthHandler(r *chi.Mux, svc usecase.UserService, log *slog.Logger, jwt *jwtauth.JWTAuth) *AuthHandler {
	h := &AuthHandler{svc: svc, log: log, jwt: jwt}
	r.Post("/auth/register", h.Register)
	r.Post("/auth/login", h.Login)
	return h
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		h.log.Error("decode register: " + err.Error())
		w.WriteHeader(400)
		return
	}
	user, err := h.svc.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		h.log.Error("register: " + err.Error())
		w.WriteHeader(500)
		return
	}
	_, token, _ := h.jwt.Encode(map[string]interface{}{"user_id": user.ID})
	render.JSON(w, r, map[string]interface{}{"user": user, "token": token})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		h.log.Error("decode login: " + err.Error())
		w.WriteHeader(400)
		return
	}
	user, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		h.log.Error("login: " + err.Error())
		w.WriteHeader(401)
		return
	}
	_, token, _ := h.jwt.Encode(map[string]interface{}{"user_id": user.ID})
	render.JSON(w, r, map[string]interface{}{"user": user, "token": token})
}
