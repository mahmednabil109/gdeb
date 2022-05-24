package consensus

import (
	"github.com/mahmednabil109/gdeb/blockchain"
	"github.com/mahmednabil109/gdeb/communication"
	"github.com/mahmednabil109/gdeb/data"
)

// update stakeDist after each block for fast look-up

type Consensus struct {
	ChanNetBlock        <-chan data.Block
	ChanNetTransaction  <-chan data.Transaction
	ChanConsBlock       chan<- data.Block
	ChanConsTransaction chan<- data.Transaction
	stakeDist           map[string]float64
	blockchain          blockchain.Blockchain
	transPool           blockchain.TransPool
}

func New(c *communication.CommunNetwCons, stakeDist map[string]float64) *Consensus {
	return &Consensus{ChanNetBlock: c.ChanNetBlock, ChanNetTransaction: c.ChanNetTransaction,
		ChanConsBlock:       c.ChanConsBlock,
		ChanConsTransaction: c.ChanConsTransaction,
		stakeDist:           stakeDist,
		blockchain:          blockchain.NewBlockchain(),
		transPool:           blockchain.NewTransPool(),
	}
}

func (c *Consensus) Init() error {
	go func() {
		// recieve from network
		for {
			select {
			case b := <-c.ChanNetBlock:
				go func() {
					//validate block & append if valid (leader and fields correct)
					if ValidateBlock(&b, c.stakeDist) {
						c.blockchain.Add(b)
					}
				}()
			case t := <-c.ChanNetTransaction:
				go func() {
					//validate transaction & append if valid (fields & money in blockchain)
					if t.Validate() {
						c.transPool.Add(t.Signature, t)
					}
				}()
			}
		}
	}()
	return nil
}
