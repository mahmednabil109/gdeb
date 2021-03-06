package node

import (
	"context"
	"log"

	"github.com/mahmednabil109/gdeb/blockchain"
	"github.com/mahmednabil109/gdeb/communication"
	"github.com/mahmednabil109/gdeb/config"
	"github.com/mahmednabil109/gdeb/consensus"
	"github.com/mahmednabil109/gdeb/data"
	"github.com/mahmednabil109/gdeb/network"
	"github.com/mahmednabil109/gdeb/network/utils"
)

const (
	totalCoins = 1000
)

// full node
type Node struct {
	// VM isntance
	Network   *network.Node
	Consensus *consensus.Consensus
}

func New(config *config.Config) *Node {
	NCComm := communication.CommunNetwCons{
		ChanNetBlock:        make(chan data.Block),
		ChanNetTransaction:  make(chan data.Transaction),
		ChanConsBlock:       make(chan data.Block),
		ChanConsTransaction: make(chan data.Transaction),
	}

	privateKey := config.NodeKey()
	log.Print(privateKey)

	dist := make(map[string]float64)
	data.LoadStakeDist("stakeDistribution.json", &dist)
	stakeDist := blockchain.NewStakeDist(totalCoins, dist)

	consensus := consensus.New(&NCComm, &stakeDist, privateKey)
	network := network.New(&NCComm)
	node := Node{
		Network:   network,
		Consensus: consensus,
	}
	return &node
}

func (n *Node) Init(ctx context.Context) {
	port := ctx.Value("port").(int)

	// init the network
	n.Network.Init(port)
	if first := ctx.Value("first").(int); first == 0 {
		bootstrapIP := ctx.Value("bootstrap").(string)
		n.Network.Join(utils.ParseIP(bootstrapIP), port)
	}

	// init the consensus
	n.Consensus.Init()
}
