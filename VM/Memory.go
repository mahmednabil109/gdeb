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

func (memory Memory) store(offset int, elements []uint8) bool {
	size := len(elements)
	if size > 0 && size+offset < max {
		copy(memory[offset:offset+size], elements)
	}
	return false
}

func (memory Memory) load(offset, size int) ([]byte, error) {
	if size+offset >= max {
		return nil, errors.New("memory out of bound")
	}

	return memory[offset : offset+size], nil

}

func (memory Memory) loadString(offset, size int) (string, error) {

	if size+offset >= max {
		return "", errors.New("memory out of bound")
	}

	return string(memory[offset : offset+size]), nil
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
