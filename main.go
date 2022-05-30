package main

import (
	"bufio"
	"crypto/ed25519"
	"fmt"
	"github.com/mahmednabil109/gdeb/VM/DataTypes"
	"github.com/mahmednabil109/gdeb/communication"
	"github.com/mahmednabil109/gdeb/config"
	"github.com/mahmednabil109/gdeb/data"
	"os"
	"strings"
)

var stakeDist map[string]float64
var deployedContracts []string
var privateKey ed25519.PrivateKey
var communNetwCons communication.CommunNetwCons

func setup() {
	config := config.New()
	pk := config.NodeKey()
	privateKey = pk

	data.LoadStakeDist("stakeDistribution.json", &stakeDist)

	communNetwCons = communication.CommunNetwCons{
		ChanNetBlock:        make(chan data.Block),
		ChanNetTransaction:  make(chan data.Transaction),
		ChanConsBlock:       make(chan data.Block),
		ChanConsTransaction: make(chan data.Transaction),
	}

	//same behavior but for deployed contracts, useful for transactions the execute a contract (contract should exist to begin with)
	// data.LoadContracts("deployedContracts.json", &deployedContracts)
}

func generateOpcodes() {
	f, _ := os.Open("ins.txt")
	defer f.Close()
	buffer := bufio.NewScanner(f)

	w, _ := os.Create("generated_opcode.txt")
	defer w.Close()
	i := 1
	for buffer.Scan() {

		line := buffer.Text()

		line = strings.TrimSpace(line)
		lineArr := strings.Split(line, ":")

		ins := lineArr[0][1 : len(lineArr[0])-1]
		out := fmt.Sprintf(`%s OPCODE = 0x%02x`, ins, i)
		i += 1
		w.WriteString(out + "\n")

	}

}

func toBinary32(x uint32) string {
	result := ""
	for bit := 0; bit < 32; bit++ {

		if x&(1<<bit) == 0 {
			result = "0" + result
		} else {
			result = "1" + result
		}
	}
	return result
}
func toBinary64(x uint64) string {
	result := ""
	for bit := 0; bit < 32; bit++ {

		if x&(1<<bit) == 0 {
			result = "0" + result
		} else {
			result = "1" + result
		}
	}
	return result
}
func main() {
	//setup()
	//generateOpcodes()
	// cons := consensus.New(&communNetwCons)
	// netw := network.New(&communNetwCons)

	d1 := DataTypes.NewData(DataTypes.Int64)
	d1.Data[1] = 1 << 31
	d2 := DataTypes.NewData(DataTypes.Int64)
	d2.Data[1] = 1 << 30
	fmt.Println(d1.Data.ToBinary(), d2.Data.ToBinary())

	fmt.Println(d2.Data.Sub(d1.Data).ToBinary())

	var a uint32 = 1<<32 - 1

	var b uint32 = 1<<32 - 1
	fmt.Println(toBinary32(a))
	fmt.Println(toBinary32(b))
	fmt.Println(toBinary32(b + a))

	////code snippet to test ValidateLeader function
	//PublicKey, _ := hex.DecodeString("bd92fd2c61027f602170bf9f6608bc80cabc2f6e6834824fa67dc7fc745cbfe0")
	//PrivateKey, _ := hex.DecodeString("e70b0983a423db62605c527109306d67e16a69d2f4d6641183242e1eac462d27bd92fd2c61027f602170bf9f6608bc80cabc2f6e6834824fa67dc7fc745cbfe0")
	//nonce := 053464
	//proofBytes, _, err := vrf.Prove(PublicKey, PrivateKey, []byte(fmt.Sprint(nonce)))
	//if err != nil {
	//	log.Println(err)
	//}
	//
	//val := consensus.ValidateLeader(nonce, PublicKey, proofBytes, stakeDist)
	//log.Println(val)

}
