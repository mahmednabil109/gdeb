package blockchain

import (
	"sync"

	"github.com/mahmednabil109/gdeb/data"
)

type TransPool struct {
	Pool map[string]*data.Transaction //key is transaction signature
	Mux  sync.Mutex
}

func NewTransPool() *TransPool {
	return &TransPool{
		Pool: make(map[string]*data.Transaction),
	}
}

func (tp *TransPool) Add(t *data.Transaction) {
	val := t.Signature
	tp.Mux.Lock()
	defer tp.Mux.Unlock()
	tp.Pool[val] = t
}

func (tp *TransPool) Remove(val string) {
	tp.Mux.Lock()
	defer tp.Mux.Unlock()
	delete(tp.Pool, val)
}

func (tp *TransPool) Update(trans []data.Transaction) {
	tp.Mux.Lock()
	defer tp.Mux.Unlock()
	for _, t := range trans {
		delete(tp.Pool, t.Signature)
	}
}

// func (tp *TransPool) TransactionsForBlock() []*data.Transaction {
// 	tp.Mux.Lock()
// 	defer tp.Mux.Unlock()

// 	counter := 0
// 	transactions := make([]*data.Transaction, transCount)
// 	moneySum := make(map[string]float64)

// 	for _, t := range tp.Pool {
// 		if counter >= transCount {
// 			break
// 		}
// 		val, ok := moneySum[t.From]
// 		if !ok && t.validTrans(&t) {
// 			transactions = append(transactions, t)
// 			moneySum[t.From] = float64(t.Amount)
// 		} else if ok {
// 			val + float64(t.Amount)

// 		}

// 	}
// 	return transactions
// }
