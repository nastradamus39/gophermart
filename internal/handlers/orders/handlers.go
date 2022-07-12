package orders

import (
	"net/http"
)

// AddOrderHandler — загрузка пользователем номера заказа для расчёта.
func AddOrderHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("загрузка пользователем номера заказа для расчёта"))
}

// GetOrdersHandler — получение списка загруженных пользователем номеров заказов,
// статусов их обработки и информации о начислениях
func GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("получение списка загруженных пользователем номеров заказов,\n" +
		"// статусов их обработки и информации о начислениях"))
}
