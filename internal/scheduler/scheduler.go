package scheduler

import (
	"github.com/chan-jui-huang/go-backend-package/pkg/booter/service"
	"github.com/chan-jui-huang/go-backend-package/pkg/scheduler"
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
