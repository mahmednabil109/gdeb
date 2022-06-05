package TimeListener

import (
	"fmt"
	"github.com/mahmednabil109/gdeb/Messages"
	"sync"
	"testing"
	"time"
)

func TestNewTimeListener(t *testing.T) {
	var timeSlice [][]byte

	timeListener := NewTimeListener()

	for i := 0; i < 10; i++ {
		timeSlice = append(timeSlice, []byte{})
		now := time.Now().UTC()
		year := intToByte(now.Year())
		timeSlice[i] = append(timeSlice[i], year[0])
		timeSlice[i] = append(timeSlice[i], year[1])
		timeSlice[i] = append(timeSlice[i], byte(now.Month()))
		timeSlice[i] = append(timeSlice[i], byte(now.Day()))
		timeSlice[i] = append(timeSlice[i], byte(now.Hour()))
		timeSlice[i] = append(timeSlice[i], byte(now.Minute()))
		duration := time.Duration(i + 1)
		timeSlice[i] = append(timeSlice[i], byte(now.Add(duration*time.Second).Second()))
	}
	var wg sync.WaitGroup
	trackWg := 0
	testFunc := func(i int) {
		fmt.Println("entered go routine", i)
		trackWg++
		receiveChan := make(chan *Messages.BroadcastMsg)
		sub := &SubscribeMsg{
			Time:         timeSlice[i],
			ResponseChan: receiveChan,
			Idx:          i,
			Id:           i,
		}
		timeListener.Subscribe(sub)
		go func() {
			for {
				select {
				case msg := <-receiveChan:
					fmt.Println("go routine:", i, "received msg", msg)
					wg.Done()
					trackWg--
					break
				}
			}
		}()
	}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go testFunc(i)
	}
	wg.Wait()
	fmt.Println("Track,", trackWg)
}

func intToByte(integer int) []byte {
	result := make([]byte, 4)
	s := ""
	for i := 0; i < 4; i++ {
		for j := 0; j < 8; j++ {
			if (1<<((i*8)+j))&(integer) != 0 {
				result[i] = result[i] | (1 << j)
				s = "1" + s
			} else {
				s = "0" + s
			}
		}
		s = " " + s
	}

	return result
}
