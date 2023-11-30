package scheduler

import (
	"context"

	"github.com/robfig/cron/v3"
)

type Job interface {
	GetFrequency() string
	Execute()
}

type scheduler struct {
	crontab     *cron.Cron
	backlogJobs map[string]Job
	jobs        map[string]cron.EntryID
}

var Scheduler *scheduler

func init() {
	Scheduler = NewScheduler(map[string]Job{})
}

func NewScheduler(backlogJobs map[string]Job) *scheduler {
	return &scheduler{
		crontab:     cron.New(cron.WithSeconds()),
		backlogJobs: backlogJobs,
		jobs:        map[string]cron.EntryID{},
	}
}

func (s *scheduler) BacklogJobs(backlogJobs map[string]Job) {
	for name, backlogJob := range backlogJobs {
		s.backlogJobs[name] = backlogJob
	}
}

func (s *scheduler) RemoveJobs(backlogJobs []string) {
	for _, name := range backlogJobs {
		delete(s.backlogJobs, name)
	}
}

func (s *scheduler) ClearBacklogJobs() {
	s.backlogJobs = map[string]Job{}
}

func (s *scheduler) AddJob(name string, job Job) error {
	id, err := s.crontab.AddFunc(job.GetFrequency(), job.Execute)
	if err != nil {
		return err
	}
	s.jobs[name] = id

	return nil
}

func (s *scheduler) RemoveJob(name string) {
	s.crontab.Remove(s.jobs[name])
	delete(s.jobs, name)
}

func (s *scheduler) Start() {
	for name, backlogJob := range s.backlogJobs {
		if err := s.AddJob(name, backlogJob); err != nil {
			panic(err)
		}
	}
	s.backlogJobs = map[string]Job{}
	s.crontab.Start()
}

func (s *scheduler) Stop() context.Context {
	return s.crontab.Stop()
}
