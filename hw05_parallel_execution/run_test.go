package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})
	t.Run("concurrency without Sleep", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		taskStarted := make(chan struct{}, tasksCount)

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				atomic.AddInt32(&runTasksCount, 1)
				taskStarted <- struct{}{}
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		var wg sync.WaitGroup
		wg.Add(1)

		// Запускаем функцию в горутине для параллельного отслеживания в require.Eventually
		go func() {
			defer wg.Done()
			err := Run(tasks, workersCount, maxErrorsCount)
			require.NoError(t, err)
			require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		}()

		// Переменная для отслеживания параллельных задач
		parallelTasks := 0
		checkTasksStarted := func() bool {
			// Проверяем, сколько задач начали выполняться параллельно
			for {
				select {
				case <-taskStarted:
					parallelTasks++
				default:
					// Выходим из цикла, если нет новых задач в канале
					return parallelTasks >= workersCount
				}
			}
		}

		// Ожидаем, что параллельно выполнится как минимум количество воркеров.
		require.Eventually(t, checkTasksStarted, 100*time.Millisecond, 10*time.Millisecond, "tasks did not run concurrently")

		// Закрываем канал после завершения всех задач
		close(taskStarted)
		// Ожидаем завершения всех горутин.
		wg.Wait()
	})

	t.Run("negative error's limit", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := -3

		err := Run(tasks, workersCount, maxErrorsCount)
		require.NoError(t, err)
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
	})
}
