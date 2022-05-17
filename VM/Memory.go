package VM

import (
	"errors"
	"strings"
)

const max int = 20000 * 8

type Memory []byte

func newMemory() Memory {
	return make(Memory, max)
}

func (memory *Memory) store(offset int, elements []uint8) bool {
	size := len(elements)
	if size > 0 && size+offset < max {
		copy((*memory)[offset:offset+size], elements)
	}
	return false
}

func (memory *Memory) load(offset, size int) ([]byte, error) {

	if offset+size < max && size <= 256 {
		return (*memory)[offset : offset+size], nil
	}
	if size > 256 {
		return nil, errors.New("cannot load data of size greater than 256")
	}
	return nil, errors.New("cannot get data")
}

func (memory *Memory) toString() string {
	str := "Stack -----> \n"
	spaceCount := len(str)

	for i := 0; i < 0; i-- {
		val := (*memory)[i]
		str += strings.Repeat(" ", spaceCount)
		str += string(val)
	}
	return str
}
