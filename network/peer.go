package network

import (
	"log"
	"net"

	pd "github.com/mahmednabil109/gdeb/network/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	retryPolicy = `{
		"methodConfig": [{
		  "name": [{"service": "rpc.Koorde"}],
		  "waitForReady": true,
		  "retryPolicy": {
			  "MaxAttempts": 4,
			  "InitialBackoff": ".01s",
			  "MaxBackoff": ".01s",
			  "BackoffMultiplier": 1.0,
			  "RetryableStatusCodes": [ "UNAVAILABLE" ]
		  }
		}]}`
)

type Peer struct {
	NetAddr  *net.TCPAddr
	NodeAddr ID
	Start    ID
	Interval []ID
	kc       pd.KoordeClient
	conn     *grpc.ClientConn
}

func (p *Peer) InitConnection() error {
	if p.kc != nil {
		return nil
	}

	conn, err := grpc.Dial(
		p.NetAddr.String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(retryPolicy),
	)

	if err != nil {
		log.Fatalf("can't dial: %v", err)
		return err
	}
	p.conn = conn

	p.kc = pd.NewKoordeClient(p.conn)
	// log.Printf("connection Done With %s", p.NetAddr.String())
	return nil
}

func (p *Peer) CloseConnection() {
	p.kc = nil
	p.conn.Close()
}
