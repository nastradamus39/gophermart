package handlers

import (
	"log"
	"net/http"

	"github.com/nastradamus39/gophermart/gophermart"
	"github.com/nastradamus39/gophermart/internal/db"
)

// InternalErrorResponse - возвращает пользователю 500 ошибку
func InternalErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Internal server error. %s", err)
	http.Error(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
	return
}

// UnauthorizedResponse - возвращает пользователю 401 ошибку
func UnauthorizedResponse(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	return
}

// AuthenticateUser создает сессию пользователя
func AuthenticateUser(user *db.User, r *http.Request, w http.ResponseWriter) error {
	// авторизуем пользователя
	session, err := gophermart.SessionStore.Get(r, "go-session")
	if err != nil {
		return err
	}
	session.Values["userId"] = user.Login
	err = session.Save(r, w)
	if err != nil {
		return err
	}
	return nil
}
