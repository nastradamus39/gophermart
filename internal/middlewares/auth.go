package middlewares

import (
	"context"
	"net/http"

	"github.com/nastradamus39/gophermart/gophermart"
	"github.com/nastradamus39/gophermart/internal/db"
	"github.com/nastradamus39/gophermart/internal/handlers"
)

func UserAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// сессия текущего пользователя
		session, err := gophermart.SessionStore.Get(r, "go-session")

		userLogin, ok := session.Values["userId"].(string)

		if !ok {
			handlers.UnauthorizedResponse(w, r)
			return
		}

		// ищем пользователя в базе
		user, err := db.Repositories().Users.Find(userLogin)

		if err != nil {
			handlers.UnauthorizedResponse(w, r)
			return
		}

		// в контекст передаем ссылку на пользователя
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "user", user)))
	})
}