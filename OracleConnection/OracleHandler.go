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
	ResChan   chan Response
}

type OraclePool struct {
	connections      map[string]*SafeConnection
	lock             sync.Mutex
	SubscribeChannel chan *SubscribeMessage
	subscribers      []*SubscribeMessage
}

type Response struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Timestamp string `json:"timestamp"`
	Error     error
}

func NewOraclePool() *OraclePool {

	pool := &OraclePool{
		connections:      map[string]*SafeConnection{},
		lock:             sync.Mutex{},
		SubscribeChannel: make(chan *SubscribeMessage, 5),
		subscribers:      make([]*SubscribeMessage, 1),
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
	pool.lock.Lock()
	pool.connections[url] = &SafeConnection{
		url:             url,
		conn:            connection,
		connectionCount: 1,
		lock:            sync.Mutex{},
	}
	pool.lock.Unlock()

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
	}
}

func (pool *OraclePool) listen() {

	for {

		select {
		case sub := <-pool.SubscribeChannel:
			pool.lock.Lock()
			pool.subscribers = append(pool.subscribers, sub)
			pool.lock.Unlock()
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

	/*



		go func() {
			for {
				res, err := safeConn.readMessage()
				if err != nil {
					log.Println(err)
					return
				}
				resChannel <- *res
				if message.IsSubOnce {
					safeConn.decCount()
					if safeConn.connectionCount == 0 {
						err := safeConn.conn.Close()
						log.Println(err)
					}
					return
				}
			}
		}()

		return err*/
}
