package VM

type Operation struct {
	execute        OperationType
	stackArgsCount int // number of arguments needed for the operation
	codeArgsCount  uint
	gasPrice       uint64
	pcJump         uint
}

const (
	lowGasPrice       = 2
	midGasPrice       = 4
	highGasPrice      = 7
	superHighGasPrice = 15
)

const (
	onePCJump = 1
)

type JumpTable [100]Operation

func newJumpTable() *JumpTable {
	var oppArray = new(JumpTable)
	(*oppArray)[ADD] =
		Operation{
			execute:        AddOP,
			stackArgsCount: 2,
			gasPrice:       lowGasPrice,
			pcJump:         onePCJump,
		}
	(*oppArray)[SUB] = Operation{
		execute:        SubOP,
		stackArgsCount: 2,
		gasPrice:       lowGasPrice,
		pcJump:         onePCJump,
	}
	(*oppArray)[MUL] = Operation{
		execute:        MulOP,
		stackArgsCount: 2,
		gasPrice:       midGasPrice,
		pcJump:         onePCJump,
	}
	(*oppArray)[DIV] = Operation{}
	//(*oppArray)[GT] = Operation{
	//	run:        GreaterOp,
	//	stackArgsCount: 2,
	//	gasPrice:       lowGasPrice,
	//	pcJump:         onePCJump,
	//}
	(*oppArray)[OR] = Operation{
		execute:        OrOP,
		stackArgsCount: 2,
		gasPrice:       lowGasPrice,
		pcJump:         onePCJump,
	}
	(*oppArray)[XOR] = Operation{
		execute:        XorOP,
		stackArgsCount: 2,
		gasPrice:       lowGasPrice,
		pcJump:         onePCJump,
	}
	(*oppArray)[AND] = Operation{
		execute:        AndOP,
		stackArgsCount: 2,
		gasPrice:       lowGasPrice,
		pcJump:         onePCJump,
	}
	(*oppArray)[NOT] = Operation{
		execute:        NotOP,
		stackArgsCount: 2,
		gasPrice:       lowGasPrice,
		pcJump:         onePCJump,
	}

	(*oppArray)[POP] = Operation{
		execute:        PopOp,
		stackArgsCount: 0,
		gasPrice:       lowGasPrice,
		pcJump:         1,
	}
	(*oppArray)[JUMP] = Operation{
		execute:        JumpOp,
		stackArgsCount: 1,
		gasPrice:       lowGasPrice,
		pcJump:         0,
	}
	(*oppArray)[JUMPI] = Operation{
		execute:        JumpIOp,
		stackArgsCount: 2,
		gasPrice:       lowGasPrice,
		pcJump:         0,
	}

	(*oppArray)[SUBSCRIBE] = Operation{
		execute:        SubscribeOp,
		stackArgsCount: 4,
		gasPrice:       highGasPrice,
		pcJump:         onePCJump,
	}

	(*oppArray)[PUSH32] = Operation{
		execute:        PushOp,
		stackArgsCount: 0,
		gasPrice:       lowGasPrice,
		codeArgsCount:  4,
		pcJump:         0,
	}
	(*oppArray)[PUSHBOOl] = Operation{
		execute:        PushOp,
		stackArgsCount: 0,
		gasPrice:       lowGasPrice,
		codeArgsCount:  4,
		pcJump:         0,
	}
	(*oppArray)[PUSHTIME] = Operation{
		execute:        PushOp,
		stackArgsCount: 0,
		gasPrice:       lowGasPrice,
		codeArgsCount:  7,
		pcJump:         0,
	}
	(*oppArray)[PUSH64] = Operation{
		execute:        PushOp,
		stackArgsCount: 0,
		gasPrice:       lowGasPrice,
		codeArgsCount:  8,
		pcJump:         0,
	}
	(*oppArray)[PUSH256] = Operation{
		execute:        PushOp,
		stackArgsCount: 0,
		gasPrice:       lowGasPrice,
		codeArgsCount:  32,
		pcJump:         0,
	}
	return oppArray
}

func (Map *JumpTable) getInstruction(index byte) Operation {
	return Map[index]
}
