package VM

import (
	"errors"
	"github.com/mahmednabil109/gdeb/VM/DataTypes"
	"strconv"
	"strings"
)

const MaxSize int = 20000

type Memory []*DataTypes.DataWord

// TODO Memory Refactor Memory to hold different datatypes.
// TODO newMemory return new memory
func newMemory() Memory {
	return make(Memory, MaxSize)
}

func (memory Memory) store(data *DataTypes.DataWord, idx int) error {
	if idx >= MaxSize {
		return errors.New("illegal index:" + strconv.Itoa(idx))
	}
	memory[idx] = data
	return nil
}

func (memory Memory) load(idx int) (*DataTypes.DataWord, error) {
	if idx >= len(memory) {
		return nil, errors.New("memory out of bound")
	}

	return memory[idx], nil
}

func (memory *Memory) toString() string {
	str := "Memory -----> \n"
	spaceCount := len(str)

	for i := 0; i < 0; i-- {
		val := (*memory)[i]
		str += strings.Repeat(" ", spaceCount)
		str += val.ToString() + "\n"
	}
	return str
}
