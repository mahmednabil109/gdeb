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
	userMoney           *blockchain.UserMoney
	blockchain          *blockchain.Blockchain
	transPool           *blockchain.TransPool
}

func New(c *communication.CommunNetwCons, stakeDist *blockchain.StakeDistribution) *Consensus {
	return &Consensus{ChanNetBlock: c.ChanNetBlock, ChanNetTransaction: c.ChanNetTransaction,
		ChanConsBlock:       c.ChanConsBlock,
		ChanConsTransaction: c.ChanConsTransaction,
		stakeDist:           stakeDist,
		userMoney:           blockchain.NewUserMoney(stakeDist),
		blockchain:          blockchain.NewBlockchain(),
		transPool:           blockchain.NewTransPool(),
	}
}

func (c *Consensus) HasEnoughMoney(user string, amount float64) bool {
	return c.userMoney.HasEnoughMoney(user, amount)
}

func (c *Consensus) validTrans(t data.Transaction) bool {
	if !t.Validate() {
		return false
	}
	valid := c.userMoney.HasEnoughMoney(t.From, float64(t.Amount))
	return valid
}

func (c *Consensus) updateChain(b data.Block) {
	c.blockchain.Update(b)
	//remove all transactions contained in the block that got appended
	c.transPool.Update(b.Transactions)
}

//note:
// 1) transaction could have been invalidated after previous block added it
// 2) several transaction could come from 1 uesr, enough money for each one seperately, but not for their some (can be invalid if validating 1 by 1)
//    --> combining transactions of same sender beforhand (looping and validating), using hashmap
func (c *Consensus) validBlock(b data.Block) bool {
	leader := ValidateLeader(b.Nonce, b.SlotLeader, b.VrfProof, c.stakeDist)
	if !leader {
		return false
	}
	for _, t := range b.Transactions {
		if !c.validTrans(t) {
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
					if c.validBlock(b) {
						c.updateChain(b)
					}
				}()
			case t := <-c.ChanNetTransaction:
				go func() {
					//validate transaction & append if valid (fields & money in blockchain)
					if c.validTrans(t) {
						c.transPool.Add(t)
					}
				}()
			}
		}
	}()
	return nil
}
