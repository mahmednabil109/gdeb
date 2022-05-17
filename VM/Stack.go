package VM

import "strings"

type Stack []DataWord

func newStack() *Stack {
	return &Stack{}
}

func (stack *Stack) IsEmpty() bool {
	return len(*stack) == 0
}

func (stack *Stack) Push(x DataWord) {
	*stack = append(*stack, x)
}

func (stack *Stack) Pop() DataWord {
	if stack.IsEmpty() {
		return NewDataWord()
	}
	popped := (*stack)[len(*stack)-1] // the popped element
	*stack = (*stack)[:len(*stack)-1] // Removing the element from the stack by slicing
	return popped
}

func (stack *Stack) Peek() (*DataWord, bool) {
	return &(*stack)[(len(*stack) - 1)], !stack.IsEmpty()
}

func (stack *Stack) Size() int {
	return len(*stack)
}

func (stack *Stack) toString() string {
	str := "Stack -----> \n"
	spaceCount := len(str)

	for i := len(*stack) - 1; i >= 0; i-- {
		val := (*stack)[i]
		str += strings.Repeat(" ", spaceCount)
		str += "0x" + val.toStringHex() + "\n"
	}
	return str
}
