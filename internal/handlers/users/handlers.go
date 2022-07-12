package users

import (
	"net/http"
)

// RegisterHandler — регистрация пользователя.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("регистрация пользователя"))
}

// LoginHandler — аутентификация пользователя.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("аутентификация пользователя"))
}

// BalanceHandler — получение текущего баланса лояльности пользователя.
func BalanceHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("получение текущего баланса лояльности пользователя"))
}

// WithdrawHandler — запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа.
func WithdrawHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа"))
}

// WithdrawalsHandler — получение информации о выводе средств с накопительного счёта пользователем.
func WithdrawalsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("получение информации о выводе средств с накопительного счёта пользователем"))
}
