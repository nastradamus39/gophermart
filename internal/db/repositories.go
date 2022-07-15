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

type WithdrawRepository struct {
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
	} else {
		_, err := r.db.NamedQuery(`UPDATE users SET "balance" = :balance
			WHERE "login" = :login`, &u)

		if err != nil {
			return err
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
		res, err := r.db.NamedQuery(`INSERT INTO orders("orderId", "userId", "status", "accrual") 
			VALUES (:orderId, :userId, :status, :accrual) on conflict ("orderId") DO NOTHING 
			RETURNING "orderId", "userId", "status", "accrual"`, &o)

		if err != nil {
			return err
		}

		o.Persist = true

		if !res.Next() {
			return fmt.Errorf("%w", gophermart.ErrOrderIDConflict)
		}
	} else {
		_, err := r.db.NamedQuery(`UPDATE orders SET "status" = :status, "accrual" = :accrual 
			WHERE "orderId" = :orderId`, o)

		if err != nil {
			return err
		}
		o.Persist = true
	}

	return nil
}

// Find поиск заказа по orderId
func (r *OrderRepository) Find(orderID string) (order *Order, err error) {
	order = &Order{}
	err = r.db.Get(order, `SELECT * FROM orders WHERE "orderId" = $1`, orderID)
	order.Persist = true
	return
}

// FindByUser поиск заказов по пользователю
func (r *OrderRepository) FindByUser(userID int) (orders []*Order, err error) {
	err = r.db.Select(&orders, `SELECT * FROM orders WHERE "userId" = $1 ORDER BY "uploadedAt"`, userID)
	return
}

// FindWithdrawalsByUser списания балов пользователя
func (r *WithdrawRepository) FindWithdrawalsByUser(userID int) (withdrawals []*Withdraw, err error) {
	err = r.db.Select(
		&withdrawals,
		`SELECT "orderId", withdraw, "date" FROM withdrawals WHERE "userId" = $1 ORDER BY "date"`,
		userID,
	)
	return
}

// WithdrawalsSumByUser сумма всех списаний пользователя
func (r *WithdrawRepository) WithdrawalsSumByUser(userID int) (sum float32, err error) {
	sum = 0
	err = r.db.Get(
		&sum,
		`SELECT CASE
						   WHEN SUM(withdraw) IS NULL
							   THEN 0
						   ELSE SUM(withdraw)
						   END
				FROM withdrawals
				WHERE "userId" = $1`,
		userID,
	)
	return
}

func (r *WithdrawRepository) Save(withdraw interface{}) error {
	w, ok := withdraw.(*Withdraw)
	if !ok {
		return errors.New("unsupported type")
	}

	if !w.Persist {
		_, err := r.db.NamedQuery(`INSERT INTO withdrawals("orderId", "userId", "withdraw") 
			VALUES (:orderId, :userId, :withdraw)`, &w)

		if err != nil {
			return err
		}

		w.Persist = true
	}

	return nil
}
