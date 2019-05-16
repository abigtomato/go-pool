package pool

import "log"

type Job func()

type Worker struct {
	pool *Pool     // 所属的池
	task chan Job  // 任务队列
	quit chan bool // 退出标记
}

func NewWorker(pool *Pool) *Worker {
	return &Worker{
		pool: pool,
		task: make(chan Job, chanSize),
		quit: make(chan bool, 0),
	}
}

func (w *Worker) start() {
	for {
		select {
		case job := <-w.task:
			job()

			w.pool.workerQueue <- w

			if p := recover(); p != nil {
				if w.pool.PanicHandler != nil {
					w.pool.PanicHandler(p)
				} else {
					log.Printf("worker exits from a panic: %v", p)
				}
			}
		case <-w.quit:
			return
		}
	}
}

func (w *Worker) stop() {
	go func() {
		w.quit <- true
	}()
}
