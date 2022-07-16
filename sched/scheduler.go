package sched

import "github.com/althk/goeasy/sched/internal/strategy"

type Task strategy.Task
type TaskScheduler strategy.TaskScheduler

func NewFIFO(maxTasks, maxConcurrency uint) strategy.TaskScheduler {
	return strategy.NewFIFO(maxTasks, maxConcurrency)
}
