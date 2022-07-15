package db

import (
	"log"

	"github.com/nastradamus39/gophermart/gophermart"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var schema = `
CREATE TABLE IF NOT EXISTS "orders" (
    "orderId" character varying(64) NOT NULL,
    "userId" integer NOT NULL,
    "accrual" double precision DEFAULT '0',
    "status" character varying(10) NOT NULL,
    "uploadedAt" timestamp DEFAULT now(),
    CONSTRAINT "orders_orderid_uindex" UNIQUE ("orderId"),
    CONSTRAINT "orders_pk" PRIMARY KEY ("orderId")
) WITH (oids = false);

CREATE SEQUENCE IF NOT EXISTS user_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE IF NOT EXISTS users (
    "login" character varying(256) NOT NULL,
    "password" character varying(256) NOT NULL,
    "balance" double precision DEFAULT '0',
    "id" integer DEFAULT nextval('user_id_seq') NOT NULL,
    CONSTRAINT "users_id_uindex" UNIQUE ("id"),
    CONSTRAINT "users_login_uindex" UNIQUE ("login"),
    CONSTRAINT "users_pk" PRIMARY KEY ("id")
) WITH (oids = false);

CREATE TABLE IF NOT EXISTS withdrawals (
    "orderId" character varying(64),
    "withdraw" double precision,
    "date" timestamp DEFAULT now(),
    "userId" integer NOT NULL
) WITH (oids = false);
`

type repositories struct {
	Users    *UsersRepository
	Orders   *OrderRepository
	Withdraw *WithdrawRepository
}

var repos *repositories

// InitDB инициализирует подключение к бд
func InitDB() (err error) {
	if gophermart.DB, err = sqlx.Open("postgres", gophermart.Cfg.DatabaseDsn); err == nil {
		err = migrate()
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		log.Println(err)
	}

	repos = &repositories{
		Users: &UsersRepository{repo{
			table: "users",
			db:    gophermart.DB,
		}},
		Orders: &OrderRepository{repo{
			table: "orders",
			db:    gophermart.DB,
		}},
		Withdraw: &WithdrawRepository{repo{
			table: "withdrawals",
			db:    gophermart.DB,
		}},
	}

	return
}

// Repositories Возвращает список всех доступных репозиториев
func Repositories() *repositories {
	return repos
}

// migrate Создает структуру базы
func migrate() (err error) {
	// миграции
	_, err = gophermart.DB.Exec(schema)
	if err != nil {
		log.Println(err)
	}

	return
}
