package data

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
)

type Transaction struct {
	// maybe introduce a transaction header type later
	Nonce           uint64 // in case same amount sent to same receiver later
	From            string
	To              string
	Amount          uint64
	Timestamp       string
	ContractCode    []byte
	GasPrice        uint64
	GasLimit        uint64
	ContractAddress string // provides consent for smart contracts involving other users (otherwise empty)
	Signature       string
}
type unsignedTransaction struct {
	Nonce           uint64
	From            string
	To              string
	Amount          uint64
	Timestamp       string
	ContractCode    []byte
	GasPrice        uint64
	GasLimit        uint64
	ContractAddress string
}

func (t *Transaction) unsigned() unsignedTransaction {
	return unsignedTransaction{
		Nonce:           t.Nonce,
		From:            t.From,
		To:              t.To,
		Amount:          t.Amount,
		Timestamp:       t.Timestamp,
		ContractCode:    t.ContractCode,
		GasPrice:        t.GasPrice,
		GasLimit:        t.GasLimit,
		ContractAddress: t.ContractAddress,
	}
}

func (t *Transaction) Validate() bool {
	unsigned := t.unsigned()
	unsignedJSON, _ := json.Marshal(unsigned)
	signature, _ := hex.DecodeString(t.Signature)
	pk, _ := hex.DecodeString(t.From)
	publicKey := ed25519.PublicKey(pk)
	return ed25519.Verify(publicKey, unsignedJSON, signature)
}

func (t *Transaction) Sign() {
	// private key of node hard-coded for now
	//mock values:
	// "PublicKey": "317f55e2a701a149d5385fe10ff8fe92762309914e85a91f59e90743f5909b52",
	// "PrivateKey": "b63d61f2c90efa388a565b39f84e903538b8347773262019b9c4f2189abca9ee317f55e2a701a149d5385fe10ff8fe92762309914e85a91f59e90743f5909b52",

	sk, _ := hex.DecodeString("b63d61f2c90efa388a565b39f84e903538b8347773262019b9c4f2189abca9ee317f55e2a701a149d5385fe10ff8fe92762309914e85a91f59e90743f5909b52")
	privateKey := ed25519.PrivateKey(sk)
	unsigned := t.unsigned()
	unsignedJSON, _ := json.Marshal(unsigned)
	signature := ed25519.Sign(privateKey, unsignedJSON)
	t.Signature = hex.EncodeToString(signature)
}
