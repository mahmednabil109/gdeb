package TimeListener

import (
	"github.com/mahmednabil109/gdeb/Messages"
	"sync"
	"time"
)

type Frequency byte

const (
	ONCE    Frequency = 0
	HOURLY  Frequency = 1
	DAILY   Frequency = 2
	WEEKLY  Frequency = 3
	MONTHLY Frequency = 4
	YEARLY  Frequency = 5
)

type SubscribeMsg struct {
	Id           int
	Time         []byte
	Frequency    Frequency
	ResponseChan chan *Messages.BroadcastMsg
	Idx          int
}

type UnsubscribeMsg struct {
	Id int
}

type subscriber struct {
	Time         *time.Time
	Frequency    Frequency
	ResponseChan chan *Messages.BroadcastMsg
	Idx          int
}

type TimeListener struct {
	Subscribers     map[int][]*subscriber
	subscribeChan   chan *SubscribeMsg
	unsubscribeChan chan *UnsubscribeMsg
	mutex           sync.Mutex
}

func NewTimeListener() *TimeListener {
	timeListener := &TimeListener{
		Subscribers:     make(map[int][]*subscriber),
		subscribeChan:   make(chan *SubscribeMsg, 5),
		unsubscribeChan: make(chan *UnsubscribeMsg, 5),
	}
	go timeListener.handleSubscribe()
	go timeListener.handleUnsubscribe()
	go timeListener.handlePublishing()

	return timeListener
}

func (listener *TimeListener) Subscribe(msg *SubscribeMsg) {
	listener.subscribeChan <- msg
}

func (listener *TimeListener) handleSubscribe() {
	for {
		select {
		case msg := <-listener.subscribeChan:
			id := msg.Id
			newSubscriber := &subscriber{
				Idx:          msg.Idx,
				Time:         arrToTime(msg.Time),
				ResponseChan: msg.ResponseChan,
				Frequency:    msg.Frequency,
			}
			listener.mutex.Lock()
			if list, isOk := listener.Subscribers[id]; isOk {
				list = append(list, newSubscriber)
			} else {
				listener.Subscribers[id] = make([]*subscriber, 0)
				listener.Subscribers[id] = append(listener.Subscribers[id], newSubscriber)
			}
			listener.mutex.Unlock()
		}
	}
}

func (listener *TimeListener) handleUnsubscribe() {
	for {
		select {
		case msg := <-listener.unsubscribeChan:
			id := msg.Id
			listener.mutex.Lock()
			if _, isOk := listener.Subscribers[id]; isOk {
				delete(listener.Subscribers, id)
			}
			listener.mutex.Unlock()
		}
	}
}

func (listener *TimeListener) handlePublishing() {
	for {
		for k, list := range listener.Subscribers {
			toBeRemoved := make([]int, 0)
			listener.mutex.Lock()
			for i, v := range list {
				if time.Now().UTC().After(*v.Time) || time.Now().UTC().Equal(*v.Time) {
					if v.Frequency == ONCE {
						toBeRemoved = append(toBeRemoved, i)
					} else {
						*v.Time = nextTime(v.Frequency, *(v.Time))
					}
					v.ResponseChan <- &Messages.BroadcastMsg{
						Type:      0,
						Value:     []byte("OK"),
						Key:       "Time",
						Index:     v.Idx,
						Error:     false,
						Timestamp: time.Now().String(),
					}
				}
			}
			for _, i := range toBeRemoved {
				list = removeFromSlice(list, i)
			}
			if len(list) == 0 {
				delete(listener.Subscribers, k)
			}
			listener.mutex.Unlock()

		}
		time.Sleep(3 * time.Second)
	}
}

func nextTime(frequency Frequency, now time.Time) time.Time {
	switch frequency {
	case HOURLY:
		return now.Add(time.Hour)
	case DAILY:
		return now.AddDate(0, 0, 1)
	case WEEKLY:
		return now.AddDate(0, 0, 7)
	case MONTHLY:
		return now.AddDate(0, 1, 0)
	case YEARLY:
		return now.AddDate(1, 0, 0)
	default:
		return now
	}
}

func removeFromSlice(s []*subscriber, i int) []*subscriber {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func arrToTime(arr []byte) *time.Time {
	var year = uint16(arr[1])<<8 + uint16(arr[0])
	t := time.Date(int(year), time.Month(int(arr[2])), int(arr[3]), int(arr[4]), int(arr[5]), int(arr[6]), 0, time.UTC)
	return &t
}
