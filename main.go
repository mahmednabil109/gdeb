package main

import (
	"github.com/mahmednabil109/gdeb/OracleConnection"
	"log"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {

	pool := OracleConnection.NewOraclePool()

	topics := []string{"Temperature", "RicePrice", "Hamda"}

	broadcastChan := make(chan *OracleConnection.BroadcastMsg)
	idx := 0
	goFunc := func(id int) {
		topic := 0
		if id < 2 {
			topic = 1
		} else if id < 4 {
			topic = 2
		}
		for j := 0; j < 3; j++ {
			subMsg1 := OracleConnection.SubscribeMsg{
				VM:            id,
				OracleKey:     topics[topic],
				Url:           "127.0.0.1:8383",
				Index:         idx,
				BroadcastChan: broadcastChan,
			}
			idx++
			pool.SubscribeChan <- &subMsg1
		}
		wg.Done()
	}

	for i := 0; i < 6; i++ {
		wg.Add(1)
		go goFunc(i)
	}
	wg.Wait()
	time.Sleep(1000)
	for _, v := range pool.Subscribers {
		log.Println(v.Index)
	}

	isBroadcast := false
	for {
		select {
		case res := <-broadcastChan:
			log.Println(res.Value, res.Index)
			if res.Value == "0" {
				for _, v := range pool.Subscribers {
					log.Println(v)
				}

				for k, v := range pool.Connections {
					log.Println(k, ":", v)
				}
			} else if res.Value == "1" && !isBroadcast {
				isBroadcast = true
				msg := &OracleConnection.UnsubscribeMsg{
					VM: 0,
				}
				log.Println("Res =", 1)
				log.Println("Msg =", msg)
				pool.UnsubscribeChan <- msg
			}
		}
	}

}
