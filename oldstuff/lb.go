package main

import (
	"container/heap"
	"math/rand"
	"time"
)

type Request struct {
	fn func() int //operation to perform
	c  chan int   // the channel to return the result
}

func requester(work chan<- Request) {
	c := make(chan int)
	for {
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		work <- Request{workFn, c} //send request
		result := <-c
		furtherProcess(result)
	}
}

type Worker struct {
	requests chan Request
	pending  int
	index    int //index in the heap
}

func (w *Worker) work(done chan *Worker) {
	for {
		req := <-w.requests //get request from load balancer
		req.c <- req.fn()   // call fn and send result
		done <- w
	}
}

type Pool []*Worker

func (p Pool) Push(x any) {
	//TODO implement me
	panic("implement me")
}

func (p Pool) Pop() any {
	//TODO implement me
	panic("implement me")
}

func (p Pool) Len() int {
	//TODO implement me
	panic("implement me")
}

func (p Pool) Swap(i, j int) {
	//TODO implement me
	panic("implement me")
}

type Balancer struct {
	pool Pool
	done chan *Worker
}

func (b *Balancer) balance(work chan Request) {
	for {
		select {
		case req := <-work: //received request
			b.dispatch(req) //send to worker
		case <-b.done: // finished work
			b.completed(w) //update info
		}
	}
}

func (p Pool) Less(i, j int) bool {
	return p[i].pending < p[j].pending
}

// Send Request to worker
func (b *Balancer) dispatch(req Request) {
	// Grab the least loaded worker...
	w := heap.Pop(&b.pool).(*Worker)
	// ...send it the task.
	w.requests <- req
	// One more in its work queue.
	w.pending++
	// Put it into its place on the heap.
	heap.Push(&b.pool, w)
}

// Job is complete; update heap
func (b *Balancer) completed(w *Worker) {
	// One fewer in the queue.
	w.pending--
	// Remove it from heap.
	heap.Remove(&b.pool, w.index)
	// Put it into its place on the heap.
	heap.Push(&b.pool, w)
}
