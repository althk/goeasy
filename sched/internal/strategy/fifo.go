package strategy

import (
	"fmt"
	"sync"
)

// FIFOSched implements TaskScheduler and runs
// tasks as they come in.
//
// The scheduler is threadsafe, and will block
// the callers by default if maxTasks are enqueued or
// if maxCurrency # of tasks are running at the moment.
type FIFOSched struct {
	maxTasks       uint
	maxConcurrency uint
	taskQ          chan Task
	workerPool     chan int
	quitCh         chan bool
	ready          bool
	mu             sync.RWMutex
	wg             sync.WaitGroup
}

var _ TaskScheduler = new(FIFOSched)

func (s *FIFOSched) Init() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.ready {
		return fmt.Errorf("FIFO Scheduler already initialized")
	}
	for i := 0; i < int(s.maxConcurrency); i++ {
		go func(i int) {
			for {
				select {
				case task := <-s.taskQ:
					fmt.Printf("(worker %d) running task\n", i)
					task()
					s.wg.Done()
					fmt.Printf("(worker %d) completed task\n", i)
				case <-s.quitCh:
					fmt.Printf("(worker %d) shutting down", i)
					return
				}
			}
		}(i)
	}
	s.ready = true
	fmt.Println("sched ready")
	return nil
}

func (s *FIFOSched) EnqueueTask(t Task) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if !s.ready {
		return fmt.Errorf("FIFO Scheduler not initialized, please call Init() before enqueueing")
	}
	fmt.Println("task enqueued")
	s.wg.Add(1)
	s.taskQ <- t
	return nil
}

func (s *FIFOSched) Shutdown(waitForAllTasks bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	fmt.Println("shutting down scheduler")
	if waitForAllTasks {
		s.wg.Wait()
	}
	s.quitCh <- true
	s.ready = false
	fmt.Println("scheduler shut down completed")
	return nil
}

func (s *FIFOSched) IsReady() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.ready
}

func NewFIFO(maxTasks, maxConcurrency uint) TaskScheduler {
	return &FIFOSched{
		maxTasks:       maxTasks,
		maxConcurrency: maxConcurrency,
		taskQ:          make(chan Task, maxTasks),
		workerPool:     make(chan int, maxConcurrency),
		quitCh:         make(chan bool),
	}
}
