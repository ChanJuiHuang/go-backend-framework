package scheduler

import (
	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/service"
	"github.com/ChanJuiHuang/go-backend-framework/pkg/scheduler"
	"go.uber.org/zap"
)

func backlogJobs() {
	scheduler.Scheduler.BacklogJobs(map[string]scheduler.Job{
		// "example": job.NewExampleJob(),
	})
}

func Start() {
	backlogJobs()
	scheduler.Scheduler.Start()

	logger := service.Registry.Get("logger").(*zap.Logger)
	logger.Info("scheduler is started")
}

func Stop() {
	<-scheduler.Scheduler.Stop().Done()

	logger := service.Registry.Get("logger").(*zap.Logger)
	logger.Info("scheduler is stopped")
}
