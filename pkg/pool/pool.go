package pool

import (
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
)

const CLOSED = 1

var (
	chanSize = func() int {
		// 如果GOMAXPROCS为1时，使用阻塞channel
		if runtime.GOMAXPROCS(0) == 1 {
			return 0
		}

		// 如果GOMAXPROCS大于1时，使用非阻塞channel
		return 1
	}()
	PoolClosedError      = errors.New("this pool has been closed")
	InvalidPoolSizeError = errors.New("invalid size for pool")
)

type Pool struct {
	cap          int32        // 池容量
	closed       int32        // 关闭标记
	jobQueue     chan Job     // 任务队列
	workerQueue  chan *Worker // worker队列
	once         sync.Once
	PanicHandler func(interface{})
}

func NewPool(size int) (*Pool, error) {
	if size <= 0 {
		return nil, InvalidPoolSizeError
	}

	pool := &Pool{
		cap:      int32(size),
		jobQueue: make(chan Job, chanSize),
	}

	if chanSize != 0 {
		pool.workerQueue = make(chan *Worker, size)
	} else {
		pool.workerQueue = make(chan *Worker, chanSize)
	}

	return pool, nil
}

func initPool(size int) *Pool {
	pool, _ := NewPool(size)
	pool.run()
	return pool
}

func (p *Pool) submit(job Job) error {
	if atomic.LoadInt32(&p.closed) == CLOSED {
		return PoolClosedError
	}

	p.jobQueue <- job

	return nil
}

func (p *Pool) run() {
	for i := 0; i < int(p.cap); i++ {
		worker := NewWorker(p)
		go worker.start()
		p.workerQueue <- worker
	}

	go p.scheduler()
}

func (p *Pool) scheduler() {
	for {
		select {
		case job := <-p.jobQueue:
			worker := <-p.workerQueue
			worker.task <- job
		}
	}
}

func (p *Pool) close() {
	p.once.Do(func() {
		atomic.StoreInt32(&p.closed, 1)

		close(p.jobQueue)
		close(p.workerQueue)

		p.jobQueue = nil
		p.workerQueue = nil
	})
}
