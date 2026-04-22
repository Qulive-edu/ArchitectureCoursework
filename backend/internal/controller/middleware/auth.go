package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/redis/go-redis/v9"
)

type contextKey string

const (
	ClaimsKey contextKey = "jwt_claims"
	UserIDKey contextKey = "user_id"
)

// JWTWithBlacklist — middleware для валидации JWT с проверкой блеклиста в Redis
func JWTWithBlacklist(jwtAuth *jwtauth.JWTAuth, rdb *redis.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Извлекаем токен (заголовок или query-параметр)
			tokenStr := jwtauth.TokenFromHeader(r)
			if tokenStr == "" {
				tokenStr = jwtauth.TokenFromQuery(r)
			}
			if tokenStr == "" {
				http.Error(w, `{"error": "missing token"}`, http.StatusUnauthorized)
				return
			}

			// 2. Проверка блеклиста в Redis
			if banned, err := rdb.Get(r.Context(), "blacklist:"+tokenStr).Result(); err == nil && banned == "1" {
				http.Error(w, `{"error": "token revoked"}`, http.StatusUnauthorized)
				return
			}

			// 3. 👇 ДЕКОДИРУЕМ ТОКЕН (возвращает 2 значения: token, error)
			token, err := jwtAuth.Decode(tokenStr) // ✅ Только 2 переменные
			if err != nil {
				http.Error(w, `{"error": "invalid token: `+err.Error()+`"}`, http.StatusUnauthorized)
				return
			}
			if token == nil {
				http.Error(w, `{"error": "invalid token"}`, http.StatusUnauthorized)
				return
			}

			// 4. 👇 ИЗВЛЕКАЕМ CLAIMS ИЗ ТОКЕНА
			// Приводим к конкретному типу для работы с методами
			jwtToken, ok := token.(jwt.Token)
			if !ok {
				http.Error(w, `{"error": "invalid token format"}`, http.StatusUnauthorized)
				return
			}

			// Извлекаем user_id (стандартный метод для jwx v2)
			userIDRaw, ok := jwtToken.Get("user_id")
			if !ok {
				http.Error(w, `{"error": "invalid token claims: missing user_id"}`, http.StatusUnauthorized)
				return
			}

			// Конвертируем в int (JWT хранит числа как float64)
			var userID int
			switch v := userIDRaw.(type) {
			case float64:
				userID = int(v)
			case int:
				userID = v
			case int64:
				userID = int(v)
			default:
				http.Error(w, `{"error": "invalid user_id type in token"}`, http.StatusUnauthorized)
				return
			}

			// 5. Добавляем данные в контекст запроса
			ctx := context.WithValue(r.Context(), ClaimsKey, jwtToken)
			ctx = context.WithValue(ctx, UserIDKey, userID)
			r = r.WithContext(ctx)

			// 6. Передаём управление дальше
			next.ServeHTTP(w, r)
		})
	}
}

// GetClaimsFromContext — хелпер для получения токена/claims в хендлерах
func GetTokenFromContext(r *http.Request) (jwt.Token, bool) {
	token, ok := r.Context().Value(ClaimsKey).(jwt.Token)
	return token, ok
}

// GetUserIDFromContext — хелпер для получения user_id в хендлерах
func GetUserIDFromContext(r *http.Request) (int, bool) {
	userID, ok := r.Context().Value(UserIDKey).(int)
	return userID, ok
}
