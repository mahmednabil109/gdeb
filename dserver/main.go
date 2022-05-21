package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/mahmednabil109/gdeb/dserver/drpc"
	"google.golang.org/grpc"
)

var (
	UIPORT   = flag.Int("ui-port", 8282, "dserver ui port")
	PORT     = flag.Int("port", 16585, "dserver drpc port")
	upgrader = websocket.Upgrader{
		CheckOrigin:     func(r *http.Request) bool { return true },
		WriteBufferSize: 1024,
	}
)

type DServer struct {
	drpc.UnimplementedDServerServer

	Data chan *drpc.Pointers
}

func (s *DServer) Init() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *PORT))
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	drpc.RegisterDServerServer(server, s)
	log.Printf("grpc start listening %v", *PORT)
	if err := server.Serve(lis); err != nil {
		panic(err)
	}
}

func (s *DServer) UpdatePointers(ctx context.Context, pointers *drpc.Pointers) (*drpc.Empty, error) {
	go func() { s.Data <- pointers }()
	return &drpc.Empty{}, nil
}

func main() {
	flag.Parse()

	server := DServer{
		Data: make(chan *drpc.Pointers),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			panic(err)
		}

		go func(s DServer) {
			for d := range s.Data {
				err := conn.WriteJSON(struct {
					Id        string `json:"id"`
					Successor string `json:"successor"`
					D         string `json:"d"`
				}{
					Id:        d.Id,
					Successor: d.Successor,
					D:         d.D,
				})
				if err != nil {
					log.Printf("unable to send update %v", err)
					conn.Close()
					return
				}
			}
		}(server)
	})

	go server.Init()
	log.Printf("UI server start listening %v", *UIPORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", *UIPORT), nil))
}
