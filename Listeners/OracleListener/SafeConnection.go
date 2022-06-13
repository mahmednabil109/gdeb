package OracleListener

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/mahmednabil109/gdeb/Messages"
	"log"
	"sync"
)

type SafeConnection struct {
	Url   string
	Conn  *websocket.Conn
	Count int
	mutex sync.Mutex
}

func (conn *SafeConnection) writeMsg(message string) error {

	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	err := conn.Conn.WriteMessage(websocket.TextMessage, []byte(message))

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (conn *SafeConnection) readMessage() (*Messages.BroadcastMsg, error) {
	var response *Messages.BroadcastMsg
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	err := conn.Conn.ReadJSON(response)
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
	err := conn.Conn.Close()
	fmt.Println("Close!!!!!!!")
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
	if conn.Count == 0 {
		conn.close()
	}

}
