package handlers

import (
	"backend/internal/usecase"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/redis/go-redis/v9"
)

type AuthHandler struct {
	svc usecase.UserService
	log *slog.Logger
	jwt *jwtauth.JWTAuth
	rdb *redis.Client // 👈 Добавлено поле для Redis
}

func NewAuthHandler(r *chi.Mux, svc usecase.UserService, log *slog.Logger, jwt *jwtauth.JWTAuth, rdb *redis.Client) *AuthHandler {
	h := &AuthHandler{svc: svc, log: log, jwt: jwt, rdb: rdb} // 👈 Сохраняем rdb

	r.Post("/auth/register", h.Register)
	r.Post("/auth/login", h.Login)
	r.Post("/auth/logout", h.Logout) // 👈 Регистрация маршрута логаута
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := h.svc.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		h.log.Error("register: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		h.log.Error("login: " + err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	_, token, _ := h.jwt.Encode(map[string]interface{}{"user_id": user.ID})
	render.JSON(w, r, map[string]interface{}{"user": user, "token": token})
}

// 👇 НОВЫЙ МЕТОД: ЛОГАУТ
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// 1. Берём токен из заголовка
	tokenStr := jwtauth.TokenFromHeader(r)
	if tokenStr == "" {
		http.Error(w, `{"error": "missing token"}`, http.StatusUnauthorized)
		return
	}

	// 2. Добавляем в блеклист Redis
	// TTL: 24 часа (в продакшене можно парсить `exp` из JWT и ставить точное время)
	err := h.rdb.Set(r.Context(), "blacklist:"+tokenStr, "1", 24*time.Hour).Err()
	if err != nil {
		h.log.Error("redis blacklist error: " + err.Error())
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}

	h.log.Info("user logged out, token blacklisted")
	render.JSON(w, r, map[string]bool{"success": true})
}
