package OracleConnection

import (
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"sync"
)

var Host = "127.0.0.1:8383/"

type OraclePool struct {
	Connections     map[string]*SafeConnection
	mutex           sync.Mutex
	SubscribeChan   chan *SubscribeMsg
	UnsubscribeChan chan *UnsubscribeMsg
	Subscribers     []*SubscribeMsg
}

func NewOraclePool() *OraclePool {

	pool := &OraclePool{
		Connections:     map[string]*SafeConnection{},
		mutex:           sync.Mutex{},
		SubscribeChan:   make(chan *SubscribeMsg),
		UnsubscribeChan: make(chan *UnsubscribeMsg),
		Subscribers:     make([]*SubscribeMsg, 0),
	}
	go pool.subscribeListener()
	go pool.unsubscribeListener()

	return pool
}

func (pool *OraclePool) unsubscribeListener() {
	log.Println("Initiated unsub listener!!")
	defer log.Println("Exit Listener")
	for {
		select {
		case msg := <-pool.UnsubscribeChan:
			log.Println("Received unsubMsg from:", msg.VM)
			vmId := msg.VM
			i := 0
			newSubscribers := make([]*SubscribeMsg, len(pool.Subscribers))
			log.Println(pool.Connections)
			for _, v := range pool.Subscribers {
				log.Println("Oracle Key", v.Index)
				safeConn := pool.Connections[v.OracleKey]
				if v.VM == vmId {
					safeConn.decCount()
					if safeConn.Count == 0 {
						safeConn.close()
						pool.removeConnection(v.OracleKey)
					}
				} else {
					newSubscribers[i] = v
					i++
				}
			}
			pool.Subscribers = newSubscribers
			log.Println("Current Connections:", pool.Connections)
			for _, v := range pool.Connections {
				log.Println(v.Count)
			}
		}
	}
}

func (pool *OraclePool) addConnection(topic string, conn *SafeConnection) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()
	pool.Connections[topic] = conn
}

func (pool *OraclePool) removeConnection(topic string) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()
	delete(pool.Connections, topic)
}

func (pool *OraclePool) getConnection(oracleUrl, topic string) (*SafeConnection, error) {

	if conn, isOk := pool.Connections[topic]; isOk {
		conn.incCount()
		return conn, nil
	}
	u := url.URL{Scheme: "ws", Host: oracleUrl, Path: "sub"}

	log.Println("Started to connect to", u.String())
	connection, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	log.Println("Connected to", u.String())
	if err != nil {
		log.Println(err)
		return nil, errors.New("no such Oracle")
	}

	newConn := &SafeConnection{
		Url:   topic,
		conn:  connection,
		Count: 1,
		mutex: sync.Mutex{},
	}
	pool.addConnection(topic, newConn)

	go pool.messageReceiver(pool.Connections[topic])

	return pool.Connections[topic], nil
}

func (pool *OraclePool) messageReceiver(safeConn *SafeConnection) {
	conn := safeConn.conn
	for {
		var res ResponseMsg
		err := conn.ReadJSON(&res)
		if err != nil {
			log.Println("Close: " + err.Error())
			break
		} else {
			pool.broadcast(res)
		}
	}
}

func (pool *OraclePool) broadcast(msg ResponseMsg) {
	log.Println("Start Broadcast")
	for _, v := range pool.Subscribers {
		if isSubscriber := v.OracleKey == msg.Key; isSubscriber {
			v.BroadcastChan <- &BroadcastMsg{
				Key:       msg.Key,
				Value:     msg.Value,
				Timestamp: msg.Timestamp,
				Index:     v.Index,
			}
		}
	}
}
func (pool *OraclePool) addSubscriber(subMsg *SubscribeMsg) {
	pool.mutex.Lock()
	pool.Subscribers = append(pool.Subscribers, subMsg)
	pool.mutex.Unlock()
}

func (pool *OraclePool) subscribeListener() {
	for {
		select {
		case sub := <-pool.SubscribeChan:
			pool.addSubscriber(sub)
			safeConn, err := pool.getConnection(sub.Url, sub.OracleKey)
			if err != nil {
				sub.BroadcastChan <- &BroadcastMsg{Error: errors.New("no such oracle")}
			}

			err = safeConn.writeMsg(sub.OracleKey)

			if err != nil {
				log.Println(err)
			}
		}
	}
}
