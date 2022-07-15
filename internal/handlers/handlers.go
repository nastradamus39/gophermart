package handlers

import (
	"encoding/json"
	"errors"
	"github.com/nastradamus39/gophermart/gophermart"
	"github.com/nastradamus39/gophermart/internal/db"
	"github.com/nastradamus39/gophermart/internal/luhn"
	"golang.org/x/crypto/bcrypt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// RegisterHandler — регистрация пользователя. Регистрация производится по паре логин/пароль.
// Каждый логин должен быть уникальным.
// После успешной регистрации должна происходить автоматическая аутентификация пользователя.
// Возможные ответы
// 200 — пользователь успешно зарегистрирован и аутентифицирован;
// 400 — неверный формат запроса;
// 409 — логин уже занят;
// 500 — внутренняя ошибка сервера.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user db.User

	// Обрабатываем входящий json
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Считаем хеш пароля для дальнейшего сохранения
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		InternalErrorResponse(w, r, err)
		return
	}
	user.Password = string(hash)

	// Сохраняем пользователя в базу
	err = db.Repositories().Users.Save(&user)
	if err != nil {
		if errors.Is(err, gophermart.ErrUserLoginConflict) {
			http.Error(w, gophermart.ErrUserLoginConflict.Error(), http.StatusConflict)
		} else {
			InternalErrorResponse(w, r, err)
		}
		return
	}

	// Аунтефицируем пользователя
	err = AuthenticateUser(&user, r, w)
	if err != nil {
		InternalErrorResponse(w, r, err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("пользователь успешно зарегистрирован"))
}

// LoginHandler — аутентификация пользователя. Аутентификация производится по паре логин/пароль.
// Возможные коды ответа:
// 200 — пользователь успешно аутентифицирован;
// 400 — неверный формат запроса;
// 401 — неверная пара логин/пароль;
// 500 — внутренняя ошибка сервера.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user db.User
	var err error

	// Обрабатываем входящий json
	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Ищем пользователя в базе
	u, err := db.Repositories().Users.Find(user.Login)
	if err != nil {
		http.Error(w, "неверный логин/пароль", http.StatusUnauthorized)
		return
	}

	// проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
	if err != nil {
		http.Error(w, "неверный логин/пароль", http.StatusUnauthorized)
		return
	}

	// Аунтефицируем пользователя
	err = AuthenticateUser(u, r, w)
	if err != nil {
		InternalErrorResponse(w, r, err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("authenticated"))
}

