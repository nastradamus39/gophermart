package db

import (
	"log"

	"github.com/nastradamus39/gophermart/gophermart"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var schema = `
create table if not exists users 
(
    login     varchar(256) not null,
    password  varchar(256) not null,
    accrual   int default 0,
    withdrawn int default 0,
    balance   int default 0
);

create unique index if not exists users_login_uindex
    on users (login);

--alter table users
--    add constraint users_pk
--        primary key (login);

create table if not exists orders
(
    "orderId"  varchar(64) not null,
    "status"   varchar(10) not null,
    "userId"   int not null,
    accural    int default 0,
    withdraw   int default 0,
    uploadedAt timestamp default now()
);

create unique index if not exists orders_orderid_uindex
    on orders ("orderId");

--alter table orders
--    add constraint orders_pk
--        primary key ("orderId");
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
