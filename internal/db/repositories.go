package db

import (
	"errors"
	"fmt"

	"github.com/nastradamus39/gophermart/gophermart"

	"github.com/jmoiron/sqlx"
)

type repository interface {
	// Save сохраняет сущность в бд
	Save(user interface{}) error
	// Delete удаляет сущность из бд
	Delete(user interface{}) error
}

type repo struct {
	table string
	db    *sqlx.DB
}

type UsersRepository struct {
	repo
}

type OrderRepository struct {
	repo
}

// Save сохраняет пользователя в базе
func (r *UsersRepository) Save(user interface{}) error {
	u, ok := user.(*User)
	if !ok {
		return errors.New("unsupported type")
	}

	if !u.Persist {
		res, err := r.db.NamedQuery(`INSERT INTO users(login, password) 
			VALUES (:login, :password) on conflict (login) DO NOTHING RETURNING login`, &u)

		if err != nil {
			return err
		}

		if !res.Next() {
			return fmt.Errorf("%w", gophermart.ErrUserLoginConflict)
		}
	}

	return nil
}

// Delete удаляет пользователя в базе
func (r *UsersRepository) Delete(user interface{}) error {
	return nil
}

// Find поиск пользователя по логину
func (r *UsersRepository) Find(login string) (user *User, err error) {
	user = &User{}
	err = r.db.Get(user, "SELECT * FROM users WHERE login = $1", login)
	user.Persist = true
	return
}

// Save сохраняет пользователя в базе
func (r *OrderRepository) Save(order interface{}) error {
	o, ok := order.(*Order)
	if !ok {
		return errors.New("unsupported type")
	}

	if !o.Persist {
		res, err := r.db.NamedQuery(`INSERT INTO orders("orderId", "userId", "status", "accrual", "withdraw") 
			VALUES (:orderId, :userId, :status, :accrual, :withdraw) on conflict ("orderId") DO NOTHING 
			RETURNING "orderId", "userId", "status", "accrual", "withdraw"`, &o)

		if err != nil {
			return err
		}

		if !res.Next() {
			return fmt.Errorf("%w", gophermart.ErrOrderIdConflict)
		}
	}

	return nil
}

// Find поиск заказа по orderId
func (r *OrderRepository) Find(orderId int) (order *Order, err error) {
	err = r.db.Get(&order, `SELECT * FROM orders WHERE "orderId" = $1`, orderId)
	order.Persist = true
	return
}

// FindByUser поиск заказов по пользователю
func (r *OrderRepository) FindByUser(userId int) (orders []*Order, err error) {
	err = r.db.Select(&orders, `SELECT * FROM orders WHERE "userId" = $1`, userId)
	return
}
