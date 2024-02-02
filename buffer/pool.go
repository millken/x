package buffer

import (
	"github.com/millken/gosync"
)

var (
	_pool = NewPool()
	// Get retrieves a buffer from the pool, creating one if necessary.
	Get = _pool.Get
)

// A Pool is a type-safe wrapper around a sync.Pool.
type Pool struct {
	p *gosync.Pool[*Buffer]
}

// NewPool constructs a new Pool.
func NewPool() Pool {
	return Pool{
		p: gosync.NewPool(func() *Buffer {
			return &Buffer{bs: make([]byte, 0, _size)}
		}),
	}
}

// NewPoolSize constructs a new Pool.
func NewPoolSize(size int) Pool {
	return Pool{
		p: gosync.NewPool(func() *Buffer {
			return &Buffer{bs: make([]byte, 0, size)}
		}),
	}
}

// Get retrieves a Buffer from the pool, creating one if necessary.
func (p Pool) Get() *Buffer {
	buf := p.p.Get()
	buf.Reset()
	buf.pool = p
	return buf
}

func (p Pool) put(buf *Buffer) {
	p.p.Put(buf)
}
