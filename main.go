package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"github.com/mahmednabil109/gdeb/config"
	"github.com/mahmednabil109/gdeb/consensus"
	"github.com/mahmednabil109/gdeb/data"
	"github.com/yoseplee/vrf"
	"log"
)

var stakeDist map[string]float64
var deployedContracts []string
var privateKey ed25519.PrivateKey

func setup() {
	config := config.New()
	pk := config.NodeKey()
	privateKey = pk

	data.LoadStakeDist("stakeDistribution.json", &stakeDist)
	//same behavior but for deployed contracts, useful for transactions the execute a contract (contract should exist to begin with)
	// data.LoadContracts("deployedContracts.json", &deployedContracts)
}

func main() {
	setup()

	//code snippet to test ValidateLeader function
	PublicKey, _ := hex.DecodeString("bd92fd2c61027f602170bf9f6608bc80cabc2f6e6834824fa67dc7fc745cbfe0")
	PrivateKey, _ := hex.DecodeString("e70b0983a423db62605c527109306d67e16a69d2f4d6641183242e1eac462d27bd92fd2c61027f602170bf9f6608bc80cabc2f6e6834824fa67dc7fc745cbfe0")
	nonce := 053464
	proofBytes, _, err := vrf.Prove(PublicKey, PrivateKey, []byte(fmt.Sprint(nonce)))
	if err != nil {
		log.Println(err)
	}

	val := consensus.ValidateLeader(nonce, PublicKey, proofBytes, stakeDist)
	log.Println(val)

}
