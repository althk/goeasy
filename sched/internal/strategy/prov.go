package strategy

// Task is the function that the scheduler needs to run.
// It has to be thread-safe.
type Task func()

// TaskOptions captures options for the current task
type TaskOption func(opts *TaskConfig)

type TaskConfig struct {
	Priority Priority
}

// TaskScheduler is the interface that wraps
// task scheduling API.
//
// Implementations of different scheduling
// algorithms are exported via this interface.
//
// Implementations of TaskScheduler API need to gaurantee thread-safety.
type TaskScheduler interface {
	// Init initializes the scheduler and brings it
	// to a ready state.
	//
	// It is an error to reinitialize an already
	// initialized scheduler.
	Init() error

	// EnqueueTask enqueues the given task onto the
	// `runqueue`.
	//
	// The task execution may vary between implementations
	// and in some cases may not be gauranteed to run (again
	// depends on implementation).
	//
	// opts are one or more `TaskOption`s which are available
	// via the sched package. Some schedulers might not need
	// any options.
	EnqueueTask(t Task, opts ...TaskOption) error

	// Shutdown shuts down the scheduler.
	//
	// if `waitForAllTasks` is true, the scheduler will
	// wait to run all the tasks already enqueued and then
	// exit.
	//
	// if `waitForAllTasks` is false, the scheduler will
	// exit immediately once all the in-flight tasks are
	// completed. The remaining tasks enqueued (if any) are
	// discarded.
	Shutdown(waitForAllTasks bool) error

	// IsReady returns true if the scheduler has been initialized
	// and ready for enqueing tasks.
	IsReady() bool
}

func taskConfig(opts ...TaskOption) *TaskConfig {
	taskCfg := &TaskConfig{}
	for _, opt := range opts {
		opt(taskCfg)
	}
	return taskCfg
}
