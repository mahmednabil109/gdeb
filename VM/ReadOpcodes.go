package VM

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GenerateByteCode(fileName string) []byte {

	file, err := os.Open("./VM/Programs/" + fileName)

	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(file)
	code := make([]byte, 0)
	data := make([]byte, 0)
	for scanner.Scan() {
		offsetIdx := 0
		ins := strings.Split(scanner.Text(), " ")
		tmp := make([]byte, 0)
		fmt.Println(ins)
		i := 1
		if ins[0] == "SUBSCRIBE" {
			tmp = append(tmp, opcodeStringToInt["PUSH32"])

			tmp = append(tmp, IntToByte(0)...) // offset
			offsetIdx = len(tmp) - 4
			tmp = append(tmp, opcodeStringToInt["PUSH32"])
			tmp = append(tmp, IntToByte(uint32(len(ins[1])))...) // size
			i++
		}

		for ; i < len(ins); i++ {
			tmp = append(tmp, opcodeStringToInt["PUSH32"])
			stringToInt, _ := strconv.Atoi(ins[i])
			tmp = append(tmp, IntToByte(uint32(stringToInt))...)
		}
		tmp = append(tmp, opcodeStringToInt[ins[0]])
		if ins[0] == "SUBSCRIBE" {
			tmp = append(tmp, GeneratePush(len(ins[1])+1)...)
			tmp = append(tmp, opcodeStringToInt["AddToPC"])
			t := IntToByte(uint32(len(tmp)))
			tmp[offsetIdx] = t[0]
			tmp[offsetIdx+1] = t[1]
			tmp[offsetIdx+2] = t[2]
			tmp[offsetIdx+3] = t[3]
			tmp = append(tmp, []byte(ins[1])...)

		}
		code = append(code, tmp...)
		fmt.Println(tmp)
	}

	return append(code, data...)
}

//func decode(code []byte) {
//	opcodeStringToInt := NewJumpTable()
//
//	for i := 0; i < len(code); {
//		s := ""
//		op := (*opcodeStringToInt)[code[i]]
//
//	}
//
//}

func GeneratePush(a int) []byte {
	result := make([]byte, 0)
	result = append(result, opcodeStringToInt["PUSH32"])
	result = append(result, IntToByte(uint32(a))...)
	return result
}
func incrementIndices(index []int, code []byte, by byte) {
	for _, v := range index {
		code[v] += by
	}
}

func IntToByte(integer uint32) []byte {
	result := make([]byte, 4)

	for i := 0; i < 4; i++ {
		for j := 0; j < 8; j++ {
			if (1<<((i*8)+j))&int(integer) != 0 {
				result[i] = result[i] | (1 << j)
			}
		}
	}

	return result
}
func Read(dir string) []byte {

	file, err := os.Open(dir)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(file)
	i := 0

	code := make([]byte, 0)
	for scanner.Scan() {
		ins := strings.Split(scanner.Text(), " ")
		code = append(code, opcodeStringToInt[ins[0]])
		fmt.Println(ins[0], i)
		i++
		for a, v := range ins {
			if a == 0 {
				continue
			}
			fmt.Println(v, i)
			i++
			num, _ := strconv.Atoi(v)
			code = append(code, byte(num))
		}
	}
	fmt.Println(code)
	return code
}
