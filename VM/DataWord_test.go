package VM

import (
	"math/big"
	"math/rand"
	"testing"
)

func TestDataWord_Mul(t *testing.T) {

	for i := 0; i < 10; i++ {

		x := generateDataWord()
		y := generateDataWord()

		xBigInt, _ := big.NewInt(1).SetString(x.ToBinary(), 2)
		yBigInt, _ := big.NewInt(0).SetString(y.ToBinary(), 2)

		result := mul(x, y)

		resultBigInt := big.NewInt(1)
		resultBigInt.Mul(xBigInt, yBigInt)

		if uintArrayToBinary(&result) != resultBigInt.Text(2) {
			t.Error("Multiplication op failed")
		}
	}
}

func uintArrayToBinary(arr *[16]uint32) string {
	result := ""
	for word := 0; word < len(*arr); word++ {
		for bit := 0; bit < 32; bit++ {

			if (*arr)[word]&(1<<bit) == 0 {
				result = "0" + result
			} else {
				result = "1" + result
			}
		}
	}
	for i := 0; i < len(result); {
		if result[i] == '0' {
			result = result[1:]
		} else {
			break
		}
	}

	return result
}

func generateDataWord() DataWord {
	dataWord := NewDataWord()

	for i := 0; i < len(dataWord); i++ {
		dataWord.SetUint32(rand.Uint32(), uint(i))
	}
	return dataWord
}

func main() {

}
