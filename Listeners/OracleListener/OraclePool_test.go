package OracleListener

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/mahmednabil109/gdeb/Messages"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
)

var wg sync.WaitGroup
var s = make([]string, 0)

func TestNewOraclePool(t *testing.T) {
	go initServer()
	wg.Wait()
	time.Sleep(3000 * time.Millisecond)
	pool := NewOraclePool()
	log.Println("Created New Pool")

	sendMsg := func(s, f string, t time.Duration) {
		wg.Add(1)
		log.Println("Entered ", f)
		for i := 0; i < 5; i++ {
			subMsg := &Messages.SubscribeMsg{
				OracleKey:     f,
				Url:           s,
				Index:         i,
				BroadcastChan: nil,
			}
			log.Println("Send Msg", i, "in", f)
			pool.SubscribeChan <- subMsg
			//	time.Sleep(t * time.Millisecond)

		}
		wg.Done()
	}

	for i := 0; i < 10; i++ {
		go sendMsg("hello", "go routine 1", time.Duration(5000))
		go sendMsg("bonjour", "go routine 2", time.Duration(4000))
	}

	//wg.Wait()
	time.Sleep(5000 * time.Millisecond)
	fmt.Println(pool.Connections)
	fmt.Println(pool.Subscribers)
	log.Println(len(s))
	log.Println(s)
}

var upgrader = websocket.Upgrader{} // use default options

func initServer() {
	wg.Add(1)
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
		}
		//
		defer c.Close()
		serverReceiveMsg(c, "hello")
	})
	http.HandleFunc("/bonjour", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
		}
		//
		defer c.Close()
		serverReceiveMsg(c, "bonjour")
	})

	log.Println("Starting server at port", Host)
	if err := http.ListenAndServe(Host, nil); err != nil {
		log.Fatal(err)
	}
	wg.Done()
}

func serverReceiveMsg(conn *websocket.Conn, server string) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
		}
		log.Println("Server:", server, "received Msg: %s", message)
		s = append(s, string(message))
	}

}
