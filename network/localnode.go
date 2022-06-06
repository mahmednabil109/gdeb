package network

import (
	"context"
	"fmt"
	"log"
	"net"
	"sort"
	"time"

	"github.com/mahmednabil109/gdeb/communication"
	"github.com/mahmednabil109/gdeb/data"
	"github.com/mahmednabil109/gdeb/network/rpc"
	"github.com/mahmednabil109/gdeb/network/utils"
	"google.golang.org/grpc"
)

const (
	MAX_REQ_TIME = 20 * time.Second
)

type Node struct {
	// RPC stuff
	rpc.UnimplementedKoordeServer
	Peer
	D            *Peer
	DParents     []*Peer
	Successor    *Peer
	Predecessor  *Peer
	NodeShutdown chan bool

	// communication
	ChanNetBlock        chan<- data.Block
	ChanNetTransaction  chan<- data.Transaction
	ChanConsBlock       <-chan data.Block
	ChanConsTransaction <-chan data.Transaction

	// private
	s *grpc.Server
}

/* Node API */

// New() constructs a new network node,
// and it setups the communication channels with the Consensus Module
func New(c *communication.CommunNetwCons) *Node {
	node := Node{
		ChanNetBlock:        c.ChanNetBlock,
		ChanNetTransaction:  c.ChanNetTransaction,
		ChanConsBlock:       c.ChanConsBlock,
		ChanConsTransaction: c.ChanConsTransaction,
	}
	return &node
}

// Init initializes the first node in the network
// It inits the Successor, D pointers with default values (node itslef)
func (ln *Node) Init(port int) error {
	myIP := utils.GetMyIP()
	ln.NetAddr = &net.TCPAddr{IP: myIP, Port: port}
	ln.NodeAddr = utils.SHA1OF(ln.NetAddr.String())
	ln.Start = ln.NodeAddr
	ln.Interval = []ID{ln.NodeAddr, ln.NodeAddr}
	ln.Successor = &ln.Peer
	ln.Predecessor = &ln.Peer
	ln.DParents = []*Peer{&ln.Peer}
	ln.D = &ln.Peer
	err := init_grpc_server(ln, port)

	// listen to the consensus
	go func() {
		for {
			select {
			case b := <-ln.ChanConsBlock:
				go func() {
					log.Print("recieved block")
					ln.BroadCast(b)
				}()
			case t := <-ln.ChanConsTransaction:
				log.Print("recieved transaction")
				go func() {
					ln.BroadCast(t)
				}()
			}
		}
	}()

	// stablize
	go func() {
		ticker := time.NewTicker(time.Second)

		for {
			select {
			case <-ticker.C:
				ln.stablize()
				// log.Printf("%v, %v", ln.Successor, ln.D)
			case <-ln.NodeShutdown:
				ticker.Stop()
				return
			}
		}
	}()

	// fix dpointers
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)

		for {
			select {
			case <-ticker.C:
				ln.fixPointers()
			case <-ln.NodeShutdown:
				ticker.Stop()
				return
			}
		}
	}()

	return err
}

