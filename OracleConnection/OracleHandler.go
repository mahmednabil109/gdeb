package OracleConnection

import (
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type SubscribeMessage struct {
	Url       string
	OracleKey string
	Index     int
	ResChan   chan Response
}

type OraclePool struct {
	connections      map[string]*SafeConnection
	mutex            sync.Mutex
	SubscribeChannel chan *SubscribeMessage
	subscribers      []*SubscribeMessage
}

type Response struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
	Error     error
}

func NewOraclePool() *OraclePool {

	pool := &OraclePool{
		connections:      map[string]*SafeConnection{},
		mutex:            sync.Mutex{},
		SubscribeChannel: make(chan *SubscribeMessage),
		subscribers:      make([]*SubscribeMessage, 0),
	}
	go pool.listen()

	return pool
}

func (pool *OraclePool) getConnection(url string) (*SafeConnection, error) {
	if conn, isOk := pool.connections[url]; isOk {
		conn.lock.Lock()
		conn.connectionCount++
		conn.lock.Unlock()
		return conn, nil
	}

	connection, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, errors.New("no such Oracle")
	}
	pool.mutex.Lock()
	pool.connections[url] = &SafeConnection{
		url:             url,
		conn:            connection,
		connectionCount: 1,
		lock:            sync.Mutex{},
	}
	pool.mutex.Unlock()

	go pool.messageReceiver(pool.connections[url])

	return pool.connections[url], nil
}

func (pool *OraclePool) messageReceiver(safeConn *SafeConnection) {
	conn := safeConn.conn
	for {
		var res struct {
			Key       string `json:"key"`
			Value     string `json:"value"`
			Timestamp string `json:"timestamp"`
		}
		err := conn.ReadJSON(&res)
		if err != nil {
			log.Println("Close: " + err.Error())
		}

		for _, v := range pool.subscribers {
			if isSubscriber := v.OracleKey == res.Key; isSubscriber {
				v.ResChan <- Response{
					Key:       res.Key,
					Value:     res.Value,
					Timestamp: res.Timestamp,
				}
			}
		}

	}
}

func (pool *OraclePool) listen() {
	for {
		select {
		case sub := <-pool.SubscribeChannel:
			pool.mutex.Lock()
			pool.subscribers = append(pool.subscribers, sub)
			pool.mutex.Unlock()
			safeConn, err := pool.getConnection(sub.Url)
			if err != nil {
				sub.ResChan <- Response{Error: errors.New("no such oracle")}
			}
			err = safeConn.writeMsg(sub.OracleKey)
			if err != nil {
				log.Println(err)
			}
		}

	}
}
