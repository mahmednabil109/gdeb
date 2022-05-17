package OracleConnection

import (
	json2 "encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type SafeConnection struct {
	conn            *websocket.Conn
	connectionCount int
	lock            sync.Mutex
}

func (conn *SafeConnection) writeJson(message Message) error {
	json, err := json2.Marshal(message)
	if err != nil {
		log.Println(err)
		return err
	}
	conn.lock.Lock()
	defer conn.lock.Unlock()
	err = conn.conn.WriteJSON(json)
	if err != nil {
		return err
	}
	return nil
}

func (conn *SafeConnection) readMessage() {
	conn.lock.Lock()
	//_, message, err := conn.conn.ReadJSON()
	//if

}
