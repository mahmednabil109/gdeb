package consensus

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/mahmednabil109/gdeb/blockchain"
	"github.com/mahmednabil109/gdeb/communication"
	"github.com/mahmednabil109/gdeb/data"
	"github.com/yoseplee/vrf"
)

const (
	transCount = 15
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
	TransPool           *blockchain.TransPool
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
		TransPool:           blockchain.NewTransPool(),
		PrivateKey:          privateKey,
		epochNonce:          "RANDOM_NONCE",
	}
}

// validate without the sign
func (c *Consensus) SendTransaction(t data.Transaction) {
	// if c.validTrans(&t) {
	log.Print("valid transaction")
	t.Sign(c.PrivateKey)
	c.ChanConsTransaction <- t
	// } else {
	// log.Print("invalid transaction")
	// }
}

func (c *Consensus) CreateBlock() {
	slot := c.Blockchain.GetSlot() + 1
	vrf_input := []byte(fmt.Sprintf("%d%s", slot, c.epochNonce))

	valid, proof := checkIfLeader(vrf_input, c.PrivateKey, c.stakeDist)
	if !valid {
		return
	}
	transactions := c.TransactionsForBlock()
	if transactions == nil {
		return
	}

	log.Print("start mining")

	block := data.Block{
		PreviousHash: c.Blockchain.GetPrevHash(),
		Slot:         slot,
		VrfOutput:    hex.EncodeToString(vrf.Hash(proof)),
		VrfProof:     hex.EncodeToString(proof),
		SlotLeader:   hex.EncodeToString(c.PrivateKey.Public().(ed25519.PublicKey)),
		Transactions: transactions, //get from transaction pool and update it
	}
	//adds signature field to block
	block.Sign(c.PrivateKey)
	c.ChanConsBlock <- block

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
	c.TransPool.Update(b.Transactions)
}

//note:
// 1) transaction could have been invalidated after previous block added it
// 2) several transaction could come from 1 uesr, enough money for each one seperately, but not for their some (can be invalid if validating 1 by 1)
//    --> combining transactions of same sender beforhand (looping and validating), using hashmap
func (c *Consensus) validBlock(b *data.Block) bool {
	leader := ValidateLeader(b.Slot, c.epochNonce, b.SlotLeader, b.VrfProof, c.stakeDist)
	if !leader {
		return false
	}
	for _, t := range b.Transactions {
		if t.ContractCode != nil {
			// TODO: verifiable VM
			// skip validation of contract transaction for no
			continue
		} else if !c.validTrans(&t) {
			return false
		} else {
			return false
		}

	}
	return true
}

//necessary as some transactions could be from same user and could potentially cause wrong valadation
func (c *Consensus) TransactionsForBlock() []data.Transaction {
	tp := c.TransPool
	tp.Mux.Lock()
	defer tp.Mux.Unlock()

	if len(tp.Pool) == 0 {
		return nil
	}

	counter := 0
	transactions := make([]data.Transaction, 0, transCount)
	moneySum := make(map[string]float64)

	for _, t := range tp.Pool {
		if counter > transCount {
			break
		}
		val, ok := moneySum[t.From]
		if !ok && c.validTrans(t) {
			transactions = append(transactions, *t)
			moneySum[t.From] = float64(t.Amount)
			delete(tp.Pool, t.Signature)
			counter += 1
		} else if ok && c.userMoney.HasEnoughMoney(t.From, val+float64(t.Amount)) {
			transactions = append(transactions, *t)
			moneySum[t.From] += val
			delete(tp.Pool, t.Signature)
			counter += 1
		}

	}
	if len(transactions) == 0 {
		return nil
	}

	return transactions
}

func (c *Consensus) Init() error {
	go func() {
		// recieve from network
		for {
			select {
			case b := <-c.ChanNetBlock:
				go func() {
					//validate block & append if valid (leader and fields correct)
					log.Printf("consensus: recieved block %+v", b)
					if c.validBlock(&b) {
						c.updateChain(&b)
					} else {
						log.Print("block is not valid")

					}
				}()
			case t := <-c.ChanNetTransaction:
				go func() {
					//validate transaction & append if valid (fields & money in blockchain)
					log.Printf("consensus: recieved transaction %+v", t)
					if c.validTrans(&t) {
						if t.ContractCode == nil {
							c.TransPool.Add(&t)
						} else {
							// Interpreter stuf

							// create channel and Interpreter instance
						}
						//
					} else {
						log.Print("transaction is not valid")
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
				// check if leader
				//    --> if yes, create block from transPool and send it to network
				go c.CreateBlock()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	return nil
}

func (c *Consensus) readTransaction(VMCons <-chan data.Transaction) {
	// sign
	// add to pool

	for t := range VMCons {
		c.ChanConsTransaction <- t
	}
}
