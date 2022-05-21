package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mahmednabil109/gdeb/network"
	"github.com/mahmednabil109/gdeb/network/utils"
)

var (
	port      = flag.Int("port", 16585, "port for grpc of the localnode")
	first     = flag.Int("first", 1, "flag to mark the node as the first in the network")
	bootstrap = flag.String("bootstrap", "127.0.0.1:16585", "ip for bootstraping node")
)

func main() {
	flag.Parse()

	var n network.Node

	err := n.Init(*port)
	if err != nil {
		log.Printf("faild to init the localnode: %v", err)
	}
	f, err := os.Create(fmt.Sprintf("%s.log", n.NetAddr.String()))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	log.SetOutput(f)

	if *first == 0 {
		log.Println("start to join")
		n.Join(utils.ParseIP(*bootstrap), *port)
		log.Printf("nodeID %v", n)
		log.Printf("Successor %v", n.Successor)
		log.Printf("D %v", n.D)
		log.Printf("Predecessor %v", n.Predecessor)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	<-sigs

	// n.NodeShutdown <- true
	log.Printf("Programe ended")
}
