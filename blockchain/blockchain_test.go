package blockchain

import (
	"os"
	"testing"
	"time"

	"github.com/mahmednabil109/gdeb/data"
)

func Test_serialize(t *testing.T) {
	var blockchain Blockchain

	blockchain.Init()
	blockchain.Add(data.Block{
		Slot: 1,
	})

	time.Sleep(time.Second)

	if os.IsNotExist(BLOCKCHAIN_PATH) {
		t.Fatal("Blockchain Is Not Serialized Correctly")
	}
}

func Test_deserialize(t *testing.T) {
	file, err := os.Create(BLOCKCHAIN_PATH)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err := writeGob(BLOCKCHAIN_PATH, []data.Block{
		data.Block{
			Slot: 2,
		},
	})
	if err != nil {
		panic(err)
	}
	
	var blockchain Blockchain
	blockchain.Init()
	
	time.Sleep(time.Second)

	iflen(blockchain.GetBlockchain()) != 1 {
		log.Fatal("hamda !!")
	}

}
