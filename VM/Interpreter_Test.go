package VM

import (
	"fmt"
	"github.com/mahmednabil109/gdeb/Listeners/OracleListener"
	"github.com/mahmednabil109/gdeb/Listeners/TimeListener"
	"log"
	"os"
	"testing"
)

func TestInterpreter_Run(t *testing.T) {
	type fields struct {
		id           int
		code         []byte
		pool         *OracleListener.OraclePool
		gaslimit     uint64
		timeListener *TimeListener.TimeListener
	}
	//tests := []struct {
	//	name    string
	//	fields  fields
	//}{
	//	{
	//		name:   "test1",
	//		fields: fields{
	//			id:           0,
	//			code:         nil,
	//			pool:         nil,
	//			gaslimit:     0,
	//			timeListener: nil,
	//		},
	//	},
	//}
	//
	//
}

func Test1() {
	f, err := os.Create("./Logs/Test1_Logs")
	defer f.Close()
	log.SetOutput(f)
	log.Println("Started Execution!!!")
	url := "127.0.0.1:8383"
	k1 := "RicePrice"
	k2 := "Stock"
	c := Read("./VM/Programs/Program1")
	fmt.Println("Size", len(c), len(url))
	c = append(c, []byte(url)...)
	fmt.Println("Size", len(c), len(k1))
	c = append(c, []byte(k1)...)
	fmt.Println("Size", len(c), len(k2))
	c = append(c, []byte(k2)...)
	pool := OracleListener.NewOraclePool()
	i := NewInterpreter(0, c, pool, 250, nil, nil)

	err = i.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Test2() {
	f, err := os.Create("./Logs/Test2_Logs")
	defer f.Close()
	log.SetOutput(f)
	log.Println("Started Execution!!!")
	a := "a"
	b := "b"
	c := Read("./VM/Programs/Program2")
	c = append(c, a...)
	c = append(c, b...)
	timeListener := TimeListener.NewTimeListener()
	i := NewInterpreter(0, c, nil, 110, timeListener, nil)

	err = i.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Test3() {
	f, err := os.Create("./Logs/Test3_Logs")
	defer f.Close()
	log.SetOutput(f)
	log.Println("Started Execution!!!")
	url := "127.0.0.1:8383"
	key := "RicePrice"
	a := "a"
	b := "b"
	c := Read("./VM/Programs/Program3")
	c = append(c, []byte(url)...)
	c = append(c, []byte(key)...)
	c = append(c, a...)
	c = append(c, b...)
	fmt.Println(url, len(url))
	fmt.Println(key, len(key))
	fmt.Println(a, len(a))
	fmt.Println(b, len(b))
	timeListener := TimeListener.NewTimeListener()
	pool := OracleListener.NewOraclePool()
	i := NewInterpreter(0, c, pool, 150, timeListener, nil)

	err = i.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Test4() {
	f, err := os.Create("./Logs/Test4_Logs")
	defer f.Close()
	log.SetOutput(f)
	log.Println("Started Execution!!!")
	url := "127.0.0.1:8383"
	key := "Stock"
	a := "a"
	b := "b"
	c := Read("./VM/Programs/Program4")
	c = append(c, []byte(url)...)
	c = append(c, []byte(key)...)
	c = append(c, a...)
	c = append(c, b...)
	fmt.Println(url, len(url))
	fmt.Println(key, len(key))
	fmt.Println(a, len(a))
	fmt.Println(b, len(b))
	timeListener := TimeListener.NewTimeListener()
	pool := OracleListener.NewOraclePool()
	i := NewInterpreter(0, c, pool, 150, timeListener, nil)

	err = i.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}
