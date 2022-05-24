package consensus

import (
	"github.com/mahmednabil109/gdeb/data"
)

// type Block struct {
// 	PreviousHash      string // aka state
// 	Slot              int    // index of slot block produced in
// 	VrfOutput         string
// 	VrfProof          string
// 	SlotLeader        string // public key of block creator
// 	Timestamp         string
// 	Transactions      []Transaction
// 	Nonce             string        //seeds slot leaders for selection process, no need to include in all blocks (just epoch transitions)
// 	StakeDistribution []Stakeholder //only in epoch transitions
// 	Signature         string        // signing all previous fields (proof that slot leader is who he claims to be)
// }

// type Transaction struct {
// 	// maybe introduce a transaction header type later
// 	Nonce           uint64 // in case same amount sent to same receiver later
// 	From            string
// 	To              string
// 	Amount          uint64
// 	Timestamp       string
// 	ContractCode    []byte
// 	GasPrice        uint64
// 	GasLimit        uint64
// 	ContractAddress string // provides consent for smart contracts involving other users (otherwise empty)
// 	Signature       string
// }

func ValidateBlock(b *data.Block, stakeDist map[string]float64) bool {

	// update stakeDist in some variable
	ValidateLeader(b.Nonce, b.SlotLeader, b.VrfProof, stakeDist)

	for _, trans := range b.Transactions {
		//note: transaction could have been invalidated after previous block added it
		if !trans.Validate() {
			return false
		}
	}
	return true
}
