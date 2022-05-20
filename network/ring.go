package network

import (
	"encoding/hex"

	"github.com/mahmednabil109/gdeb/network/utils"
)

// Fixed Identifer space of 160 bits
//* [Note] this works for hex base ids
// the identifier is a []byte each entry contains 2 digits
// TODO make it works for generic id radixes
var (
	MAX_ID  ID = []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
	ZERO_ID ID = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
)

type ID []byte

// function checks if the id is in left execlusive range
// id is in the range (a, b]
func (id ID) InLXRange(a, b ID) bool {
	if a.Less(b) {
		return a.Less(id) && (id.Less(b) || id.Equal(b))
	} else {
		return id.InLXRange(a, MAX_ID) || id.Equal(ZERO_ID) || id.InLXRange(ZERO_ID, b)
	}
}

func (id ID) InLRXRange(a, b ID) bool {
	return id.InLXRange(a, b) && !id.Equal(b)
}

// checks that id is befor a in the ring
func (id ID) Less(a ID) bool {
	if len(id) == len(a) {
		for idx, digit := range id {
			if digit != a[idx] {
				return digit < a[idx]
			}
		}
	}
	return len(id) < len(a)
}

func (id ID) Equal(a ID) bool {
	for idx, digit := range id {
		if digit != a[idx] {
			return false
		}
	}
	return true
}

func (id ID) LeftShift() (ID, byte) {
	var (
		_id   []byte = make([]byte, len(id))
		carry byte
	)

	copy(_id, id)
	carry = ((_id[0] >> 0x04) & 0x0f)
	for i, digit := range _id {
		_id[i] = ((digit << 0x04) & 0xf0)
		if i != len(_id)-1 {
			_id[i] |= ((_id[i+1] >> 0x04) & 0x0f)
		}
	}
	return _id, carry
}

func (id ID) TopShift(a ID) ID {
	_id, _ := id.LeftShift()
	_id[len(_id)-1] |= ((a[0] >> 0x04) & 0x0f)
	return _id
}

// changes the lower i bits of id with the heighest i bits of x
func (id ID) MaskLowerWith(x ID, i int) ID {
	_id := make([]byte, len(id))
	copy(_id, id)

	if i%2 == 0 {
		i /= 2
		copy(_id[len(_id)-i:], x[:i])
	} else {
		i = i/2 + 1
		_tmp := make([]byte, 0)

		// get the higher digits from x
		_tmp = append(_tmp, (x[0]>>0x04)&0x0f)
		for j := 1; j < i; j++ {
			_digit := (x[j-1] << 0x04) | (x[j] >> 0x04)
			_tmp = append(_tmp, _digit)
		}

		// set the digits to the lower part of id
		_id[len(_id)-len(_tmp)] &= 0xf0
		_id[len(_id)-len(_tmp)] |= _tmp[0]

		copy(_id[len(_id)-len(_tmp)+1:], _tmp[1:])
	}
	return _id
}

// AddOne adds one starting from a specific digit and carrys up
func (id ID) AddOne(from int) ID {
	var (
		_id   []byte = make([]byte, len(id))
		carry byte   = 1
	)
	copy(_id, id)

	// special case odd digits
	if from%2 == 1 {
		carry = 1 << 0x04
	}

	// start adding and carring if nessacery
	for i := len(id) - 1 - from/2; i >= 0; i-- {
		_id[i], carry = utils.Add8(_id[i], 0, carry)

		// quick exit
		if carry == 0 {
			break
		}
	}

	return _id
}

func (id ID) String() string {
	return hex.EncodeToString(id)
}
