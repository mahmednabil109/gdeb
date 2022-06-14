package network

import (
	"github.com/mahmednabil109/gdeb/data"
	"github.com/mahmednabil109/gdeb/network/rpc"
	"github.com/mahmednabil109/gdeb/network/utils"
)

/* Packet parsers and forms */

// parse_lookup_reply parses the rpc.PeerPacket into a Peer struct
func parse_peer_packet(reply *rpc.PeerPacket) *Peer {
	if reply == nil || reply.SrcId == "" {
		return nil
	}
	return &Peer{
		NodeAddr: utils.ParseID(reply.SrcId),
		NetAddr:  utils.ParseIP(reply.SrcIp),
		Start:    utils.ParseID(reply.Start),
		Interval: []ID{utils.ParseID(reply.Interval[0]), utils.ParseID(reply.Interval[1])},
	}
}

func form_peer_packet(peer *Peer) *rpc.PeerPacket {
	if peer == nil {
		return &rpc.PeerPacket{}
	}
	return &rpc.PeerPacket{
		SrcId:    peer.NodeAddr.String(),
		SrcIp:    peer.NetAddr.String(),
		Start:    peer.Start.String(),
		Interval: []string{peer.Interval[0].String(), peer.Interval[1].String()},
	}
}

func parse_block_packet(b *rpc.Block) *data.Block {
	transactions := []data.Transaction{}
	stakes := []data.Stakeholder{}

	for _, s := range stakes {
		stakes = append(stakes, data.Stakeholder{
			PublicKey: s.PublicKey,
			Stake:     s.Stake,
		})
	}

	for _, t := range b.Transactions {
		transactions = append(transactions, *parse_transaction_packet(t))
	}

	return &data.Block{
		PreviousHash:      b.PreviousHash,
		Slot:              int(b.Slot),
		VrfOutput:         b.VrfOutput,
		VrfProof:          b.VrfProof,
		SlotLeader:        b.SlotLeader,
		Transactions:      transactions,
		StakeDistribution: stakes,
		Nonce:             b.Nonce,
		Signature:         b.Signature,
	}
}

func form_block_packet(b *data.Block) *rpc.Block {
	transactions := []*rpc.Transaction{}
	stakes := []*rpc.StackHolder{}

	for _, t := range b.Transactions {
		transactions = append(transactions, form_trans_packet(&t))
	}

	for _, s := range b.StakeDistribution {
		stakes = append(stakes, &rpc.StackHolder{
			PublickKey: s.PublicKey,
			Stack:      s.Stake,
		})
	}

	return &rpc.Block{
		PreviousHash:      b.PreviousHash,
		Slot:              int32(b.Slot),
		VrfOutput:         b.VrfOutput,
		VrfProof:          b.VrfProof,
		SlotLeader:        b.SlotLeader,
		Transactions:      transactions,
		StakeDistribution: stakes,
		Nonce:             b.Nonce,
		Signature:         b.Signature,
	}
}

func parse_transaction_packet(t *rpc.Transaction) *data.Transaction {

	return &data.Transaction{
		Nonce:           t.Nonce,
		From:            t.From,
		To:              t.To,
		Amount:          t.Amount,
		Timestamp:       t.Timestamp,
		ContractCode:    t.ContractCode,
		GasPrice:        t.GasPrice,
		GasLimit:        t.GasLimit,
		ConsumedGas:     t.ConsumedGas,
		ContractAddress: t.ContractAddress,
		Signature:       t.Signature,
	}
}
func form_trans_packet(t *data.Transaction) *rpc.Transaction {
	return &rpc.Transaction{
		Nonce:           t.Nonce,
		From:            t.From,
		To:              t.To,
		Amount:          t.Amount,
		Timestamp:       t.Timestamp,
		ContractCode:    t.ContractCode,
		GasPrice:        t.GasPrice,
		GasLimit:        t.GasLimit,
		ConsumedGas:     t.ConsumedGas,
		ContractAddress: t.ContractAddress,
		Signature:       t.Signature,
	}
}
