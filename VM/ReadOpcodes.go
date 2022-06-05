package VM

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readOpcodes(fileName string) (res []byte) {

	file, err := os.Open("./Programs/" + fileName)

	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		ins := strings.Split(scanner.Text(), " ")

		res = append(res, opcodeMap[ins[0]])
		if ins[0] == "PUSH" {
			values := make([]byte, 32)
			for i, j := 1, 0; i < len(ins); i, j = i+1, j+1 {
				byteVal, _ := strconv.Atoi(ins[i])
				values[j] = uint8(byteVal)
			}
			res = append(res, values...)
		}
	}
	return res
}

func Read() {

	//byteArr := readOpcodes("Program1")
	//code := make(GlobalData, len(byteArr))
	//code = byteArr
	//fmt.Println(len(code))
	//inter := newInterpreter(&code, 1000)
	//err := inter.run()
	//if err != nil {
	//	return
	//}
}
