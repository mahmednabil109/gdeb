package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/mahmednabil109/gdeb/blockchain"
	"github.com/mahmednabil109/gdeb/communication"
	"github.com/mahmednabil109/gdeb/config"
	"github.com/mahmednabil109/gdeb/consensus"
	"github.com/mahmednabil109/gdeb/data"
	"github.com/yoseplee/vrf"
)

var (
	totalCoins        uint64
	stakeDist         blockchain.StakeDistribution
	deployedContracts []string
	privateKey        ed25519.PrivateKey
	communNetwCons    communication.CommunNetwCons
)

func setup() {
	config := config.New()
	pk := config.NodeKey()
	privateKey = pk

	dist := make(map[string]float64)
	data.LoadStakeDist("stakeDistribution.json", &dist)
	stakeDist = blockchain.NewStakeDist(totalCoins, dist)

	communNetwCons = communication.CommunNetwCons{
		ChanNetBlock:        make(chan data.Block),
		ChanNetTransaction:  make(chan data.Transaction),
		ChanConsBlock:       make(chan data.Block),
		ChanConsTransaction: make(chan data.Transaction),
	}

	//same behavior but for deployed contracts, useful for transactions the execute a contract (contract should exist to begin with)
	// data.LoadContracts("deployedContracts.json", &deployedContracts)
}

func main() {
	setup()

	// suggestion to set up communication/ pass channels between different modules
	// cons := consensus.New(&communNetwCons, &stakeDist)
	// netw := network.New(&communNetwCons)

	//code snippet to test ValidateLeader function
	PublicKey, _ := hex.DecodeString("bd92fd2c61027f602170bf9f6608bc80cabc2f6e6834824fa67dc7fc745cbfe0")
	PrivateKey, _ := hex.DecodeString("e70b0983a423db62605c527109306d67e16a69d2f4d6641183242e1eac462d27bd92fd2c61027f602170bf9f6608bc80cabc2f6e6834824fa67dc7fc745cbfe0")
	pub := "bd92fd2c61027f602170bf9f6608bc80cabc2f6e6834824fa67dc7fc745cbfe0"
	nonce := "053464"
	proofBytes, _, err := vrf.Prove(PublicKey, PrivateKey, []byte(fmt.Sprint(nonce)))
	proof := hex.EncodeToString(proofBytes)
	if err != nil {
		log.Println(err)
	}
	val := consensus.ValidateLeader(nonce, pub, proof, &stakeDist)
	log.Println(val)

}
