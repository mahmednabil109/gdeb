package blockchain

import (
	"encoding/gob"
	"os"
	"sync"
	"time"

	"github.com/mahmednabil109/gdeb/data"
)

// handle fork (appropriate data structure)
type Blockchain struct {
	Chain []*data.Block
	Mux   sync.Mutex
}

/*
				 / [slot 6] <= [slot 7] <= [slot 8] <= [slot 9]
[] <= [] <= [solt 5]
		           \ [slot 6] <= [slot 7] <= [slot 8]
*/

func NewBlockchain(dist map[string]float64) *Blockchain {
	chain := make([]*data.Block, 0)
	chain = append(chain, data.GenesisBlock(dist))
	return &Blockchain{
		Chain: chain,
	}
}

var (
	BLOCKCHAIN_PATH string        = "./.blockchain"
	SAVE_FREQ       time.Duration = 800 * time.Millisecond
)

// write to disk (goroutine)
// func (bc *Blockchain) Init() {
// 	bc.chain = []data.Block{}

// 	if _, err := os.Stat(BLOCKCHAIN_PATH); err != nil {
// 		if os.IsNotExist(err) {
// 			// file does not exist
// 		} else {
// 			go func() {
// 				bc.mux.Lock()
// 				defer bc.mux.Unlock()
// 				// err := readGob(BLOCKCHAIN_PATH, bc.chain)
// 			}()
// 		}
// 	}

// 	go func() {
// 		ticker := time.NewTicker(SAVE_FREQ)
// 		for {
// 			select {
// 			case <-ticker.C:
// 				// serialize the blockchain in the fs

// 				// to minimize the time we hold the mutex
// 				var chain_copy []data.Block
// 				bc.mux.Lock()
// 				copy(chain_copy, bc.chain)
// 				bc.mux.Unlock()

// 				err := writeGob(BLOCKCHAIN_PATH, chain_copy)
// 				if err != nil {
// 					panic(err)
// 				}

// 			}
// 		}
// 	}()
// }

func (bc *Blockchain) GetSlot() int {
	return bc.Chain[len(bc.Chain)-1].Slot
}
func (bc *Blockchain) GetPrevHash() string {
	return bc.Chain[len(bc.Chain)-1].PreviousHash
}

// block paramater has already been validated, add to blockchain and handle forks (keeping track of side chains and longest one in them)
func (bc *Blockchain) Update(block *data.Block) {
	bc.Mux.Lock()
	defer bc.Mux.Unlock()

	bc.Chain = append(bc.Chain, block)
}

// returns a copy of the blockchain
func (bc *Blockchain) GetBlockchain() []*data.Block {
	var chain_copy []*data.Block
	copy(chain_copy, bc.Chain)
	return chain_copy
}

// serialization utils

func writeGob(filePath string, object interface{}) error {
	file, err := os.Create(filePath)
	if err != nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}

	file.Close()
	return err
}

func readGob(filePath string, object interface{}) error {
	file, err := os.Open(filePath)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}

	file.Close()
	return err
}
