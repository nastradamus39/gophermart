package db

import "time"

// OrderStatusNew
// NEW — заказ загружен в систему, но не попал в обработку;
// PROCESSING — вознаграждение за заказ рассчитывается;
// INVALID — система расчёта вознаграждений отказала в расчёте;
// PROCESSED — данные по заказу проверены и информация о расчёте успешно получена.
const (
	OrderStatusNew = "NEW"
	//ORDER_STATUS_PROCESSING = "PROCESSING"
	//ORDER_STATUS_INVALID    = "INVALID"
	//ORDER_STATUS_PROCESSED  = "PROCESSED"
)

type User struct {
	Persist  bool    `db:"-" json:"-"`
	ID       int     `db:"id"`
	Login    string  `db:"login" json:"login"`
	Password string  `db:"password" json:"password"`
	Balance  float32 `db:"balance"`
}

type Order struct {
	Persist    bool      `db:"-" json:"-"`
	OrderID    string    `db:"orderId" json:"number"`
	Status     string    `db:"status" json:"status"`
	UserId     int       `db:"userId" json:"-"`
	Accrual    float32   `db:"accrual" json:"accrual,omitempty"`
	UploadedAt time.Time `db:"uploadedAt" json:"uploaded_at"`
}

type Withdraw struct {
	Persist     bool      `db:"-" json:"-"`
	Order       string    `db:"orderId" json:"order"`
	UserID      int       `db:"userId" json:"-"`
	Sum         float32   `db:"withdraw" json:"sum"`
	ProcessedAt time.Time `db:"date" json:"processed_at"`
}
