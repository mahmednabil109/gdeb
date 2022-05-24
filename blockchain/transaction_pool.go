package blockchain

import (
	"sync"

	"github.com/mahmednabil109/gdeb/data"
)

type TransPool struct {
	pool map[string]data.Transaction //key is transaction signature
	mux  sync.Mutex
}

func NewTransPool() TransPool {
	return TransPool{
		pool: make(map[string]data.Transaction),
	}
}

func (tp *TransPool) Add(val string, trans data.Transaction) {
	tp.mux.Lock()
	defer tp.mux.Unlock()
	tp.pool[val] = trans
}

func (tp *TransPool) Remove(val string) {
	tp.mux.Lock()
	defer tp.mux.Unlock()
	delete(tp.pool, val)
}
