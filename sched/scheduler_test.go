package sched

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/althk/goeasy/sched/internal/strategy"
	"github.com/stretchr/testify/require"
)

var s strategy.TaskScheduler
var testTask strategy.Task
var c uint32
var mu sync.Mutex

func TestFIFOScheduler_Init(t *testing.T) {
	s = NewFIFO(5, 1)
	v := s.IsReady()
	require.False(t, v, "IsReady(): want: false, got: %v", v)
	// Initialize and ensure no error was returned
	require.NoError(t, s.Init(), "FIFO Scheduler init returned error; want no error")
	v = s.IsReady()
	require.True(t, v, "IsReady(): want: true, got: %v", v)
	require.NoError(t, s.Shutdown(true))
}

func newFIFOWithTasks(t *testing.T) strategy.TaskScheduler {
	s = NewFIFO(5, 1)
	s.Init()
	// Enqueue 4 tasks, verify by the counter.
	for i := 0; i < 4; i++ {
		err := s.EnqueueTask(testTask)
		require.NoError(t, err)
	}
	return s
}
func TestFIFOScheduler_WaitForAllTasks(t *testing.T) {

	s = newFIFOWithTasks(t)
	err := s.Shutdown(true)
	require.NoError(t, err)

	mu.Lock()
	require.EqualValues(t, 4, c) // scheduler should wait for all tasks to complete and then shutdown.
	c = 0
	mu.Unlock()
}

func TestFIFOScheduler_NoWait(t *testing.T) {
	s = newFIFOWithTasks(t)
	err := s.Shutdown(false)
	require.NoError(t, err)

	mu.Lock()
	// scheduler should close down immediately so one or more enqueued tasks will not run.
	require.Less(t, c, uint32(4))
	mu.Unlock()

}

func TestMain(m *testing.M) {
	testTask = func() {
		fmt.Println("running..")
		mu.Lock()
		atomic.AddUint32(&c, 1)
		mu.Unlock()
		time.Sleep(time.Duration(rand.Intn(100) * int(time.Millisecond)))
	}

}
