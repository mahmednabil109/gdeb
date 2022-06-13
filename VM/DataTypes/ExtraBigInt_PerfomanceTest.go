package DataTypes

import (
	"math/big"
	"math/rand"
	"testing"
)

type Test struct {
	a ExtraBigInt
	b ExtraBigInt
	x *big.Int
	y *big.Int
}

func GenerateTests256() []Test {

	tests := make([]Test, 0)
	for i := 0; i < 2000; i++ {
		a := generateDataWord1(8)
		b := generateDataWord1(8)
		x, _ := new(big.Int).SetString(a.ToBinary(), 2)
		y, _ := new(big.Int).SetString(b.ToBinary(), 2)
		tests = append(tests, Test{
			a: a,
			b: b,
			x: x,
			y: y,
		})
	}
	return tests
}

func generateDataWord1(size int) ExtraBigInt {
	dataWord := NewExtraBigInt(size)

	for i := 0; i < size; i++ {
		dataWord.SetUint32(rand.Uint32(), uint(i))
	}
	return dataWord
}

var tests = GenerateTests256()

func TestExtraBigInt256(b *testing.B) {
	for _, test := range tests {
		_, _ = test.a.Multiply(test.b)
	}
}

var resultBigInt = new(big.Int)

func TestBigInt256(b *testing.B) {
	for _, test := range tests {
		resultBigInt.Mul(test.x, test.y)
	}
}
