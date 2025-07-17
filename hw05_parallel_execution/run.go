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
		errorsInFirstM int64 // Ошибки в первых M задачах.
		totalTasks     int64 // Всего выполненных задач.
		errorsCount    int64 // Всего ошибок.
		stop           int64 // Флаг остановки.
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
			if atomic.LoadInt64(&stop) == 1 {
				return
			}

			select {
			case task, ok := <-taskChan:
				if !ok {
					return
				}

				currentTaskNumber := atomic.AddInt64(&totalTasks, 1)
				err := task()

				if err != nil {
					handleError(currentTaskNumber, m, &errorsInFirstM, &errorsCount, &stop, &closeOnce, doneChan)
				}

			case <-doneChan:
				return
			}

			if atomic.LoadInt64(&stop) == 1 {
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
	if atomic.LoadInt64(&errorsCount) >= int64(m) && m > 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}

// handleError - Обрабатывает ошибку внутри воркера.
func handleError(
	currentTaskNumber int64,
	m int,
	errorsInFirstM *int64,
	errorsCount *int64,
	stop *int64,
	closeOnce *sync.Once,
	doneChan chan struct{},
) {
	newErrors := atomic.AddInt64(errorsCount, 1)

	if currentTaskNumber <= int64(m) {
		newErrorsInFirstM := atomic.AddInt64(errorsInFirstM, 1)
		if newErrorsInFirstM >= int64(m) && newErrors >= int64(m) {
			closeOnce.Do(func() { close(doneChan) })
			atomic.StoreInt64(stop, 1)
			return
		}
	} else if newErrors >= int64(m) && m > 0 {
		closeOnce.Do(func() { close(doneChan) })
		atomic.StoreInt64(stop, 1)
	}
}
