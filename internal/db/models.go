package db

type Model interface {
	// Table название таблицы
	isModel() bool
}

type User struct {
	Persist   bool   `db:"-" json:"-"`
	Id        int    `db:"id"`
	Login     string `db:"login" json:"login"`
	Password  string `db:"password" json:"password"`
	Accrual   int    `db:"accrual"`
	Withdrawn int    `db:"withdrawn"`
	Balance   int    `db:"balance"`
}

func (u *User) isModel() bool { return true }

// NEW — заказ загружен в систему, но не попал в обработку;
// PROCESSING — вознаграждение за заказ рассчитывается;
// INVALID — система расчёта вознаграждений отказала в расчёте;
// PROCESSED — данные по заказу проверены и информация о расчёте успешно получена.
const (
	ORDER_STATUS_NEW        = "NEW"
	ORDER_STATUS_PROCESSING = "PROCESSING"
	ORDER_STATUS_INVALID    = "INVALID"
	ORDER_STATUS_PROCESSED  = "PROCESSED"
)

type Order struct {
	Persist    bool   `db:"-" json:"-"`
	OrderId    int    `db:"orderId" json:"number"`
	Status     string `db:"status" json:"status"`
	UserId     int    `db:"userId" json:"-"`
	Accrual    int    `db:"accrual" json:"accrual,omitempty"`
	Withdrawn  int    `db:"withdraw" json:"-"`
	UploadedAt string `db:"uploadedAt" json:"uploaded_at"`
}

func (o *Order) isModel() bool { return true }
