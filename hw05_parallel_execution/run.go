package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	// Игнорируем ошибки в случае передачи отрицательного или нулевого лимита ошибок
	if m <= 0 {
		m = len(tasks) + 1
	}

	// Устанавливаем флаг и счетчик
	var stopFlag int32
	counter := 0

	ch := make(chan Task, len(tasks))
	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	for i := 0; i < n; i++ {
		wg.Add(1)

		// Запускаем воркер в горутине
		go func() {
			defer wg.Done()
			for task := range ch {
				// Проверяем флаг остановки перед выполнением задачи
				// и прекращаем выполнение горутины при наличии флага
				if atomic.LoadInt32(&stopFlag) == 1 {
					return
				}

				// Выполняем задачу
				err := task()
				// Если задача завершилась ошибкой, увеличиваем счётчик ошибок
				if err != nil {
					mu.Lock()
					counter++
					// Если счетчик ошибок превысил лимит ошибок
					// устанавливаем флаг остановки
					if counter >= m {
						atomic.StoreInt32(&stopFlag, 1)
					}
					mu.Unlock()
				}
			}
		}()
	}

	for _, task := range tasks {
		ch <- task
	}
	close(ch)

	wg.Wait()

	// После завершения всех горутин возвращаем ошибку при превышении лимита
	if counter >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
