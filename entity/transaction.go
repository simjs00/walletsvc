package entity

import "time"

type Deposit struct{
	Id string `json:"id"`
	DepositedBy string `json:"deposited_by"`
	Status string `json:"status"`
	DepositedAt time.Time `json:"deposited_at"`
	Amount float64 `json:"amount"`
	ReferenceID string `json:"reference_id"`
}

type Withdrawal struct{
	Id string `json:"id"`
	WithdrawnBy string `json:"withdrawn_by"`
	Status string `json:"status"`
	WithdrawnAt time.Time `json:"withdrawn_at"`
	Amount float64 `json:"amount"`
	ReferenceID string `json:"reference_id"`
}

type Transactions struct{
	Deposits []Deposit `json:"deposit"`
	Withdrawal []Withdrawal `json:"withdrawal"`	
}