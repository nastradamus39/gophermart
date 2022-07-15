package gophermart

import (
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
)

// Конфигурирование сервиса накопительной системы лояльности
// Сервис должен поддерживать конфигурирование следующими методами:
// адрес и порт запуска сервиса: переменная окружения ОС RUN_ADDRESS или флаг -a;
// адрес подключения к базе данных: переменная окружения ОС DATABASE_URI или флаг -d;
// адрес системы расчёта начислений: переменная окружения ОС ACCRUAL_SYSTEM_ADDRESS или флаг -r.

const (
	SessionName               = "gopherMarketSid"
	ContextUserKey ContextKey = iota
)

type ContextKey int8

// Cfg конфиг приложения
var Cfg Config

// DB подключение к базе
var DB *sqlx.DB

// SessionStore хранилище сессий
var SessionStore = sessions.NewCookieStore([]byte("aaa"))

// Config конфиг приложения
type Config struct {
	ServerAddress  string `env:"RUN_ADDRESS" envDefault:"localhost:8080"`
	DatabaseDsn    string `env:"DATABASE_URI" envDefault:""`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:""`
}
