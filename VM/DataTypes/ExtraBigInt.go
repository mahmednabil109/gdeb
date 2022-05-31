package DataTypes

import (
	"errors"
	"fmt"
	"math/bits"
	"strconv"
)

type ExtraBigInt []uint32

func NewExtraBigInt(size int) ExtraBigInt {
	return make([]uint32, size)
}

func arrayToDataWord(array []uint32) ExtraBigInt {
	result := NewExtraBigInt(len(array))

	for i, v := range array {
		result[i] = v
	}
	return result
}

func ByteArrToBigInt(arr []byte) ExtraBigInt {
	bigInt := ExtraBigInt{}
	bigInt.SetDataWord(arr)
	return bigInt
}

func (x ExtraBigInt) SetDataWord(byteArr []byte) {
	for i, j := 0, 0; i < len(byteArr); i = i + 1 {
		x[j] = x[j] | (uint32(byteArr[i]) << ((i % 4) * 8))
		if i%4 == 3 {
			j++
		}
	}
}

func intToByte(integer uint32) []byte {
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

func (x ExtraBigInt) toByteArray() []byte {

	result := make([]byte, 1)

	for i := 0; i < len(x); i++ {
		result = append(result, intToByte(x[i])...)
	}
	return result
}

func (x ExtraBigInt) toString() string {
	return string(x.toByteArray())
}
func (x ExtraBigInt) ToBinary() string {
	result := ""
	for word := 0; word < len(x); word++ {
		for bit := 0; bit < 32; bit++ {

			if x[word]&(1<<bit) == 0 {
				result = "0" + result
			} else {
				result = "1" + result
			}
		}
	}
	return result
}

func (x ExtraBigInt) toStringHex() string {
	newX := x.ToBinary()
	xInHex, _ := strconv.ParseInt(newX, 2, 64)
	return fmt.Sprintf("%x", xInHex)
}

func (x ExtraBigInt) ToInt32() uint32 {
	return x[0]
}

func (x ExtraBigInt) toDecimal() string {
	binary := x.ToBinary()
	result := ""
	for i := 0; i < len(binary); i++ {

	}
	return result
}

func (x ExtraBigInt) SetUint32(a uint32, i uint) {
	x[i] = a
}

func (x ExtraBigInt) Add(y ExtraBigInt) ExtraBigInt {
	var carry uint32 = 0
	result := NewExtraBigInt(len(x))
	for i := 0; i < len(result); i++ {
		(result)[i], carry = bits.Add32(x[i], y[i], carry)
	}

	return result
}

func (x ExtraBigInt) Sub(y ExtraBigInt) ExtraBigInt {
	var borrow uint32 = 0
	result := NewExtraBigInt(len(x))
	for i := 0; i < len(result); i++ {
		result[i], borrow = bits.Sub32(x[i], y[i], borrow)
	}
	return result
}

func (x ExtraBigInt) Multiply(y ExtraBigInt) (ExtraBigInt, bool) {
	result := mul(x, y)
	size := len(x)
	ans := result[:len(x)]

	var isOverFlow = false
	for i := size; i < size*2; i++ {
		isOverFlow = isOverFlow || result[i] != 0
	}

	return arrayToDataWord(ans), isOverFlow
}

func mul(x, y ExtraBigInt) []uint32 {
	result := make([]uint32, len(x)*2)
	for Yi := 0; Yi < len(y); Yi++ {
		var carry uint32 = 0
		Ri := Yi
		Xi := 0
		for ; Xi < len(x); Xi = Xi + 1 {
			var lastRes = result[Xi+Ri]

			carry, result[Xi+Ri] = multiplyHelper(lastRes, x[Xi], y[Yi], carry)
		}
		result[Ri+Xi] = carry
	}
	return result
}

func multiplyHelper(z, x, y, carry uint32) (hi, lo uint32) {
	hi, lo = bits.Mul32(x, y)
	lo, carry = bits.Add32(lo, carry, 0)
	hi, _ = bits.Add32(hi, 0, carry)
	lo, carry = bits.Add32(lo, z, 0)
	hi, _ = bits.Add32(hi, 0, carry)
	return hi, lo
}

func MyDiv(u, v []uint32) ([]uint32, []uint32, error) {
	var (
		b            uint64 = 1 << 32
		un, vn, q, r []uint32
		qhat         uint64
		rhat         uint64
		p            uint64
		t, k         int64
		s            int32
	)
	isAllZeros := true
	for i := 0; i < len(v); i++ {
		if v[i] != 0 {
			isAllZeros = false
		}
	}
	if isAllZeros {
		return nil, nil, errors.New("divide by zero")
	}
	m := len(u)
	n := len(v)
	if m < n {
		return nil, nil, errors.New("invalid")
	}

	s = int32(bits.LeadingZeros32(v[n-1]))

	vn = make([]uint32, n)
	for i := n - 1; i > 0; i-- {
		vn[i] = (v[i] << s) | (v[i-1] >> (32 - s))
	}
	vn[0] = v[0] << s

	un = make([]uint32, m+1)
	un[m] = u[m-1] >> (32 - s)
	for i := m - 1; i > 0; i-- {
		un[i] = (u[i] << s) | (u[i-1] >> (32 - s))
	}
	un[0] = u[0] << s

	q = make([]uint32, m-n+1)
	r = make([]uint32, n)

	for j := m - n; j >= 0; j-- {
		qhat = (uint64(un[j+n])*b + uint64(un[j+n-1])) / uint64(vn[n-1])
		rhat = (uint64(un[j+n])*b + uint64(un[j+n-1])) - qhat*uint64(vn[n-1])

		for {
			if qhat >= b || qhat*uint64(vn[n-2]) > b*rhat+uint64(un[j+n-2]) {
				qhat -= 1
				rhat += uint64(vn[n-1])
				if rhat < b {
					continue
				}
			}
			break
		}
		k = 0
		for i := 0; i < n; i++ {
			p = qhat * uint64(vn[i])
			t = int64(un[i+j]) - k - int64((p & 0xFFFFFFFF))
			un[i+j] = uint32(t)
			k = int64(p>>32) - (t >> 32)
		}
		t = int64(un[j+n]) - k
		un[j+n] = uint32(t)

		q[j] = uint32(qhat)

		if t < 0 {
			q[j] = q[j] - 1 // much, add back.
			k = 0
			for i := 0; i < n; i++ {
				t = int64(un[i+j]) + int64(vn[i]) + k
				un[i+j] = uint32(t)
				k = t >> 32
			}
			un[j+n] = un[j+n] + uint32(k)
		}

	}

	for i := 0; i < n-1; i++ {
		r[i] = (un[i] >> s) | (un[i+1] << (32 - s))
		r[n-1] = un[n-1] >> s
	}

	return q, r, nil
}

func (x ExtraBigInt) Div(y ExtraBigInt) []uint32 {
	fmt.Println("y =", y)

	//b := 1 << 32
	var ylen int
	for i := len(y) - 1; i >= 0; i-- {
		if y[i] != 0 {
			ylen = i + 1
			break
		}
	}
	if cap(y) >= ylen+1 {
		y = y[:ylen+1]
	} else {
		y = y[:ylen]
	}
	fmt.Println("y =", y)
	q, _, err := MyDiv(x, y)
	if err != nil {
		return make([]uint32, 1)
	}
	return q
}

func (x ExtraBigInt) Mod(y ExtraBigInt) (result ExtraBigInt) {

	return
}

func (x ExtraBigInt) GT(y ExtraBigInt) bool {
	_, borrow := bits.Sub32(x[0], y[0], 0)
	for i := 1; i < len(x); i++ {
		_, borrow = bits.Sub32(x[i], y[i], borrow)
	}
	return borrow == 0
}

func (x ExtraBigInt) LT(y ExtraBigInt) bool {
	return !x.GT(y) && !x.Eq(y)
}

func (x ExtraBigInt) SLT(y ExtraBigInt) bool {
	dataWordSign := x.sign()
	xSign := y.sign()

	if xSign > dataWordSign {
		return true
	} else if xSign < dataWordSign {
		return false
	} else {
		return x.LT(y)
	}
}

func (x ExtraBigInt) SGT(y ExtraBigInt) bool {
	dataWordSign := x.sign()
	xSign := y.sign()

	if xSign < dataWordSign {
		return true
	} else if xSign > dataWordSign {
		return false
	} else {
		return x.GT(y)
	}
}

/*
	Returns the sign of the dataWord
	if dataWord > 0 return 1
	if dataWord < 0 return -1
	if dataWord == 0 return 0
*/
func (x ExtraBigInt) sign() int {
	if x.IsZero() {
		return 0
	}
	if x[len(x)-1]&1<<31 != 0 {
		return -1
	}
	return 1
}

func (x ExtraBigInt) Eq(y ExtraBigInt) bool {
	isEqual := true
	for i := 0; i < len(x); i++ {
		isEqual = isEqual && x[i] == y[i]
	}
	return isEqual
}

func (x ExtraBigInt) IsZero() bool {
	for i := 0; i < len(x); i++ {
		if x[i] != 0 {
			return false
		}
	}
	return true
}

func (x ExtraBigInt) And(y ExtraBigInt) ExtraBigInt {
	result := NewExtraBigInt(len(x))
	for i := 0; i < len(x); i++ {
		(result)[i] = x[i] & y[i]
	}
	return result
}

func (x ExtraBigInt) Or(y ExtraBigInt) ExtraBigInt {
	result := NewExtraBigInt(len(x))
	for i := 0; i < len(x); i++ {
		(result)[i] = x[i] | y[i]
	}
	return result
}

func (x ExtraBigInt) Not() (result ExtraBigInt) {
	result = NewExtraBigInt(len(x))
	for i := 0; i < len(x); i++ {
		(result)[i] = ^x[i]
	}
	return
}

func (x ExtraBigInt) Xor(y ExtraBigInt) ExtraBigInt {
	result := NewExtraBigInt(len(x))
	for i := 0; i < len(x); i++ {
		(result)[i] = x[i] ^ y[i]
	}
	return result
}
