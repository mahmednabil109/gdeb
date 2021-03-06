package VM

import (
	"errors"
	"fmt"
	"github.com/mahmednabil109/gdeb/Listeners/OracleListener"
	"github.com/mahmednabil109/gdeb/Listeners/TimeListener"
	"github.com/mahmednabil109/gdeb/Messages"
	"log"
)

type Frequency byte

// Interpreter TODO Blocked is initially false
type Interpreter struct {
	Id                 int
	state              *State
	ContractCode       []byte
	OraclePool         *OracleListener.OraclePool
	TimeListener       *TimeListener.TimeListener
	ReceiveChan        chan *Messages.BroadcastMsg
	TransactionChan    chan interface{}
	getStatusChan      chan bool
	jumpTable          *JumpTable
	gasLimit           uint64
	oracleTransactions map[int]*Messages.BroadcastMsg
	reservedIndex      []int
	IsBlocked          bool
	IsTimeDependent    bool
}

func NewInterpreter(id int, contractCode []byte, pool *OracleListener.OraclePool, gasLimit uint64, timeListener *TimeListener.TimeListener) *Interpreter {
	return &Interpreter{
		Id:                 id,
		state:              newState(),
		ContractCode:       contractCode,
		OraclePool:         pool,
		TimeListener:       timeListener,
		ReceiveChan:        make(chan *Messages.BroadcastMsg),
		TransactionChan:    make(chan interface{}),
		getStatusChan:      make(chan bool),
		jumpTable:          NewJumpTable(),
		gasLimit:           gasLimit,
		oracleTransactions: make(map[int]*Messages.BroadcastMsg),
		reservedIndex:      make([]int, 0),
		IsBlocked:          false,
	}
}

// Run TODO when error happens --> unsubscribe from the oraclePool
func (interpreter *Interpreter) Run() error {

	for {

		consumedGas, pc := &interpreter.state.consumedGas, &interpreter.state.pc
		gasLimit := interpreter.gasLimit

		curInstruction := (interpreter.ContractCode)[*pc]
		operationInfo := interpreter.jumpTable.getInstruction(curInstruction)
		fmt.Println("Current Instruction:", OpcodeIntToString[curInstruction])
		fmt.Println("Consumed Gas", interpreter.state.consumedGas)
		if *consumedGas+operationInfo.gasPrice > gasLimit {
			log.Println("Out of gas error")
			log.Println("Consumed gas:", *consumedGas)
			return errors.New("out of gas error")
		}

		//stack underflow exception
		if interpreter.state.Stack.Size() < operationInfo.stackArgsCount {
			return errors.New("stack underflow error")
		}

		if curInstruction == 0 {
			interpreter.OraclePool.Unsubscribe(&Messages.UnsubscribeMsg{
				VM: interpreter.Id,
			})
			interpreter.TimeListener.Unsubscribe(interpreter.Id)
			log.Println("Consumed gas:", *consumedGas)
			return nil
		}
		err := operationInfo.execute(interpreter)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("----------->")
		for i := len(*interpreter.state.Stack) - 1; i >= 0; i-- {
			fmt.Println((*interpreter.state.Stack)[i])
		}
		fmt.Println()
		//Todo check that the operation is not transfer
		*consumedGas += operationInfo.gasPrice

		*pc += operationInfo.pcJump

		if int(*pc) >= len(interpreter.ContractCode) {
			return errors.New("cannot reach the globalData")
		}
	}
}
