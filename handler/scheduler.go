package handler

import (
	"log"
	"time"
	"walletsvc/usecase"
)
type SchedulerHandler struct{
	WallentUsecase *usecase.Wallet
	SchedulerUsecase *usecase.Scheduler
	UpdateAfterSecond int64
}

func InitScheduler(	WallentUsecase *usecase.Wallet,
	SchedulerUsecase *usecase.Scheduler,
	UpdateAfterSecond int64) *SchedulerHandler{
		return &SchedulerHandler{
			WallentUsecase:WallentUsecase,
			SchedulerUsecase: SchedulerUsecase,
			UpdateAfterSecond: UpdateAfterSecond,
		}	

}
func(s *SchedulerHandler) Start(){
	go func ()  {
		for (true){
			job:=s.SchedulerUsecase.GetFirstJob()
			if job!=nil{
				log.Println("update balance transaction id",job.TransactionID," user id :",job.UserID)
				err:=s.WallentUsecase.UpdateWalletBalance(job.UserID,job.Amount,job.Operation)
				if err!=nil{
					log.Println(err)
				}
			}

			time.Sleep(time.Duration(s.UpdateAfterSecond) * time.Second)
		}
	}()
}