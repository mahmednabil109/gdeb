package TimeListener

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Frequency byte

const (
	ONCE        Frequency = 0
	HOURLY      Frequency = 1
	DAILY       Frequency = 2
	WEEKLY      Frequency = 3
	MONTHLY     Frequency = 4
	YEARLY      Frequency = 5
	EveryMinute Frequency = 6
)

//type PeriodicExecution struct {
//	StartTime         time.TimeArr
//	Frequency         Frequency
//	ExecutionInterval time.Duration // in minutes
//}

type SubscribeMsg struct {
	Id            int
	TimeArr       []byte
	Frequency     Frequency
	Interval      time.Duration
	GetStatusChan chan bool
}

type UnsubscribeMsg struct {
	Id int
}

type subscriber struct {
	startTime *time.Time
	isBlocked bool
	*SubscribeMsg
}

type TimeListener struct {
	Subscribers     []*subscriber
	subscribeChan   chan *SubscribeMsg
	unsubscribeChan chan int
	mutex           sync.Mutex
}

func NewTimeListener() *TimeListener {
	timeListener := &TimeListener{
		Subscribers:     make([]*subscriber, 0),
		subscribeChan:   make(chan *SubscribeMsg, 5),
		unsubscribeChan: make(chan int, 5),
	}
	go timeListener.handleSubscribe()
	go timeListener.handleUnsubscribe()
	go timeListener.handlePublishing()

	return timeListener
}

func (listener *TimeListener) Subscribe(msg *SubscribeMsg) {
	listener.subscribeChan <- msg
}

func (listener *TimeListener) Unsubscribe(id int) {
	listener.unsubscribeChan <- id
}

func (listener *TimeListener) handleSubscribe() {
	for {
		select {
		case msg := <-listener.subscribeChan:
			msg.Interval = time.Second * 10
			now := time.Now()
			var startTime time.Time
			fmt.Println("Now", now.String(), "Contract code", *arrToTime(msg.TimeArr))
			if now.After(*arrToTime(msg.TimeArr)) {
				startTime = now
			} else {
				startTime = *arrToTime(msg.TimeArr)
			}

			newSubscriber := &subscriber{
				SubscribeMsg: msg,
				startTime:    &startTime,
				isBlocked:    true,
			}
			listener.mutex.Lock()
			listener.Subscribers = append(listener.Subscribers, newSubscriber)
			listener.mutex.Unlock()
		}
	}
}

func (listener *TimeListener) handleUnsubscribe() {
	for {
		select {
		case msg := <-listener.unsubscribeChan:
			id := msg
			listener.mutex.Lock()
			for i, v := range listener.Subscribers {
				if v.Id == id {
					listener.Subscribers = append(listener.Subscribers[:i], listener.Subscribers[i+1:]...)
					break
				}
			}
			listener.mutex.Unlock()
		}
	}
}

func (listener *TimeListener) handlePublishing() {
	for {
		listener.mutex.Lock()
		for _, v := range listener.Subscribers {
			now := time.Now().UTC()
			if v.isBlocked {
				if now.After(*v.startTime) {
					v.isBlocked = false
					v.GetStatusChan <- false
					fmt.Println("VM is now Free :)")
					log.Println("VM is now Free :)")
				}
			} else {
				if now.After((*v.startTime).Add(v.Interval)) {
					v.isBlocked = true
					v.GetStatusChan <- true
					*v.startTime = getNextTime(v.Frequency, *v.startTime)
					log.Println("VM is Blocked :(")
					fmt.Println("VM is Blocked :(")
				}
			}
		}
		listener.mutex.Unlock()
		time.Sleep(time.Second)
	}
}

func getNextTime(frequency Frequency, now time.Time) time.Time {
	switch frequency {
	case EveryMinute:
		return now.Add(60 * time.Second)
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

func arrToTime(arr []byte) *time.Time {
	var year = uint16(arr[1])<<8 + uint16(arr[0])
	t := time.Date(int(year), time.Month(int(arr[2])), int(arr[3]), int(arr[4]), int(arr[5]), int(arr[6]), 0, time.Local)
	return &t
}
