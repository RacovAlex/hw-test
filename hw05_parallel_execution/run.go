package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	// Игнорируем ошибки в случае передачи отрицательного или нулевого лимита ошибок
	if m <= 0 {
		m = len(tasks) + 1
	}
	// Не запускаем лишние горутины, если количество задач меньше количества горутин
	if len(tasks) < n {
		n = len(tasks)
	}

	// Устанавливаем канал для завершения воркеров
	done := make(chan struct{})
	// Канал для получения результатов работы воркеров
	results := make(chan error, m)
	// Канал для принятия Tasks из массива
	work := make(chan Task, len(tasks))

	wg := &sync.WaitGroup{}

	// Счетчик ошибок
	var errCount int32

	for i := 0; i < n; i++ {
		wg.Add(1)

		// Запускаем воркер в горутине
		go func() {
			defer wg.Done()
			// Используем паттерн concurrency для работы с Tasks
			for {
				select {
				case <-done:
					return
				default:
				}

				select {
				case task := <-work:
					results <- task()
				default:
					time.Sleep(time.Millisecond * 10)
				}
			}
		}()
	}

	for _, task := range tasks {
		work <- task
	}
	// Закрывать канал work не требуется, так как имеется механизм
	// завершения воркеров через закрытие канала done. Если же канал будет
	// закрываться - требуется проверка на наличия Task при чтении из него.

	go func() {
		// Устанавливаем счетчик выполненных заданий для проверки на необходимость завершения
		// работы через закрытие канала done
		doneTaskCount := 0
		var once sync.Once

		for err := range results {
			doneTaskCount++
			// Безопасный доступ нужен, так как errCount читается из основной горутины
			if err != nil {
				atomic.AddInt32(&errCount, 1)
			}
			// Если достигнут лимит ошибок или выполнены все задания - отправляем сигнал на завершение
			if m == int(atomic.LoadInt32(&errCount)) || doneTaskCount == len(tasks) {
				// Закрываем канал через once, так как данный цикл не прекращается
				once.Do(func() {
					close(done)
				})
			}
		}
	}()

	// Ожидаем завершение всех горутин
	wg.Wait()
	close(results)

	if int(atomic.LoadInt32(&errCount)) >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
