package usecase

import (
	"fmt"
	"time"
	"walletsvc/entity"

	"github.com/google/uuid"
)

type Transaction struct{
	depositDatabase map[string] []entity.Deposit
	withdrawalDatabase map[string] []entity.Withdrawal
	referenceIDDatabase map[string] bool
}
func InitTransactionUsecase(depositDatabase map[string] []entity.Deposit,withdrawalDatabase map[string] []entity.Withdrawal,referenceIDDatabase map[string] bool ) *Transaction{
	return &Transaction{
		depositDatabase: depositDatabase,
		withdrawalDatabase:withdrawalDatabase,
		referenceIDDatabase:referenceIDDatabase,
	}
}

func(t *Transaction) AddVirtualMoney(userID string, amount float64, referenceID string) (*entity.Deposit,error){
	if t.referenceIDDatabase[referenceID]  {
		return nil,fmt.Errorf("reference id existed")	
	}
	t.referenceIDDatabase[referenceID] =true
	id :=  uuid.New()
	newDeposit := &entity.Deposit{
		Id:id.String(),
		DepositedBy: userID,
		Status: "success",
		DepositedAt: time.Now(),
		Amount: amount,
		ReferenceID: referenceID,
	}
	t.depositDatabase[userID] = append(t.depositDatabase[userID], *newDeposit)
	return newDeposit,nil
}

func(t *Transaction) UseVirtualMoney(userID string, amount float64, referenceID string) (*entity.Withdrawal,error){
	if t.referenceIDDatabase[referenceID]  {
		return nil,fmt.Errorf("reference id existed")	
	}
	t.referenceIDDatabase[referenceID] =true
	id :=  uuid.New()
	newWithdrawal := &entity.Withdrawal{
		Id:id.String(),
		WithdrawnBy: userID,
		Status: "success",
		WithdrawnAt: time.Now(),
		Amount: amount,
		ReferenceID: referenceID,
	}
	t.withdrawalDatabase[userID] = append(t.withdrawalDatabase[userID], *newWithdrawal)
	return newWithdrawal,nil
}


func(t *Transaction) GetTransactions(userID string) *entity.Transactions{
	deposits := t.depositDatabase[userID]
	witdrawals := t.withdrawalDatabase[userID]
	transactions := entity.Transactions{
		Deposits: deposits,
		Withdrawal: witdrawals,
	}
	return &transactions
}

