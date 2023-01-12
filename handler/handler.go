package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
	"walletsvc/entity"
	"walletsvc/usecase"
)

type Handler struct{
	AccountUseCase *usecase.Account
	WallentUsecase *usecase.Wallet
	TransactionUsecase *usecase.Transaction
	SchedulerUsecase *usecase.Scheduler
}

func Init(AccountUseCase *usecase.Account,WallentUsecase *usecase.Wallet,TransactionUsecase *usecase.Transaction,
	SchedulerUsecase *usecase.Scheduler) *Handler{
	return &Handler{
		AccountUseCase:AccountUseCase ,
		WallentUsecase: WallentUsecase,
		TransactionUsecase: TransactionUsecase,
		SchedulerUsecase: SchedulerUsecase,
	}
}

func (e *Handler) HandlerInitilizeAccount(w http.ResponseWriter, req *http.Request) {
	requestBody := entity.RequestInitAccount{}
	responsBody := entity.DataInitAccount{}
	if err := req.ParseForm(); err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(&entity.ResponseError{"fail", err.Error()})
		return
	}
	requestBody.CustomerXid=req.Form.Get("customer_xid")

	if requestBody.CustomerXid != "" {
		if e.AccountUseCase.GetUser(requestBody.CustomerXid) == false {
			e.AccountUseCase.CreateUser(requestBody.CustomerXid)
			e.WallentUsecase.CreateWallet(requestBody.CustomerXid)
		}
		token, err:=	e.AccountUseCase.GenerateToken(requestBody.CustomerXid)
		if err != nil {
			log.Println(err)
			json.NewEncoder(w).Encode(&entity.ResponseError{"fail", err.Error()})	
			return
		}
		responsBody.Token = token
	}else{
		json.NewEncoder(w).Encode(&entity.ResponseError{"fail", "user id is empty"})
		return
	}

	json.NewEncoder(w).Encode(&entity.ResponseInitAccount{"success",responsBody})
}

func (e *Handler) HandlerEnableWallet(w http.ResponseWriter, req *http.Request) {
	userID := req.Header["user"]
	wallet,err:=e.WallentUsecase.EnableWallet(userID[0])
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(&entity.ResponseError{"fail", err.Error()})	
		return
	}
	json.NewEncoder(w).Encode(&entity.ResponseGetWallet{
		"success",
		entity.DataInitGetWallet{wallet}})

}

func (e *Handler) HandlerDisableWallet(w http.ResponseWriter, req *http.Request) {
	userID := req.Header["user"]
	wallet,err:=e.WallentUsecase.DisabledWallet(userID[0])
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(&entity.ResponseError{"fail", err.Error()})	
		return
	}
	json.NewEncoder(w).Encode(&entity.ResponseGetWallet{
		"success",
		entity.DataInitGetWallet{wallet}})

}

func (e *Handler) HandlerGetWalletBalance(w http.ResponseWriter, req *http.Request) {
	userID := req.Header["user"]
	wallet,err:= e.WallentUsecase.GetWallet(userID[0])
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(&entity.ResponseError{"fail", err.Error()})	
		return
	}

	json.NewEncoder(w).Encode(&entity.ResponseGetWallet{
		"success",
		entity.DataInitGetWallet{wallet}})
}

func(e *Handler) HandlerAddVirtualMoney(w http.ResponseWriter, req *http.Request) {
	userID := req.Header["user"]
	_,err:= e.WallentUsecase.GetWallet(userID[0])
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(&entity.ResponseError{"fail", err.Error()})	
		return
	}
	if err := req.ParseForm(); err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(&entity.ResponseError{"fail", err.Error()})
		return
	}
	amount,err := strconv.ParseFloat(req.Form.Get("amount"),64)
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(&entity.ResponseError{"fail", err.Error()})
		return
	}
	deposit,err:=e.TransactionUsecase.AddVirtualMoney(userID[0],amount,req.Form.Get("reference_id"))
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(&entity.ResponseError{"fail", err.Error()})
		return
	}
	job:= entity.SchedulerUpdateBalanceJob{
		Amount: amount,
		UserID: userID[0],
		Operation: "add",
		TransactionID: deposit.Id,
		Timestamp: time.Now(),
	}
	e.SchedulerUsecase.AddJob(&job)
	json.NewEncoder(w).Encode(&entity.ResponseAddVirtualMoney{
		"success",
		entity.DataDeposit{deposit}})
}

func(e *Handler) HandlerUseVirtualMoney(w http.ResponseWriter, req *http.Request) {
	userID := req.Header["user"]
	_,err:= e.WallentUsecase.GetWallet(userID[0])
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(&entity.ResponseError{"fail", err.Error()})	
		return
	}
	if err := req.ParseForm(); err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(&entity.ResponseError{"fail", err.Error()})
		return
	}
	
	amount,err := strconv.ParseFloat(req.Form.Get("amount"),64)
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(&entity.ResponseError{"fail", err.Error()})
		return
	}
	withdrawal,err:=e.TransactionUsecase.UseVirtualMoney(userID[0],amount,req.Form.Get("reference_id"))

	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(&entity.ResponseError{"fail", err.Error()})
		return
	}

	job:= entity.SchedulerUpdateBalanceJob{
		Amount: amount,
		UserID: userID[0],
		Operation: "sub",
		TransactionID: withdrawal.Id,
		Timestamp: time.Now(),
	}
	e.SchedulerUsecase.AddJob(&job)

	json.NewEncoder(w).Encode(&entity.ResponseUseVirtualMoney{
		"success",
		entity.DataWithdrawal{withdrawal}})
}

func(e *Handler) HandlerGetTransactions(w http.ResponseWriter, req *http.Request) {
	userID := req.Header["user"]
	_,err:= e.WallentUsecase.GetWallet(userID[0])
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(&entity.ResponseError{"fail", err.Error()})	
		return
	}

	transactions := e.TransactionUsecase.GetTransactions(userID[0])
	json.NewEncoder(w).Encode(&entity.ResponseGetTransactions{
		"success",
		entity.DataTransactions{transactions}})
}