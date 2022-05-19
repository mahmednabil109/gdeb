package OracleConnection

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type SafeConnection struct {
	url             string
	conn            *websocket.Conn
	connectionCount int
	lock            sync.Mutex
}

func (conn *SafeConnection) writeMsg(message string) error {

	conn.lock.Lock()
	defer conn.lock.Unlock()
	err := conn.conn.WriteMessage(websocket.TextMessage, []byte(message))

	log.Println("Connection send Msg:", message)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (conn *SafeConnection) readMessage() (*Response, error) {
	var response *Response
	conn.lock.Lock()
	defer conn.lock.Unlock()
	err := conn.conn.ReadJSON(response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (conn *SafeConnection) incCount() {
	conn.lock.Lock()
	defer conn.lock.Unlock()
	conn.connectionCount++
}

func (conn *SafeConnection) decCount() {
	conn.lock.Lock()
	defer conn.lock.Unlock()
	if conn.connectionCount > 0 {
		conn.connectionCount--
	}
}
