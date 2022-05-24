package blockchain

import (
	"os"
	"sync"
	"time"

	"github.com/mahmednabil109/gdeb/data"
)

type Blockchain struct {
	// !handle fork ?? correct data structure?
	chain []data.Block
	mux   sync.Mutex
}

/*
				 / [slot 6] <= [slot 7] <= [slot 8] <= [slot 9]
[] <= [] <= [solt 5]
		           \ [slot 6] <= [slot 7] <= [slot 8]
*/

var (
	BLOCKCHAIN_PATH string = './.blockchain'
	SAVE_FREQ time.Duration = 800 * time.Millisecond
)

// write to disk (goroutine)
func (bc *Blockchain) Init() {
	bc.chain = []data.Block{}

	// read the blockchain  from disk if it exists
	if os.IsExist(BLOCKCHAIN_PATH) {
		go func(){
			bc.mux.Lock()
			defer bc.mux.Unlock()

			err := readGob(BLOCKCHAIN_PATH, bc.chain)
		}()
	}

	go func() {
		ticker := time.NewTicker(SAVE_FREQ)
		for {
			select {
			case <-ticker:
				// serialize the blockchain in the fs

				// to minimize the time we hold the mutex
				var chain_copy []data.Block
				bc.mux.Lock()
				copy(chain_copy, bc.chain)
				bc.mux.Unlock()

				err := writeGob(BLOCKCHAIN_PATH, chain_copy)
				if err != nil {
					panic(err)
				}

			}
		}
	}()
}

//! handle fork ??
func (bc *Blockchain) Add(block data.Block) {
	bc.mux.Lock()
	defer bc.mux.Unlock()

	bc.chain = append(bc.chain, block)
}

// returns a copy of the blockchain
func (bc *Blockchain) GetBlockchain() []data.Block {
	var chain_copy []data.Block
	copy(chain_copy, bc.chain)
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

func readGob(filgePath string, object interface{}) error {
	file, err := os.Open(filePath)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}

	file.Close()
	return err
}
