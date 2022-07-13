package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/nastradamus39/gophermart/gophermart"
	"github.com/nastradamus39/gophermart/internal/db"
	"github.com/nastradamus39/gophermart/internal/handlers"
	"github.com/nastradamus39/gophermart/internal/middlewares"

	"github.com/caarlos0/env/v6"
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

	// Переменные окружения в конфиг
	err = env.Parse(&gophermart.Cfg)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Параметры командной строки в конфиг
	flag.StringVar(&gophermart.Cfg.ServerAddress, "a", gophermart.Cfg.ServerAddress, "Адрес и порт запуска сервиса")
	flag.StringVar(&gophermart.Cfg.DatabaseDsn, "d", gophermart.Cfg.DatabaseDsn, "Адрес подключения к базе данных")
	flag.StringVar(&gophermart.Cfg.AccrualAddress, "r", gophermart.Cfg.AccrualAddress, "Адрес системы расчёта начислений")
	flag.Parse()

	err = db.InitDB()
	if err != nil {
		log.Fatal(err)
		return
	}

	// запускаем сервер
	err = http.ListenAndServe(gophermart.Cfg.ServerAddress, r)
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

	r.Post("/api/user/register", handlers.RegisterHandler)
	r.Post("/api/user/login", handlers.LoginHandler)

	// закрытые авторизацией эндпоинты
	r.Mount("/api/user", privateRouter())

	return r
}

// privateRouter Роутер для закрытых авторизацией эндпоинтов
func privateRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(middlewares.UserAuth) // проверка авторизации

	r.Post("/orders", handlers.AddOrderHandler)
	r.Get("/orders", handlers.GetOrdersHandler)
	r.Get("/balance", handlers.BalanceHandler)
	r.Post("/balance/withdraw", handlers.WithdrawHandler)
	r.Get("/balance/withdrawals", handlers.WithdrawalsHandler)

	return r
}
