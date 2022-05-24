package consensus

import (
	"github.com/mahmednabil109/gdeb/communication"
	"github.com/mahmednabil109/gdeb/data"
)

// transaction pool (store transactions)

type Consensus struct {
	ChanNetBlock        <-chan data.Block
	ChanNetTransaction  <-chan data.Transaction
	ChanConsBlock       chan<- data.Block
	ChanConsTransaction chan<- data.Transaction
	stakeDist           map[string]float64
}

func New(c *communication.CommunNetwCons, stakeDist map[string]float64) *Consensus {
	return &Consensus{
		ChanNetBlock:        c.ChanNetBlock,
		ChanNetTransaction:  c.ChanNetTransaction,
		ChanConsBlock:       c.ChanConsBlock,
		ChanConsTransaction: c.ChanConsTransaction,
		stakeDist:           c.stakstakeDist,
	}
}

func (c *Consensus) Init() error {
	go func() {
		// recieve from network
		for {
			select {
			case b := <-c.ChanNetBlock:
				go func() {
					//handle block
					ValidateBlock(b)
					//validate if leader correct, yes -> blockchain
				}()
			case t := <-c.ChanNetTransaction:
				go func() {
					//handle transaction
					//validate trans -> put in pool
				}()
			}
		}
	}()
}
