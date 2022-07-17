package sched

import "github.com/althk/goeasy/sched/internal/strategy"

type Task strategy.Task
type TaskScheduler strategy.TaskScheduler
type TaskOption strategy.TaskOption
type Priority strategy.Priority

func NewFIFO(maxTasks, maxConcurrency uint) strategy.TaskScheduler {
	return strategy.NewFIFO(maxTasks, maxConcurrency)
}

func WithPriority(p Priority) TaskOption {
	return func(opts *strategy.TaskConfig) {
		opts.Priority = strategy.Priority(p)
	}
}
