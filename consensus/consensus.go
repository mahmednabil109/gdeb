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
	stakeDist           *blockchain.StakeDistribution
	blockchain          *blockchain.Blockchain
	transPool           *blockchain.TransPool
}

func New(c *communication.CommunNetwCons, stakeDist *blockchain.StakeDistribution) *Consensus {
	return &Consensus{ChanNetBlock: c.ChanNetBlock, ChanNetTransaction: c.ChanNetTransaction,
		ChanConsBlock:       c.ChanConsBlock,
		ChanConsTransaction: c.ChanConsTransaction,
		stakeDist:           stakeDist,
		blockchain:          blockchain.NewBlockchain(),
		transPool:           blockchain.NewTransPool(),
	}
}

func (c *Consensus) enoughMoney(t data.Transaction) bool {
	valid, _ := c.stakeDist.Valid(t)
	return valid
}
func (c *Consensus) appendBlockchain(b data.Block) {
	c.blockchain.Add(b)
	//remove all transactions contained in the block that got appended
	c.transPool.Update(b.Transactions)
}
func (c *Consensus) updateStakedist(t data.Transaction) {
	c.stakeDist.Update(t)
}

func (c *Consensus) ValidateBlock(b data.Block) bool {

	leader := ValidateLeader(b.Nonce, b.SlotLeader, b.VrfProof, c.stakeDist)
	if !leader {
		return false
	}

	//note: transaction could have been invalidated after previous block added it
	//we need to check if enough money exists (infering blokchain, but continously changing stakeDist for fast access)
	//TODO:
	//In this case it could be that, one sender's has several transactions in block such that their combined sum is above  value saved in stakeDistribution
	//thus, this validation naive and incorrect --> needs to be addressed by combining transactions of same sender beforhand
	for _, t := range b.Transactions {
		if !t.Validate() || !c.enoughMoney(t) {
			return false
		}
	}
	return true
}

func (c *Consensus) Init() error {
	go func() {
		// recieve from network
		for {
			select {
			case b := <-c.ChanNetBlock:
				go func() {
					//validate block & append if valid (leader and fields correct)
					if c.ValidateBlock(b) {
						c.appendBlockchain(b)
					}
				}()
			case t := <-c.ChanNetTransaction:
				go func() {
					//validate transaction & append if valid (fields & money in blockchain)
					if t.Validate() && c.enoughMoney(t) {
						c.transPool.Add(t)
						c.stakeDist.Update(t)
					}
				}()
			}
		}
	}()
	return nil
}
