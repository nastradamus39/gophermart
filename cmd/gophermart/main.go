package main

import (
	"log"
	"net/http"
	"os"

	"github.com/nastradamus39/gophermart/internal/handlers/orders"
	"github.com/nastradamus39/gophermart/internal/handlers/users"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := Router()

	// Logger
	flog, err := os.OpenFile(`server.log`, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer flog.Close()

	log.SetOutput(flog)

	log.Printf("Starting server on %s", "127.0.0.1:8081")

	// запускаем сервер
	err = http.ListenAndServe("127.0.0.1:8081", r)
	if err != nil {
		log.Printf("Не удалось запустить сервер. %s", err)
		return
	}
}

func Router() (r *chi.Mux) {
	r = chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))

	r.Post("/api/user/register", users.RegisterHandler)
	r.Post("/api/user/login", users.LoginHandler)

	r.Post("/api/user/orders", orders.AddOrderHandler)
	r.Get("/api/user/orders", orders.GetOrdersHandler)

	r.Get("/api/user/balance", users.BalanceHandler)
	r.Post("/api/user/balance/withdraw", users.WithdrawHandler)
	r.Get("/api/user/balance/withdrawals", users.WithdrawalsHandler)

	return r
}
