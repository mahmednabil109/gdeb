package consensus

import (
	"crypto/ed25519"
	"github.com/mahmednabil109/gdeb/blockchain"
	"github.com/mahmednabil109/gdeb/communication"
	"github.com/mahmednabil109/gdeb/data"
	"time"
)

// update stakeDist after each block for fast look-up

type Consensus struct {
	ChanNetBlock        <-chan data.Block
	ChanNetTransaction  <-chan data.Transaction
	ChanConsBlock       chan<- data.Block
	ChanConsTransaction chan<- data.Transaction
	stakeDist           *blockchain.StakeDistribution
	userMoney           *blockchain.UserMoney
	Blockchain          *blockchain.Blockchain
	transPool           *blockchain.TransPool
	epochNonce          string
	PrivateKey          ed25519.PrivateKey
}

func New(c *communication.CommunNetwCons, stakeDist *blockchain.StakeDistribution, privateKey ed25519.PrivateKey) *Consensus {
	return &Consensus{ChanNetBlock: c.ChanNetBlock, ChanNetTransaction: c.ChanNetTransaction,
		ChanConsBlock:       c.ChanConsBlock,
		ChanConsTransaction: c.ChanConsTransaction,
		stakeDist:           stakeDist,
		userMoney:           blockchain.NewUserMoney(stakeDist),
		Blockchain:          blockchain.NewBlockchain(stakeDist.Distribution),
		transPool:           blockchain.NewTransPool(),
		PrivateKey:          privateKey,
		epochNonce:          "RANDOM_NONCE",
	}
}

func (c *Consensus) HasEnoughMoney(user string, amount float64) bool {
	return c.userMoney.HasEnoughMoney(user, amount)
}

func (c *Consensus) validTrans(t *data.Transaction) bool {
	if !t.Validate() {
		return false
	}
	valid := c.userMoney.HasEnoughMoney(t.From, float64(t.Amount))
	return valid
}

func (c *Consensus) updateChain(b *data.Block) {
	c.Blockchain.Update(b)
	//remove all transactions contained in the block that got appended
	c.transPool.Update(b.Transactions)
}

//note:
// 1) transaction could have been invalidated after previous block added it
// 2) several transaction could come from 1 uesr, enough money for each one seperately, but not for their some (can be invalid if validating 1 by 1)
//    --> combining transactions of same sender beforhand (looping and validating), using hashmap
func (c *Consensus) validBlock(b *data.Block) bool {
	leader := ValidateLeader(b.Nonce, b.SlotLeader, b.VrfProof, c.stakeDist)
	if !leader {
		return false
	}
	for _, t := range b.Transactions {
		if !c.validTrans(&t) {
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
					if c.validBlock(&b) {
						c.updateChain(&b)
					}
				}()
			case t := <-c.ChanNetTransaction:
				go func() {
					//validate transaction & append if valid (fields & money in blockchain)
					if c.validTrans(&t) {
						c.transPool.Add(&t)
					}
				}()
			}
		}
	}()
	ticker := time.NewTicker(3 * time.Second)
	quit := make(chan struct{})
	go func() {
		// send to network
		for {
			select {
			case <-ticker.C:
				// do stuff
				// check if leader
				// create block from transPool and send it to network
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	return nil
}
