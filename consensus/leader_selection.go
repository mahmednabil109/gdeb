package consensus

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/mahmednabil109/gdeb/blockchain"
	"github.com/yoseplee/vrf"
)

//optimal threshold calculation constant, acceptable placeholder until simulation is finished
const threshConst = 0.02

// quick validation for supposed mined block that node received from network
// determines whether network will broadcast this block further or not
func ValidateLeader(slot int, nonce string, public string, vrfProof string, stakeDist *blockchain.StakeDistribution) bool {
	proof, _ := hex.DecodeString(vrfProof)
	pub, _ := hex.DecodeString(public)
	vrf_input := []byte(fmt.Sprintf("%d%s", slot, nonce))

	valid, _ := vrf.Verify(pub, proof, vrf_input)
	if !valid {
		log.Printf("Proof of user %s is invalid!\n", hex.EncodeToString(pub))
		return false
	}

	stake := stakeDist.Get(hex.EncodeToString(pub))
	if stake <= 0 {
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

func checkIfLeader(vrf_input []byte, priv ed25519.PrivateKey, stakeDist *blockchain.StakeDistribution) (bool, []byte) {
	pub := priv.Public().(ed25519.PublicKey)
	proof, hash, err := vrf.Prove(pub, priv, vrf_input)

	if err != nil {
		return false, nil
	}

	stake := stakeDist.Get(hex.EncodeToString(pub))
	if stake <= 0 {
		return false, nil
	}

	threshold := calculateThresh(stake)
	output := hex.EncodeToString(hash)
	value := hexToBigInt(output)
	result := value.Cmp(threshold)
	if result == 1 {
		return false, nil
	}
	return true, proof
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
