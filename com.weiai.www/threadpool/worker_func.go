// MIT License

// Copyright (c) 2018 Andy Pan

package ants

import (
	"log"
	"runtime"
	"time"
)

// goWorkerWithFunc is the actual executor who runs the tasks,
// it starts a goroutine that accepts tasks and
// performs function calls.
type goWorkerWithFunc struct {
	// pool who owns this worker.
	pool *PoolWithFunc

	// args is a job should be done.
	args chan interface{}

	// recycleTime will be update when putting a worker back into queue.
	recycleTime time.Time
}

// run starts a goroutine to repeat the process
// that performs the function calls.
func (w *goWorkerWithFunc) run() {
	w.pool.incRunning()
	go func() {
		defer func() {
			if p := recover(); p != nil {
				w.pool.decRunning()
				w.pool.workerCache.Put(w)
				if w.pool.panicHandler != nil {
					w.pool.panicHandler(p)
				} else {
					log.Printf("worker with func exits from a panic: %v\n", p)
					var buf [4096]byte
					n := runtime.Stack(buf[:], false)
					log.Printf("worker with func exits from panic: %s\n", string(buf[:n]))
				}
			}
		}()

		for args := range w.args {
			if args == nil {
				w.pool.decRunning()
				w.pool.workerCache.Put(w)
				return
			}
			w.pool.poolFunc(args)
			if ok := w.pool.revertWorker(w); !ok {
				break
			}
		}
	}()
}