// Join initializes the node by executing Chord Join Algorithm
// It inits the Successor, D pointers
func (ln *Node) Join(nodeAddr *net.TCPAddr, port int) error {
	// log.Printf("Join %s", nodeAddr.String())

	if ln.s == nil {
		err := ln.Init(port)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	if nodeAddr == nil {
		return nil
	}

	peer := Peer{NetAddr: nodeAddr}
	peer.InitConnection()

	ctx, cancel := context.WithTimeout(context.Background(), MAX_REQ_TIME)
	defer cancel()
	reply, err := peer.kc.DLKup(ctx, &rpc.PeerPacket{SrcId: ln.NodeAddr.String()})
	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Printf("bootstrap result: %v", reply)

	ln.Predecessor = nil
	ln.D = nil
	ln.Successor = parse_peer_packet(reply)

	err = ln.Successor.InitConnection()
	if err != nil {
		log.Fatal("err", err)
		return err
	}

	return nil
}

func (ln *Node) Lookup(k ID) (*Peer, error) {
	kShift, i := select_imaginary_node(k, ln.NodeAddr, ln.Successor.NodeAddr)
	// log.Printf("init %s %s %s", k.String(), kShift.String(), i.String())

	lookupPacket := &rpc.LookupPacket{
		SrcId:  ln.NodeAddr.String(),
		SrcIp:  ln.NetAddr.String(),
		K:      k.String(),
		KShift: kShift.String(),
		I:      i.String()}
	ctx, cancel := context.WithTimeout(context.Background(), MAX_REQ_TIME)
	defer cancel()

	reply, err := ln.LookupRPC(ctx, lookupPacket)

	return parse_peer_packet(reply), err
}

func (ln *Node) BroadCast(thing interface{}) error {
	// constructe broadcast packtes
	BPacket := rpc.BroadCastPacket{}

	switch b := thing.(type) {
	case data.Block:
		BPacket.Type = rpc.PacketType_BlockT
		BPacket.Block = form_block_packet(&b)
	case data.Transaction:
		log.Print("broadcasting transaction!!")
		BPacket.Type = rpc.PacketType_TransT
		BPacket.Trans = form_trans_packet(&b)
	}

	log.Printf("init broadcasting of %+v", BPacket)

	// start the braodcast
	if !ln.Successor.NodeAddr.Equal(ln.D.NodeAddr) {
		err := ln.Successor.InitConnection()
		if err != nil {
			log.Fatal(err)
			return err
		}

		// log.Printf("Broadcast %s --> %s limit %s", info, ln.Successor.NetAddr.String(), ln.D.NetAddr.String())
		BPacket.Limit = ln.D.NodeAddr.String()
		ctx, cancel := context.WithTimeout(context.Background(), MAX_REQ_TIME)
		defer cancel()
		_, err = ln.Successor.kc.BroadCastRPC(ctx, &BPacket)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	err := ln.D.InitConnection()
	if err != nil {
		log.Fatal(err)
		return err
	}

	// log.Printf("Broadcast %s --> %s limit %s", info, ln.D.NetAddr.String(), ln.NetAddr.String())
	ctx, cancel := context.WithTimeout(context.Background(), MAX_REQ_TIME)
	BPacket.Limit = ln.NodeAddr.String()
	defer cancel()
	_, err = ln.D.kc.BroadCastRPC(ctx, &BPacket)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil

}

func (ln *Node) stablize() {
	if ln.Successor == nil {
		return
	}

	ln.Successor.InitConnection()

	ctx, cancel := context.WithTimeout(context.Background(), MAX_REQ_TIME)
	defer cancel()
	pred, err := ln.Successor.kc.GetPredecessorRPC(ctx, &rpc.Empty{})
	if err != nil {
		log.Fatal(err)
		return
	}

	predecessor_peer := parse_peer_packet(pred)

	if predecessor_peer != nil && predecessor_peer.NodeAddr.InLRXRange(ln.NodeAddr, ln.Successor.NodeAddr) {
		// log.Printf("better successor %v x %v", predecessor_peer.NetAddr, ln.Successor.NetAddr)

		if !ln.Successor.NodeAddr.Equal(predecessor_peer.NodeAddr) {

			ln.Successor = predecessor_peer
			err := ln.Successor.InitConnection()
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// notify only when there is a change
	ctx, cancel = context.WithTimeout(context.Background(), MAX_REQ_TIME)
	defer cancel()
	currentNode_peer := form_peer_packet(&ln.Peer)
	// log.Printf("%+v", ln.Successor)
	ln.Successor.InitConnection()
	_, err = ln.Successor.kc.NotifyRPC(ctx, currentNode_peer)
	if err != nil {
		log.Fatal(err)
	}
}

func (ln *Node) fixPointers() {
	D_id, _ := ln.NodeAddr.LeftShift()

	peer, err := ln.Lookup(D_id)
	if err != nil {
		log.Fatal(err)
		return
	}

	// if peer.NodeAddr.Equal(ln.NodeAddr) {
	// 	return
	// }

	err = peer.InitConnection()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer peer.CloseConnection()

	ctx, cancle := context.WithTimeout(context.Background(), MAX_REQ_TIME)
	defer cancle()
	peer_pred_packet, err := peer.kc.GetPredecessorRPC(ctx, &rpc.Empty{})
	if err != nil {
		log.Fatal(err)
		return
	}

	peer_pred := parse_peer_packet(peer_pred_packet)
	// log.Printf("%+v", peer_pred)
	if peer_pred == nil || peer_pred.NodeAddr == nil || (ln.D != nil && peer_pred.NodeAddr.Equal(ln.D.NodeAddr)) {
		return
	}
	// prevD := ln.D
	ln.D = peer_pred

	err = ln.D.InitConnection()
	if err != nil {
		log.Fatal(err)
	}
}

/* Helper Methods */

// Select the best imaginary node to start the lookup from
// that is in the range (m, m.Successor] in the ring
func select_imaginary_node(k, m, successor ID) (ID, ID) {

	for i := 2*len(m) - 1; i >= 0; i-- {
		_id := m.MaskLowerWith(k, i).AddOne(i)

		if ID(_id).InLXRange(m, successor) {
			for j := 0; j < i; j++ {
				k, _ = k.LeftShift()
			}
			return k, _id
		}
	}
	// no Match
	return k, m.AddOne(0)
}

// init_grpc_server creates a tcp socket and registers
// a new grpc server for Node.s
func init_grpc_server(ln *Node, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("faild to listen to %v", err)
		return err
	}

	ln.s = grpc.NewServer()
	rpc.RegisterKoordeServer(ln.s, ln)
	go func() {
		log.Printf("grpc start listening %v", ln.NetAddr)
		if err := ln.s.Serve(lis); err != nil {
			log.Fatalf("faild to serve %v", err)
		}
	}()
	return nil
}

/* RPC impelementation */

func (ln *Node) BootStarpRPC(bctx context.Context, bootstrapPacket *rpc.BootStrapPacket) (*rpc.BootStrapReply, error) {
	src_id := ID(utils.ParseID(bootstrapPacket.SrcId))

	successor, err := ln.Lookup(src_id)
	if err != nil {
		return nil, err
	}

	d_id, _ := src_id.LeftShift()
	// lookup returns the successor
	d, err := ln.Lookup(d_id)
	if err != nil {
		return nil, err
	}
	// getting the predecessor
	d.InitConnection()
	ctx, cancel := context.WithTimeout(context.Background(), MAX_REQ_TIME)
	defer cancel()

	d_pre, err := d.kc.GetPredecessorRPC(ctx, &rpc.Empty{})
	if err != nil {
		return nil, err
	}

	return &rpc.BootStrapReply{
		Successor: form_peer_packet(successor),
		D:         d_pre,
	}, nil
}

func (ln *Node) LookupRPC(bctx context.Context, lookupPacket *rpc.LookupPacket) (*rpc.PeerPacket, error) {

	k := ID(utils.ParseID(lookupPacket.K))
	kShift := ID(utils.ParseID(lookupPacket.KShift))
	i := ID(utils.ParseID(lookupPacket.I))

	if k.Equal(ln.NodeAddr) {
		log.Printf("Me || %s", ln.NetAddr)
		return form_peer_packet(&ln.Peer), nil
	}

	if k.InLXRange(ln.NodeAddr, ln.Successor.NodeAddr) {
		return form_peer_packet(ln.Successor), nil
	}

	// log.Printf("second %s in (%s %s] %v !!", i, ln.NodeAddr, ln.Successor.NodeAddr, i.InLXRange(ln.NodeAddr, ln.Successor.NodeAddr))

	if ln.D != nil && i.InLXRange(ln.NodeAddr, ln.Successor.NodeAddr) {

		// TODO handle failer and pointer replacemnet
		ln.D.InitConnection()
		ctx, cancel := context.WithTimeout(context.Background(), MAX_REQ_TIME)
		defer cancel()

		KShift, _ := kShift.LeftShift()
		lookupPacket := &rpc.LookupPacket{
			SrcId:  ln.NodeAddr.String(),
			SrcIp:  ln.NetAddr.String(),
			K:      k.String(),
			KShift: KShift.String(),
			I:      i.TopShift(kShift).String()}
		reply, err := ln.D.kc.LookupRPC(ctx, lookupPacket)

		if err != nil {
			log.Printf("lookup faild: %v", err)
			return nil, err
		}
		return reply, nil
	}

	// TODO handle failer and pointer replacemnet
	ln.Successor.InitConnection()
	ctx, cancel := context.WithTimeout(context.Background(), MAX_REQ_TIME)
	defer cancel()

	reply, err := ln.Successor.kc.LookupRPC(ctx, lookupPacket)

	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (ln *Node) UrpcatePredecessorRPC(bctx context.Context, p *rpc.PeerPacket) (*rpc.PeerPacket, error) {
	old_predecessor := form_peer_packet(ln.Predecessor)
	ln.Predecessor = parse_peer_packet(p)

	return old_predecessor, nil
}

func (ln *Node) UrpcateSuccessorRPC(bctx context.Context, p *rpc.PeerListPacket) (*rpc.PeerListPacket, error) {
	ln.Successor = parse_peer_packet(p.Peers[0])
	Succ_id := ID(utils.ParseID(p.Peers[1].SrcId))

	// slice of the pointers to hand out
	pointers := make([]*rpc.PeerPacket, 0)
	for _, n := range ln.DParents {
		D_id, _ := n.NodeAddr.LeftShift()
		if D_id.InLXRange(ln.Successor.NodeAddr, Succ_id) {
			pointers = append(pointers, form_peer_packet(n))
		}
	}
	// log.Printf("hand over %v %d", pointers, len(ln.DParents))
	return &rpc.PeerListPacket{Peers: pointers}, nil
}

func (ln *Node) UrpcateDPointerRPC(bctx context.Context, p *rpc.PeerPacket) (*rpc.Empty, error) {
	ln.D = parse_peer_packet(p)

	return &rpc.Empty{}, nil
}

func (ln *Node) AddDParentRPC(bctx context.Context, p *rpc.PeerPacket) (*rpc.Empty, error) {
	peer := parse_peer_packet(p)
	ln.DParents = append(ln.DParents, peer)

	// sort the Dparents Pointers by the id
	sort.Slice(ln.DParents, func(i, j int) bool {
		return !ln.DParents[i].NodeAddr.InLXRange(ln.DParents[j].NodeAddr, MAX_ID)
	})
	// log.Print(ln.DParents)
	return &rpc.Empty{}, nil
}

func (ln *Node) RemoveDParentRPC(btcx context.Context, p *rpc.PeerPacket) (*rpc.Empty, error) {
	peer := parse_peer_packet(p)
	peer_idx := -1

	for i := range ln.DParents {
		if ln.DParents[i].NodeAddr.Equal(peer.NodeAddr) {
			peer_idx = i
			break
		}
	}

	if peer_idx != -1 {
		// remove that peer form the list
		ln.DParents = append(ln.DParents[:peer_idx], ln.DParents[peer_idx+1:]...)
	}
	return &rpc.Empty{}, nil
}

func (ln *Node) NotifyRPC(ctx context.Context, p *rpc.PeerPacket) (*rpc.Empty, error) {
	predecessor_peer := parse_peer_packet(p)

	if ln.Predecessor == nil || predecessor_peer.NodeAddr.InLRXRange(ln.Predecessor.NodeAddr, ln.NodeAddr) {
		// for the first node in the network
		if ln.Successor.NodeAddr.Equal(ln.NodeAddr) {
			ln.Successor = predecessor_peer
		}
		ln.Predecessor = predecessor_peer
	}
	return &rpc.Empty{}, nil
}

func (ln *Node) BroadCastRPC(ctx context.Context, b *rpc.BroadCastPacket) (*rpc.Empty, error) {
	// notify the consensus code
	log.Printf("recieved somthing %+v", b.Type)
	go func() {
		switch b.Type {
		case rpc.PacketType_BlockT:
			log.Print("recieved block")
			ln.ChanNetBlock <- *parse_block_packet(b.Block)
		case rpc.PacketType_TransT:
			log.Print("recieved transaction")
			ln.ChanNetTransaction <- *parse_transaction_packet(b.Trans)
		}
	}()

	LimitID := ID(utils.ParseID(b.Limit))

	if ln.Successor.NodeAddr.Equal(ln.D.NodeAddr) && ln.Successor.NodeAddr.InLRXRange(ln.NodeAddr, LimitID) {
		ctx, cancel := context.WithTimeout(context.Background(), MAX_REQ_TIME)
		defer cancel()

		err := ln.Successor.InitConnection()
		if err != nil {
			return nil, err
		}

		// log.Printf("Broadcast %s --> %s limit %s", b.Info, ln.Successor.NetAddr.String(), b.Limit)
		_, err = ln.Successor.kc.BroadCastRPC(ctx, b)
		return &rpc.Empty{}, err

	} else if ln.Successor.NodeAddr.InLRXRange(ln.NodeAddr, LimitID) {
		NewLimit := b.Limit
		flag := false

		if ln.D.NodeAddr.InLRXRange(ln.NodeAddr, LimitID) {
			NewLimit = ln.D.NodeAddr.String()
			flag = true
		}

		err := ln.Successor.InitConnection()
		if err != nil {
			return nil, err
		}

		// log.Printf("Broadcast %s --> %s limit %s", b.Info, ln.Successor.NetAddr.String(), NewLimit)
		b.Limit = NewLimit
		ctx, cancel := context.WithTimeout(context.Background(), MAX_REQ_TIME)
		defer cancel()
		_, err = ln.Successor.kc.BroadCastRPC(ctx, b)
		if err != nil {
			return nil, err
		}

		if flag {
			err := ln.D.InitConnection()
			if err != nil {
				return nil, err
			}

			b.Limit = LimitID.String()
			// log.Printf("Broadcast %s --> %s limit %s", b.Info, ln.D.NetAddr.String(), b.Limit)
			ctx, cancel := context.WithTimeout(context.Background(), MAX_REQ_TIME)
			defer cancel()
			_, err = ln.D.kc.BroadCastRPC(ctx, b)
			if err != nil {
				return nil, err
			}
		}
	}
	return &rpc.Empty{}, nil
}

func (ln *Node) GetSuccessorRPC(ctx context.Context, e *rpc.Empty) (*rpc.PeerPacket, error) {
	// TODO make sure that the pointer is valid
	return form_peer_packet(ln.Successor), nil
}

func (ln *Node) GetPredecessorRPC(ctx context.Context, e *rpc.Empty) (*rpc.PeerPacket, error) {
	// TODO make sure that the pointer is valid
	return form_peer_packet(ln.Predecessor), nil
}

/* DEBUG RPC */

func (n *Node) InitBroadCastRPC(ctx context.Context, b *rpc.BroadCastPacket) (*rpc.Empty, error) {
	// log.Printf("Init Broadcast %s", b.Info)

	// n.ConsensusAPI.AddBlock(mock.Block{Info: b.Info})
	err := n.BroadCast("")

	return &rpc.Empty{}, err
}

func (n *Node) DJoin(ctx context.Context, p *rpc.PeerPacket) (*rpc.Empty, error) {
	// Debug ports
	n.Join(utils.ParseIP(p.SrcIp), 8081)

	return &rpc.Empty{}, nil
}

func (n *Node) DSetSuccessor(ctx context.Context, p *rpc.PeerPacket) (*rpc.Empty, error) {
	n.Successor = parse_peer_packet(p)
	n.Successor.InitConnection()
	return &rpc.Empty{}, nil
}

func (n *Node) DSetD(ctx context.Context, p *rpc.PeerPacket) (*rpc.Empty, error) {
	n.D = parse_peer_packet(p)
	n.D.InitConnection()
	return &rpc.Empty{}, nil
}

func (n *Node) DGetID(ctx context.Context, e *rpc.Empty) (*rpc.PeerPacket, error) {
	return &rpc.PeerPacket{SrcId: n.NodeAddr.String()}, nil
}

func (n *Node) DGetPointers(ctx context.Context, e *rpc.Empty) (*rpc.Pointers, error) {
	return &rpc.Pointers{Succ: n.Successor.NodeAddr.String(), D: n.D.NodeAddr.String()}, nil
}

func (n *Node) DLKup(ctx context.Context, p *rpc.PeerPacket) (*rpc.PeerPacket, error) {
	reply, err := n.Lookup(utils.ParseID(p.SrcId))
	return form_peer_packet(reply), err
}

func (n *Node) DGetBlocks(ctx context.Context, e *rpc.Empty) (*rpc.BlocksPacket, error) {
	reply := rpc.BlocksPacket{}

	// for _, b := range n.ConsensusAPI.GetBlocks() {
	// 	reply.Block = append(reply.Block, &rpc.BlockPacket{Info: b.Info})
	// }

	return &reply, nil
}