// BalanceHandler — Хендлер доступен только авторизованному пользователю. В ответе содержатся данные о текущей сумме
// баллов лояльности, а также сумме использованных за весь период регистрации баллов.
// 200 — успешная обработка запроса.
// 401 — пользователь не авторизован.
// 500 — внутренняя ошибка сервера.
func BalanceHandler(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value("user")
	user, ok := u.(*db.User)

	if !ok {
		UnauthorizedResponse(w, r)
	}

	type data struct {
		Current   float32 `json:"current"`
		Withdrawn float32 `json:"withdrawn"`
	}

	// сумма всех списанных балов пользователя
	withDrawSum, err := db.Repositories().Withdraw.WithdrawalsSumByUser(user.Id)
	if err != nil {
		InternalErrorResponse(w, r, err)
		return
	}

	b := data{
		Current:   user.Balance,
		Withdrawn: withDrawSum,
	}

	response, err := json.Marshal(b)

	if err != nil {
		InternalErrorResponse(w, r, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Accept", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// WithdrawHandler — запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа.
// 200 — успешная обработка запроса;
// 401 — пользователь не авторизован;
// 402 — на счету недостаточно средств;
// 422 — неверный номер заказа;
// 500 — внутренняя ошибка сервера.
func WithdrawHandler(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value("user")
	user, ok := u.(*db.User)

	if !ok {
		UnauthorizedResponse(w, r)
	}

	body, _ := ioutil.ReadAll(r.Body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("AddOrderHandler. %s", err)
		}
	}(r.Body)

	type data struct {
		Order string  `json:"order"`
		Sum   float32 `json:"sum"`
	}

	b := data{}
	err := json.Unmarshal(body, &b)
	if err != nil {
		InternalErrorResponse(w, r, err)
		return
	}

	// запрос на списание денег больше чем есть
	if b.Sum > user.Balance {
		w.WriteHeader(http.StatusPaymentRequired)
		w.Write([]byte("на счету недостаточно средств"))
	} else {
		// регистрируем списание в базе
		err := db.Repositories().Withdraw.Save(&db.Withdraw{
			Order:  b.Order,
			Sum:    b.Sum,
			UserId: user.Id,
		})
		if err != nil {
			InternalErrorResponse(w, r, err)
			return
		}

		// списываем с баланса пользователя
		user.Balance = user.Balance - b.Sum
		err = db.Repositories().Users.Save(user)
		if err != nil {
			InternalErrorResponse(w, r, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}

	return
}

// WithdrawalsHandler — получение информации о выводе средств с накопительного счёта пользователем.
// Возможные коды ответа:
// 200 — успешная обработка запроса.
// 204 — нет ни одного списания.
// 401 — пользователь не авторизован.
// 500 — внутренняя ошибка сервера.
func WithdrawalsHandler(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value("user")
	user, ok := u.(*db.User)

	if !ok {
		UnauthorizedResponse(w, r)
		return
	}

	withdrawals, err := db.Repositories().Withdraw.FindWithdrawalsByUser(user.Id)

	if err != nil {
		InternalErrorResponse(w, r, err)
		return
	}

	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	response, err := json.Marshal(withdrawals)

	if err != nil {
		InternalErrorResponse(w, r, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Accept", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// AddOrderHandler - загрузка пользователем номера заказа для расчёта.
// Хендлер доступен только аутентифицированным пользователям. Номером заказа является последовательность цифр
// произвольной длины. Номер заказа может быть проверен на корректность ввода с помощью алгоритма Луна.
// Возможные коды ответа:
// 200 — номер заказа уже был загружен этим пользователем;
// 202 — новый номер заказа принят в обработку;
// 400 — неверный формат запроса;
// 401 — пользователь не аутентифицирован;
// 409 — номер заказа уже был загружен другим пользователем;
// 422 — неверный формат номера заказа;
// 500 — внутренняя ошибка сервера.
func AddOrderHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := ioutil.ReadAll(r.Body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("AddOrderHandler. %s", err)
		}
	}(r.Body)

	orderId := string(id)

	if isValid := luhn.Validate(orderId); !isValid {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("неверный формат номера заказа"))
	}

	u := r.Context().Value("user")
	user, ok := u.(*db.User)

	if !ok {
		UnauthorizedResponse(w, r)
	}

	order := db.Order{
		Persist: false,
		OrderId: orderId,
		Status:  db.ORDER_STATUS_NEW,
		UserId:  user.Id,
		Accrual: 0,
	}

	err := db.Repositories().Orders.Save(&order)

	if errors.Is(err, gophermart.ErrOrderIdConflict) {
		// дополнительно нужно проверить кем был ранее загружен заказ
		o, err := db.Repositories().Orders.Find(orderId)
		if err != nil {
			InternalErrorResponse(w, r, err)
		}

		if o.UserId != user.Id {
			http.Error(w, "номер заказа уже был загружен другим пользователем", http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("номер заказа уже был загружен этим пользователем"))
		}
		return
	}

	if err != nil {
		InternalErrorResponse(w, r, err)
	}

	// отправляем в обработку
	go Accrual(&order, user)

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("новый номер заказа принят в обработку"))
	return
}

// GetOrdersHandler — получение списка загруженных пользователем номеров заказов, статусов их обработки
// и информации о начислениях
// Возможные коды ответа:
// 200 — успешная обработка запроса.
// 204 — нет данных для ответа.
// 401 — пользователь не авторизован.
// 500 — внутренняя ошибка сервера.
func GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value("user")
	user, ok := u.(*db.User)

	if !ok {
		UnauthorizedResponse(w, r)
		return
	}

	orders, err := db.Repositories().Orders.FindByUser(user.Id)

	if err != nil {
		InternalErrorResponse(w, r, err)
		return
	}

	response, err := json.Marshal(orders)

	if err != nil {
		InternalErrorResponse(w, r, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Accept", "application/json")

	if len(orders) > 0 {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}

	w.Write(response)
	return
}
