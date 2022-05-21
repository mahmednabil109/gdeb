package client

import (
	"context"
	"log"

	"github.com/mahmednabil109/gdeb/dserver/drpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	RetryPolicy = `{
		"methodConfig": [{
		  "name": [{"service": "drpc.DServer"}],
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

type Client struct {
	dc   drpc.DServerClient
	conn *grpc.ClientConn
}

func (c *Client) Init(addr string) error {
	if c.conn != nil {
		return nil
	}

	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(RetryPolicy),
	)
	if err != nil {
		log.Print(err)
		return err
	}

	c.conn = conn
	c.dc = drpc.NewDServerClient(c.conn)
	log.Printf("conn done with %s", addr)
	return nil
}

func (c *Client) Update(id, successor, d string) {
	_, err := c.dc.UpdatePointers(context.Background(), &drpc.Pointers{Id: id, Successor: successor, D: d})
	if err != nil {
		log.Print(err)
	}
}
