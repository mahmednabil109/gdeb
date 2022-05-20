package OracleConnection

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type SafeConnection struct {
	Url   string
	conn  *websocket.Conn
	Count int
	mutex sync.Mutex
}

func (conn *SafeConnection) writeMsg(message string) error {

	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	err := conn.conn.WriteMessage(websocket.TextMessage, []byte(message))

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (conn *SafeConnection) readMessage() (*BroadcastMsg, error) {
	var response *BroadcastMsg
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	err := conn.conn.ReadJSON(response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (conn *SafeConnection) incCount() {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	conn.Count++
}

func (conn *SafeConnection) close() {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	err := conn.conn.Close()
	if err != nil {
		log.Println("Error in closing Safe Conn:", err)
		return
	}
}

func (conn *SafeConnection) decCount() {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	if conn.Count > 0 {
		conn.Count--
	}
}
