package gophermart

import "errors"

var ErrUserLoginConflict = errors.New(`данный логин уже занят`)

var ErrOrderIDConflict = errors.New(`данный заказ уже загружен`)
