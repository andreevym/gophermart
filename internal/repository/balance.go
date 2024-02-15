package repository

// BalanceRepository
//
//go:generate mockgen -source=balance.go -destination=./mock/balance.go -package=mock
type BalanceRepository interface {
	WithdrawAmount(userID int64, orderID int64, amount float32)
	AccrualAmount(userID int64, orderID int64, amount float32)
}
