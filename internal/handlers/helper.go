package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/nastradamus39/gophermart/gophermart"
	"github.com/nastradamus39/gophermart/internal/db"
	"log"
	"net/http"
	"time"
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

// Accrual запрашивает число назначенных балов за заказ. Вносит их на баланс пользователя
func Accrual(order *db.Order, user *db.User) {
	url := fmt.Sprintf("%s/api/orders/%s", gophermart.Cfg.AccrualAddress, order.OrderId)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	status := resp.StatusCode

	// Возможные коды ответа:
	// 200 — успешная обработка запроса.
	// 429 — превышено количество запросов к сервису.
	// 500 — внутренняя ошибка сервера.
	if status == http.StatusOK {
		type respData struct {
			Order   string  `json:"order"`
			Status  string  `json:"status"`
			Accrual float32 `json:"accrual"`
		}

		incomingData := respData{}

		// Обрабатываем входящий json
		if err := json.NewDecoder(resp.Body).Decode(&incomingData); err != nil {
			log.Printf(err.Error())
			return
		}

		// меняем статус заказа
		order.Status = incomingData.Status
		order.Accrual = incomingData.Accrual
		err = db.Repositories().Orders.Save(order)
		if err != nil {
			log.Printf(err.Error())
		}

		// начисляем пользователю балы
		user.Balance = user.Balance + incomingData.Accrual
		err = db.Repositories().Users.Save(user)
		if err != nil {
			log.Printf(err.Error())
		}
	}
	if status == http.StatusTooManyRequests {
		log.Printf("Accrual system response - StatusTooManyRequests")
		time.Sleep(time.Second)
	}
	if status == http.StatusInternalServerError {
		log.Printf("Accrual system response - StatusInternalServerError")
	}

	return
}
