package OracleListener

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/mahmednabil109/gdeb/Messages"
	"log"
	"net/url"
	"sync"
)

var Host = "127.0.0.1:8383/"

type OraclePool struct {
	Connections     map[string]*SafeConnection
	mutex           sync.Mutex
	SubscribeChan   chan *Messages.SubscribeMsg
	UnsubscribeChan chan *Messages.UnsubscribeMsg
	Subscribers     map[int][]*Messages.SubscribeMsg
}

func NewOraclePool() *OraclePool {

	pool := &OraclePool{
		Connections:     make(map[string]*SafeConnection, 0),
		mutex:           sync.Mutex{},
		SubscribeChan:   make(chan *Messages.SubscribeMsg),
		UnsubscribeChan: make(chan *Messages.UnsubscribeMsg),
		Subscribers:     make(map[int][]*Messages.SubscribeMsg, 0),
	}
	go pool.subscribeListener()
	go pool.unsubscribeListener()

	return pool
}

func (pool *OraclePool) Unsubscribe(msg *Messages.UnsubscribeMsg) {
	pool.UnsubscribeChan <- msg
}

func (pool *OraclePool) Subscribe(msg *Messages.SubscribeMsg) {
	pool.SubscribeChan <- msg
}

func (pool *OraclePool) IsSubscribed(id int) bool {
	_, isOk := pool.Subscribers[id]
	return isOk
}

func (pool *OraclePool) unsubscribeListener() {
	for {
		select {
		case msg := <-pool.UnsubscribeChan:
			vmId := msg.VM
			topics := pool.Subscribers[vmId]
			for _, topic := range topics {
				pool.mutex.Lock()
				safeConn := pool.Connections[topic.OracleKey]
				safeConn.decCount()
				pool.mutex.Unlock()
			}
			pool.unsubscribeVM(vmId)
		}
	}
}

func (pool *OraclePool) unsubscribeVM(id int) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()
	delete(pool.Subscribers, id)
}

func (pool *OraclePool) addConnection(topic string, conn *SafeConnection) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()
	pool.Connections[topic] = conn
	fmt.Println("Connnnnection", conn)
}

func (pool *OraclePool) getConnection(oracleUrl, topic string) (*SafeConnection, error) {

	if conn, isOk := pool.Connections[topic]; isOk {
		pool.Connections[topic].incCount()
		return conn, nil
	}
	u := url.URL{Scheme: "ws", Host: oracleUrl, Path: "sub"}

	log.Println("Connected to", u.String())
	connection, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, errors.New("no such Oracle")
	}

	newConn := &SafeConnection{
		Url:   oracleUrl,
		Conn:  connection,
		Count: 1,
		mutex: sync.Mutex{},
	}
	pool.addConnection(topic, newConn)

	go pool.messageReceiver(pool.Connections[topic])

	return pool.Connections[topic], nil
}

func (pool *OraclePool) messageReceiver(safeConn *SafeConnection) {
	conn := safeConn.Conn
	for {
		var res Messages.OracleMsg
		err := conn.ReadJSON(&res)
		if err != nil {
			fmt.Println("Close: " + err.Error())
			break
		} else {
			pool.broadcast(res)
		}
	}
}

func (pool *OraclePool) broadcast(msg Messages.OracleMsg) {

	for _, v := range pool.Subscribers {
		for _, subscribeMsg := range v {
			if isSubscriber := subscribeMsg.OracleKey == msg.Key; isSubscriber {
				subscribeMsg.BroadcastChan <- &Messages.BroadcastMsg{
					Key:       msg.Key,
					Value:     msg.Value,
					Timestamp: msg.Timestamp,
					Index:     subscribeMsg.Index,
					Type:      subscribeMsg.KeyType,
					Error:     msg.Error,
				}
			}
		}
	}
}

func (pool *OraclePool) addSubscriber(subMsg *Messages.SubscribeMsg) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	id := subMsg.VmId
	if topics, isOk := pool.Subscribers[id]; isOk {
		pool.Subscribers[id] = append(pool.Subscribers[id], subMsg)
	} else {
		topics = make([]*Messages.SubscribeMsg, 0)
		topics = append(topics, subMsg)
		pool.Subscribers[id] = topics
	}
}

func (pool *OraclePool) subscribeListener() {
	fmt.Println("Launched the subscribe listener")
	for {
		select {
		case sub := <-pool.SubscribeChan:
			pool.addSubscriber(sub)
			safeConn, err := pool.getConnection(sub.Url, sub.OracleKey)
			if err != nil {
				sub.BroadcastChan <- &Messages.BroadcastMsg{Error: true}
				fmt.Println("Error")
				return
			}
			err = safeConn.writeMsg(sub.OracleKey) // sending message to oracle server
			if err != nil {
				log.Println(err)
			}
		}
	}
}
