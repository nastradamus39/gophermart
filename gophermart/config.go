package gophermart

// Конфигурирование сервиса накопительной системы лояльности
// Сервис должен поддерживать конфигурирование следующими методами:
// адрес и порт запуска сервиса: переменная окружения ОС RUN_ADDRESS или флаг -a;
// адрес подключения к базе данных: переменная окружения ОС DATABASE_URI или флаг -d;
// адрес системы расчёта начислений: переменная окружения ОС ACCRUAL_SYSTEM_ADDRESS или флаг -r.

// Cfg конфиг приложения
var Cfg Config

// Config конфиг приложения
type Config struct {
	ServerAddress  string `env:"RUN_ADDRESS" envDefault:"localhost:8080"`
	DatabaseDsn    string `env:"DATABASE_URI" envDefault:""`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:""`
}
