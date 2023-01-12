package entity

import "time"

type SchedulerUpdateBalanceJob struct{
	Amount float64
	UserID  string
	Operation string
	TransactionID string
	Timestamp time.Time
}