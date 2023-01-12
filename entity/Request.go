package entity


type RequestInitAccount struct{
	CustomerXid string   `json:"customer_xid"`
}
type ResponseInitAccount struct{
	Status       string `json:"status"`
	Data       	 DataInitAccount `json:"data"` 
}
type DataInitAccount struct{
	Token      string `json:"token"`
}

type ResponseError struct{
	Status       string `json:"status"`
	Error       string `json:"error"`
}

type ResponseGetWallet struct{
	Status       string `json:"status"`
	Data       	 DataInitGetWallet `json:"data"` 
}

type DataInitGetWallet struct{
	Wallet      *Wallet `json:"wallet"`
}

type RequestTransactionVirtualMoney struct{
	Amount string   `json:"amount"`
	ReferenceID string   `json:"reference_id"`
}

type ResponseAddVirtualMoney struct{
	Status       string `json:"status"`
	Data       	 DataDeposit `json:"data"`
}

type DataDeposit struct{
	Deposit      *Deposit `json:"deposit"`
}

type ResponseUseVirtualMoney struct{
	Status       string `json:"status"`
	Data       	 DataWithdrawal `json:"data"`
}

type DataWithdrawal struct{
	Withdrawal      *Withdrawal `json:"withdrawal"`
}

type ResponseGetTransactions struct{
	Status       string `json:"status"`
	Data       	 DataTransactions `json:"data"`
}

type DataTransactions struct{
	Transactions      *Transactions `json:"transactions"`
}