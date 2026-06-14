package pool

import (
	"sync"
)

// Pool 通用协程池
type Pool struct {
	workerCount int
	taskChan    chan func()
	wg          sync.WaitGroup
}

// New 创建协程池
func New(workerCount int) *Pool {
	p := &Pool{
		workerCount: workerCount,
		taskChan:    make(chan func(), workerCount*10),
	}
	p.start()
	return p
}

func (p *Pool) start() {
	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for task := range p.taskChan {
				task()
			}
		}()
	}
}

// Submit 提交任务
func (p *Pool) Submit(task func()) {
	p.taskChan <- task
}

// Wait 等待所有任务完成
func (p *Pool) Wait() {
	close(p.taskChan)
	p.wg.Wait()
}
