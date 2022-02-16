package main

import (
	"errors"
	"fmt"
	"io"
	"sync"
)

var (
	ErrPoolClosed = errors.New("pool has been closed")
)

// Pool资源池，让goroutine们来安全地共享资源
type Pool struct {
	factory   func() (io.Closer, error)
	resources chan io.Closer
	mtx       sync.Mutex
	closed    bool
}

// New Pool的constructor
func New(factory func() (io.Closer, error), size int) (*Pool, error) {
	if size < 0 {
		return nil, errors.New("invalid size for the resources pool")
	}

	return &Pool{
		factory:   factory,
		resources: make(chan io.Closer, size),
		closed:    false,
	}, nil
}

func (p *Pool) AcquireResource() (io.Closer, error)  {
	select {
	case resource, ok := <-p.resources:
		if !ok {
			return nil, ErrPoolClosed
		}
		fmt.Println("acquire resource from pool")
		return resource, nil
	default:
		fmt.Println("acquire new resource")
		return p.factory()
	}
}

func (p *Pool) ReleaseResource(resource io.Closer)  {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if p.closed {
		resource.Close()
		return
	}

	select {
	case p.resources <- resource:
		fmt.Println("release resource back to the pool")
	default:
		fmt.Println("release resource closed")
		resource.Close()
	}
}

func (p *Pool) Close()  {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if p.closed {
		return
	}

	p.closed = true
	// 关闭管道
	close(p.resources)
	// 关闭资源
	for resource := range p.resources {
		resource.Close()
	}
}