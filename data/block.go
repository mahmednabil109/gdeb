package data

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

func GenesisBlock(stakeDist map[string]float64) *Block {
	stakeholders := make([]Stakeholder, len(stakeDist))
	for pub, stake := range stakeDist {
		sh := Stakeholder{PublicKey: pub, Stake: stake}
		stakeholders = append(stakeholders, sh)
	}

	return &Block{
		Slot:              1,
		StakeDistribution: stakeholders,
		Nonce:             "RANDOM_NONCE",
	}
}

func (b *Block) Add(t Transaction) {
	b.Transactions = append(b.Transactions, t)
}
