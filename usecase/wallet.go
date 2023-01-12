package usecase

import (
	"fmt"
	"time"
	"walletsvc/entity"
	"github.com/google/uuid"
)

type Wallet struct{
	walletDatabase map[string] *entity.Wallet
}

func InitWalletUsecase(WalletData map[string] *entity.Wallet) *Wallet{
	return &Wallet{
		walletDatabase: WalletData,
	}
}
func(w *Wallet) EnableWallet(userId string) (*entity.Wallet,error){
	if w.walletDatabase[userId] != nil{
		if w.walletDatabase[userId].Status == "enabled"{
			return nil,fmt.Errorf("wallet is enabled already")
		}
		w.walletDatabase[userId].Status = "enabled"
		w.walletDatabase[userId].DisabledAt = nil
		currentTime :=  time.Now().String()
		w.walletDatabase[userId].EnabledAt = &currentTime
	}else{
		return nil,fmt.Errorf("user not found")
	}
	return w.walletDatabase[userId],nil
}


func(w *Wallet) CreateWallet(userId string) {
	id :=  uuid.New()
	w.walletDatabase[userId] = &entity.Wallet{
		Id: id.String(),
		OwnedBy: userId,
		Status: "disabled",
		Balance: 0,
	}

}

func(w *Wallet) DisabledWallet(userId string) (*entity.Wallet,error){
	if w.walletDatabase[userId] != nil{
		w.walletDatabase[userId].Status = "disabled"
		currentTime :=  time.Now().String()
		w.walletDatabase[userId].EnabledAt = nil
		w.walletDatabase[userId].DisabledAt = &currentTime
	}else{
		return nil,fmt.Errorf("user not found")
	}
	return w.walletDatabase[userId],nil
}

func(w *Wallet) GetWallet(userId string) (*entity.Wallet,error){
	if (w.walletDatabase[userId].Status == "enabled"){
		return w.walletDatabase[userId],nil
	}

	return nil,fmt.Errorf("wallet disabled")
}

func (w *Wallet) UpdateWalletBalance (userId string, amount float64, operation string) (error) {
	if (w.walletDatabase[userId].Status == "enabled"){
		if operation == "add"{
			w.walletDatabase[userId].Balance += amount
		}else{
			if w.walletDatabase[userId].Balance == 0 || (w.walletDatabase[userId].Balance-amount) < 0 {
				return fmt.Errorf("balance not sufficient")
			}
			w.walletDatabase[userId].Balance -= amount
		}
		
		return nil
	}
	return fmt.Errorf("wallet disabled")	
}
