package consensus

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"github.com/yoseplee/vrf"
	"log"
	"math/big"
)

//optimal threshold calculation constant, acceptable placeholder until simulation is finished
const threshConst = 0.02

// quick validation for supposed mined block that node received from network
// determines whether network will broadcast this block further or not
func ValidateLeader(nonce int, pub ed25519.PublicKey, proof []byte, stakeDist map[string]float64) bool {
	valid, _ := vrf.Verify(pub, proof, []byte(fmt.Sprint(nonce)))
	if !valid {
		log.Printf("Proof of user %s is invalid!\n", hex.EncodeToString(pub))
		return false
	}

	stake := stakeDist[hex.EncodeToString(pub)]
	if stake == 0 {
		return false
	}

	threshold := calculateThresh(stake)
	output := hex.EncodeToString(vrf.Hash(proof))
	value := hexToBigInt(output)
	result := value.Cmp(threshold)
	if result == 1 {
		return false
	}
	return true
}

func calculateThresh(stake float64) *big.Int {
	return thresholdBigInt((1 + threshConst) * stake)
}

func thresholdBigInt(stake float64) *big.Int {
	bigNum := new(big.Int)
	hex := "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
	bigNum.SetString(hex, 16)
	bigFloat := big.NewFloat(0).SetInt(bigNum)
	s := bigFloat.Mul(bigFloat, big.NewFloat(stake))
	backToInt := new(big.Int)
	s.Int(backToInt)
	return backToInt
}
func hexToBigInt(hex string) *big.Int {
	bigNum := new(big.Int)
	bigNum.SetString(hex, 16)
	bigFloat := big.NewFloat(0).SetInt(bigNum)
	backToInt := new(big.Int)
	bigFloat.Int(backToInt)
	return backToInt
}
