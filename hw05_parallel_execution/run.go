package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run запускает задачи параллельно в n горутинах.
// Останавливается при достижении m ошибок (если m > 0).
// Возвращает ошибку, если лимит ошибок превышен.
func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return errors.New("кол-во воркеров должно быть больше 0")
	}

	var (
		errorsInFirstM int32 // Ошибки в первых M задачах.
		totalTasks     int32 // Всего выполненных задач.
		errorsCount    int32 // Всего ошибок.
		stop           int32 // Флаг остановки.
	)

	// Канал для задач.
	taskChan := make(chan Task, len(tasks))

	// Канал завершения, ничем не заполняем.
	doneChan := make(chan struct{})

	var closeOnce sync.Once

	// Заполняем канал задачами.
	for _, task := range tasks {
		taskChan <- task
	}

	close(taskChan)

	var wg sync.WaitGroup

	// Воркер выполняет работу в горутинах.
	runWorker := func() {
		defer wg.Done()

		for {
			// Если достигли лимита ошибок в первых M задачах, то это повод остановиться.
			if atomic.LoadInt32(&stop) == 1 {
				return
			}

			select {
			case task, ok := <-taskChan:
				if !ok {
					return
				}

				currentTaskNumber := int(atomic.AddInt32(&totalTasks, 1))
				err := task()

				if err != nil {
					newErrors := atomic.AddInt32(&errorsCount, 1)

					// Ошибка в первых M задачах.
					if currentTaskNumber <= m {
						newErrorsInFirstM := atomic.AddInt32(&errorsInFirstM, 1)

						if newErrorsInFirstM >= int32(m) && newErrors >= int32(m) {
							// Достигли лимита ошибок в первых M задачах и всего ошибок.
							closeOnce.Do(func() { close(doneChan) })
							atomic.StoreInt32(&stop, 1)
							return
						}
					} else if errorsCount >= int32(m) && m > 0 {
						closeOnce.Do(func() { close(doneChan) })
						atomic.StoreInt32(&stop, 1)
						return
					}
				}

			case <-doneChan:
				return
			}

			if atomic.LoadInt32(&stop) == 1 {
				return
			}
		}
	}

	wg.Add(n)

	for i := 0; i < n; i++ {
		go runWorker()
	}

	wg.Wait()

	// Если достигли максимально допустимого кол-ва ошибок в горутинах.
	if atomic.LoadInt32(&errorsCount) >= int32(m) && m > 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}
