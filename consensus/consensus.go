package consensus

import (
	"github.com/mahmednabil109/gdeb/communication"
	"github.com/mahmednabil109/gdeb/data"
)

type Consensus struct {
	ChanNetBlock        <-chan data.Block
	ChanNetTransaction  <-chan data.Transaction
	ChanConsBlock       chan<- data.Block
	ChanConsTransaction chan<- data.Transaction
}

func New(c *communication.CommunNetwCons) *Consensus {
	return &Consensus{
		ChanNetBlock:        c.ChanNetBlock,
		ChanNetTransaction:  c.ChanNetTransaction,
		ChanConsBlock:       c.ChanConsBlock,
		ChanConsTransaction: c.ChanConsTransaction,
	}
}
