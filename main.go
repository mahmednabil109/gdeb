package main

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mahmednabil109/gdeb/config"
	"github.com/mahmednabil109/gdeb/data"
	"github.com/mahmednabil109/gdeb/node"
)

var (
	PORT      = flag.Int("port", 8080, "port")
	FIRST     = flag.Int("first", 1, "flag to mark the node as the first in the network")
	BOOTSTRAP = flag.String("bootstrap", "127.0.0.1:16585", "ip for bootstraping node")
	PK        = flag.Int("pk", 0, "pk")
	Send      = flag.Int("send", 0, "send transaction")
)

func main() {
	flag.Parse()

	config := config.New(*PK)
	node := node.New(config)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "port", *PORT)
	ctx = context.WithValue(ctx, "first", *FIRST)
	ctx = context.WithValue(ctx, "bootstrap", *BOOTSTRAP)

	f, err := os.Create(fmt.Sprintf("%d.log", *PORT))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	log.SetOutput(f)

	node.Init(ctx)

	rand.Seed(time.Now().UnixNano())
	if *Send == 1 {

		go func() {
			time.Sleep(2 * time.Second)
			log.Print("start sending")
			node.Consensus.SendTransaction(data.Transaction{
				From:   hex.EncodeToString(node.Consensus.PrivateKey.Public().(ed25519.PublicKey)),
				To:     "eb71de478e31020245677e9c4dab62200ce59dd8b45fd0462673822f73f807d0",
				Amount: uint64(rand.Intn(10)),
				Nonce:  uint64(rand.Intn(10)),
			})
		}()
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)

	<-sig
}
