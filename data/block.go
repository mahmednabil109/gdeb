package data

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
)

type Block struct {
	PreviousHash      string // aka state
	Slot              int    // index of slot block produced in
	VrfOutput         string
	VrfProof          string
	SlotLeader        string // public key of block creator
	Timestamp         string
	Transactions      []Transaction
	Nonce             string        //seeds slot leaders for selection process, no need to include in all blocks (just epoch transitions)
	StakeDistribution []Stakeholder //only in epoch transitions
	Signature         string        // signing all previous fields (proof that slot leader is who he claims to be)
}

type unsignedBlock struct {
	PreviousHash      string // aka state
	Slot              int    // index of slot block produced in
	VrfOutput         string
	VrfProof          string
	SlotLeader        string // public key of block creator
	Timestamp         string
	Transactions      []Transaction
	Nonce             string        //seeds slot leaders for selection process, no need to include in all blocks (just epoch transitions)
	StakeDistribution []Stakeholder //only in epoch transitions
}

func (b *Block) unsigned() unsignedBlock {
	return unsignedBlock{
		PreviousHash:      b.PreviousHash,
		Slot:              b.Slot,
		VrfOutput:         b.VrfOutput,
		VrfProof:          b.VrfProof,
		SlotLeader:        b.SlotLeader,
		Timestamp:         b.Timestamp,
		Transactions:      b.Transactions,
		Nonce:             b.Nonce,
		StakeDistribution: b.StakeDistribution,
	}
}

func GenesisBlock(stakeDist map[string]float64) *Block {
	stakeholders := make([]Stakeholder, len(stakeDist))
	for pub, stake := range stakeDist {
		sh := Stakeholder{PublicKey: pub, Stake: stake}
		stakeholders = append(stakeholders, sh)
	}

	return &Block{
		PreviousHash:      "RANDOM_PREVIOUS_HASH",
		Slot:              1,
		StakeDistribution: stakeholders,
		Nonce:             "RANDOM_NONCE",
	}
}

func (b *Block) Add(t Transaction) {
	b.Transactions = append(b.Transactions, t)
}

func (b *Block) Sign(privateKey ed25519.PrivateKey) {
	unsigned := b.unsigned()
	unsignedJSON, _ := json.Marshal(unsigned)
	signature := ed25519.Sign(privateKey, unsignedJSON)
	b.Signature = hex.EncodeToString(signature)
}
