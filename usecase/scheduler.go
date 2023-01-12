package usecase

import (
	"walletsvc/entity"
)
type Scheduler struct{
	jobDatabase []*entity.SchedulerUpdateBalanceJob
}

func InitSchedulerUsecase(jobDatabase []*entity.SchedulerUpdateBalanceJob) *Scheduler{
	return &Scheduler{
		jobDatabase: jobDatabase,
	}
}

func(s *Scheduler) AddJob(job *entity.SchedulerUpdateBalanceJob){
	s.jobDatabase = append(s.jobDatabase, job)
}

func(s *Scheduler) GetFirstJob() *entity.SchedulerUpdateBalanceJob{
	if len(s.jobDatabase) >0 {
		firstJob := s.jobDatabase[0]
		if  len(s.jobDatabase) >=1{
			s.jobDatabase = s.jobDatabase[1:]
		}else{
			s.jobDatabase = make([]*entity.SchedulerUpdateBalanceJob, 0)
		}
	
		return firstJob
	}
	return nil
}