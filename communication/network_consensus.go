package communication

import (
	"github.com/mahmednabil109/gdeb/data"
)

type CommunNetwCons struct {
	ChanNetBlock        chan data.Block
	ChanNetTransaction  chan data.Transaction
	ChanConsBlock       chan data.Block
	ChanConsTransaction chan data.Transaction
}

// func (c *CommunNetwCons) NetworkChannels() *CommunNetwCons {
// 	return &CommunNetwCons{
// 		ChanNetBlock:        c.ChanNetBlock,
// 		ChanNetTransaction:  c.ChanNetTransaction,
// 		ChanConsBlock:       c.ChanConsBlock,
// 		ChanConsTransaction: c.ChanConsTransaction,
// 	}
// }

// func (c *CommunNetwCons) ConsensusChannels() *CommunNetwCons {
// 	return nil
// }
