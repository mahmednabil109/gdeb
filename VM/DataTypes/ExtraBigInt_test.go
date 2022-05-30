package DataTypes

import (
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"testing"
)

//
//import (
//	"math/big"
//	"math/rand"
//	"testing"
//)
//
//func TestDataWord_Mul(t *testing.T) {
//
//	for i := 0; i < 10; i++ {
//
//		x := generateDataWord()
//		y := generateDataWord()
//
//		xBigInt, _ := big.NewInt(1).SetString(x.ToBinary(), 2)
//		yBigInt, _ := big.NewInt(0).SetString(y.ToBinary(), 2)
//
//		result := mul(x, y)
//
//		resultBigInt := big.NewInt(1)
//		resultBigInt.Mul(xBigInt, yBigInt)
//
//		if uintArrayToBinary(&result) != resultBigInt.Text(2) {
//			t.Error("Multiplication op failed")
//		}
//	}
//}

func uintArrayToBinary(arr []uint32) string {
	result := ""
	for word := 0; word < len(arr); word++ {
		for bit := 0; bit < 32; bit++ {

			if arr[word]&(1<<bit) == 0 {
				result = "0" + result
			} else {
				result = "1" + result
			}
		}
	}

	return result
}

func TestExtraBigInt_SetDataWord(t *testing.T) {
	type args struct {
		byteArr []byte
	}
	tests := []struct {
		name string
		x    ExtraBigInt
		args args
	}{
		// TODO: Add test cases.
		{},
	}

	for i := 1; i < 70; i++ {
		var byteArr []byte
		for j := 1; j <= i; j++ {
			byteArr = append(byteArr, byte(rand.Intn(256)))
		}
		testStruct := struct {
			name string
			x    ExtraBigInt
			args args
		}{
			name: "test" + strconv.Itoa(i),
			x:    NewExtraBigInt(i/4 + 1),
			args: args{
				byteArr: byteArr,
			},
		}
		tests = append(tests, testStruct)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.x.SetDataWord(tt.args.byteArr)

			s1 := tt.x.ToBinary()
			s2 := ByteArrToBinary(tt.args.byteArr)
			s1 = TrimZeros(s1)
			s2 = TrimZeros(s2)

			if s1 != s2 {
				t.Error("Not Equal")
			}
		})
	}
}

func TrimZeros(s string) string {
	if len(s) == 0 {
		return s
	}
	b := []byte(s)
	for b[0] == '0' {
		b = b[1:]
	}
	return string(b)
}

func ByteArrToBinary(x []byte) string {
	result := ""
	for i := 0; i < len(x); i++ {
		for bit := 0; bit < 8; bit++ {

			if x[i]&(1<<bit) == 0 {
				result = "0" + result
			} else {
				result = "1" + result
			}
		}
	}
	return result
}

func generateDataWord(size int) ExtraBigInt {
	dataWord := NewExtraBigInt(size + 1)

	for i := 0; i < size; i++ {
		dataWord.SetUint32(uint32(rand.Intn(10)), uint(i))
	}
	return dataWord
}

func generateTests(i int) (a ExtraBigInt, b ExtraBigInt, aBig *big.Int, bBig *big.Int) {

	a = generateDataWord(i)
	b = generateDataWord(i)
	aBig, _ = new(big.Int).SetString(a.ToBinary(), 2)
	bBig, _ = new(big.Int).SetString(b.ToBinary(), 2)

	return
}

func TestExtraBigInt_Add(t *testing.T) {
	type args struct {
		y ExtraBigInt
	}

	tests := []struct {
		name              string
		x                 ExtraBigInt
		args              args
		ExtraBigIntResult string
		BigInt            string
	}{
		// TODO: Add test cases.
		{},
	}

	for i := 1; i < 70; i++ {
		a, b, aBig, bBig := generateTests(i)
		result := a.Add(b)
		resultBigInt := new(big.Int)
		resultBigInt.Add(aBig, bBig)
		tStruct := struct {
			name              string
			x                 ExtraBigInt
			args              args
			ExtraBigIntResult string
			BigInt            string
		}{
			// TODO: Add test cases.
			name: "Test" + strconv.Itoa(i),
			x:    a,
			args: args{
				y: b,
			},
			ExtraBigIntResult: TrimZeros(result.ToBinary()),
			BigInt:            TrimZeros(resultBigInt.Text(2)),
		}
		tests = append(tests, tStruct)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.BigInt != tt.ExtraBigIntResult {
				fmt.Println(tt.BigInt)
				fmt.Println(tt.ExtraBigIntResult)
				t.Errorf("Add() = %v, ExtraBigIntResult %v", tt.BigInt, tt.ExtraBigIntResult)
			}
		})
	}
}

func TestExtraBigInt_Mul(t *testing.T) {
	type args struct {
		y ExtraBigInt
	}

	tests := []struct {
		name              string
		x                 ExtraBigInt
		args              args
		ExtraBigIntResult string
		BigInt            string
	}{
		// TODO: Add test cases.
		{},
	}

	for i := 1; i < 70; i++ {
		a, b, aBig, bBig := generateTests(i)
		result := mul(a, b)
		resultBigInt := new(big.Int)
		resultBigInt.Mul(aBig, bBig)
		tStruct := struct {
			name              string
			x                 ExtraBigInt
			args              args
			ExtraBigIntResult string
			BigInt            string
		}{
			// TODO: Add test cases.
			name: "Test" + strconv.Itoa(i),
			x:    a,
			args: args{
				y: b,
			},
			ExtraBigIntResult: TrimZeros(uintArrayToBinary(result)),
			BigInt:            TrimZeros(resultBigInt.Text(2)),
		}
		tests = append(tests, tStruct)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.BigInt)
			fmt.Println(tt.ExtraBigIntResult)
			if tt.BigInt != tt.ExtraBigIntResult {
				t.Errorf("Add() = %v, ExtraBigIntResult %v", tt.BigInt, tt.ExtraBigIntResult)
			}
		})
	}
}

func TestExtraBigInt_Sub(t *testing.T) {
	type args struct {
		y ExtraBigInt
	}

	tests := []struct {
		name              string
		x                 ExtraBigInt
		args              args
		ExtraBigIntResult string
		BigInt            string
	}{
		// TODO: Add test cases.
		{},
	}

	for i := 1; i < 70; i++ {
		a, b, aBig, bBig := generateTests(i)
		r1 := a.Sub(b)
		r2 := new(big.Int)
		r2.Sub(aBig, bBig)
		r1String := r1.ToBinary()
		r2String := r2.Text(2)
		if len(r2String) > 0 && r2String[0] == '-' {
			r2String = r2String[1:]
			one := NewExtraBigInt(len(a))
			one[0] = 1
			r1String = r1.Not().Add(one).ToBinary()
		}

		tStruct := struct {
			name              string
			x                 ExtraBigInt
			args              args
			ExtraBigIntResult string
			BigInt            string
		}{
			// TODO: Add test cases.
			name: "Test" + strconv.Itoa(i),
			x:    a,
			args: args{
				y: b,
			},
			ExtraBigIntResult: TrimZeros(r1String),
			BigInt:            TrimZeros(r2String),
		}
		tests = append(tests, tStruct)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.BigInt)
			fmt.Println(tt.ExtraBigIntResult)
			if tt.BigInt != tt.ExtraBigIntResult {
				t.Errorf("Sub() = %s, ExtraBigIntResult %s", tt.BigInt, tt.ExtraBigIntResult)
			}
		})
	}
}
