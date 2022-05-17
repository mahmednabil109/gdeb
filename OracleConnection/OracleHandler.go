package OracleConnection

import (
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type Message struct {
	OracleKey string `json:"oracle_key"`
	Condition string `json:"condition"`
	Value     string `json:"value"`
	IsSubOnce bool   `json:"isSubOnce"`
}

type OraclePool struct {
	connections     map[string]*SafeConnection
	connectionCount int
	lock            sync.Mutex
}

type Response struct {
	value     string
	timestamp string
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
		conn:            connection,
		connectionCount: 1,
		lock:            sync.Mutex{},
	}
	pool.lock.Unlock()

	return pool.connections[url], nil
}

func (pool *OraclePool) Subscribe(url string, message Message, res chan Response) error {

	safeConnection, err := pool.getConnection(url)
	if err != nil {
		return err
	}
	err = safeConnection.writeJson(message)
	if err != nil {
		log.Println(err)
		return err
	}

	go func() {

	}()
	return err
}
