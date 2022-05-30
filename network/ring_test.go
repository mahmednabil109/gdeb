package network

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/mahmednabil109/gdeb/network/utils"
)

var (
	id = ID(utils.ParseID("a09b0ce42948043810a1f2cc7e7079aec7582f29"))
	a  = ID(utils.ParseID("a09b0ce42948043810a1f2cc7e7079aec7582f20"))
	b  = ID(utils.ParseID("a09b0ce42948043810a1f2cc7e7079aec7582fff"))
	_b = ID(utils.ParseID("00000000000000000001"))
	c  = ID(utils.ParseID("4818927895bba843f9b86cd89ab44f885a558f17"))
	d  = ID(utils.ParseID("4818927895bba843f9b86cd89ab44f885a558f18"))
	e  = ID(utils.ParseID("48b88659ffe9828f7f69ab99273bbbcad1b5bcc5"))
)

func Test_inlxrange(t *testing.T) {

	if !id.InLXRange(a, b) {
		t.Fatalf("id in the range but the function returns false")
	}

	if a.InLXRange(id, b) {
		t.Fatal("id is not in the range but the function returns true")
	}

	if !id.InLXRange(a, _b) {
		t.Fatal("id is in the range but the function return false")
	}

	if !id.InLXRange(ZERO_ID, MAX_ID) {
		t.Fatal("id in the range but the function returns false")
	}

	if ZERO_ID.InLXRange(ZERO_ID, MAX_ID) {
		t.Fatal("id is not in the range but the function returns true")
	}

	if !ZERO_ID.InLXRange(id, ZERO_ID) {
		t.Fatal("id is in range but the function returns false")
	}

	if ZERO_ID.InLXRange(ZERO_ID, _b) {
		t.Fatal("id is not in the range but the function returns true")
	}

	if !ZERO_ID.InLXRange(ZERO_ID, ZERO_ID) {
		t.Fatal("in closed ring all values must exists")
	}

	if !a.InLXRange(ZERO_ID, ZERO_ID) {
		t.Fatal("in closed ring all values must exists")
	}

	if !b.InLXRange(ZERO_ID, ZERO_ID) {
		t.Fatal("in closed ring all values must exists")
	}

	if !_b.InLXRange(ZERO_ID, ZERO_ID) {
		t.Fatal("in closed ring all values must exists")
	}

	if !d.InLXRange(c, e) {
		t.Fatal("in range but INXLRange returned false")
	}

}

func Test_equal(t *testing.T) {
	if !a.Equal(a) {
		t.Fatal("equal returns false when it should return true")
	}

	if a.Equal(b) {
		t.Fatal("equal returns true when it should return false")
	}

	if ZERO_ID.Equal(_b) {
		t.Fatal("equal returns true when it should return false")

	}
}

func Test_leftshift(t *testing.T) {
	_id, _ := id.LeftShift()
	_id_str := hex.EncodeToString(_id)
	if _id_str != "09b0ce42948043810a1f2cc7e7079aec7582f290" {
		t.Fatal("leftshift dose not work correctly for hex base")
	}

	copy(_id, id)
	for i := 0; i < 40; i++ {
		_id, _ = _id.LeftShift()
	}
	t.Log(_id)
	if !_id.Equal(ZERO_ID) {
		t.Fatal("40 shift must result in ZERO_ID")
	}
}

func Test_topshift(t *testing.T) {
	_id := id.TopShift(id)
	_id_str := hex.EncodeToString(_id)
	if _id_str != "09b0ce42948043810a1f2cc7e7079aec7582f29a" {
		t.Fatal("topShift must replace lower digit by the top digit of the other id")
	}

	_a := make([]byte, len(a))

	copy(_id, id)
	copy(_a, a)
	for i := 0; i < 40; i++ {
		_id = _id.TopShift(_a)
		_a, _ = ID(_a).LeftShift()
	}

	if !_id.Equal(a) {
		t.Fatalf("after 40 topShifts the id must equal to the other one %v", _id)
	}
}

func Test_addone(t *testing.T) {
	cases := []struct {
		id   ID
		from int
		want string
	}{
		{utils.ParseID("a09b0ce42948043810a1f2cc7e7079aec7582f2f"), 0, "a09b0ce42948043810a1f2cc7e7079aec7582f30"},
		{id, 0, "a09b0ce42948043810a1f2cc7e7079aec7582f2a"},
		{b, 0, "a09b0ce42948043810a1f2cc7e7079aec7583000"},
		{b, 2, "a09b0ce42948043810a1f2cc7e7079aec75830ff"},
		{_b, 0, "0000000000000000000000000000000000000002"},
		{MAX_ID, 0, ZERO_ID.String()},
		{ZERO_ID, 0, "0000000000000000000000000000000000000001"},
		{ZERO_ID, 1, "0000000000000000000000000000000000000010"},
		{ZERO_ID, 2, "0000000000000000000000000000000000000100"},
	}

	for _, tt := range cases {
		t.Run(fmt.Sprintf("AddOne with %s from %d", tt.id.String(), tt.from), func(t *testing.T) {
			res := tt.id.AddOne(tt.from).String()

			if tt.want != res {
				t.Fatalf("add one faild: %s != %s", res, tt.want)
			}
		})
	}
}

func Test_masklowerwith(t *testing.T) {
	cases := []struct {
		id, x ID
		n     int
		want  string
	}{
		{id, a, 0, "a09b0ce42948043810a1f2cc7e7079aec7582f29"},
		{id, a, 2, "a09b0ce42948043810a1f2cc7e7079aec7582fa0"},
		{id, a, 4, "a09b0ce42948043810a1f2cc7e7079aec758a09b"},
		{id, a, 1, "a09b0ce42948043810a1f2cc7e7079aec7582f2a"},
		{id, a, 3, "a09b0ce42948043810a1f2cc7e7079aec7582a09"},
		{id, a, 5, "a09b0ce42948043810a1f2cc7e7079aec75a09b0"},
	}

	for _, tt := range cases {
		t.Run(fmt.Sprintf("mask length %d", tt.n), func(t *testing.T) {
			res := tt.id.MaskLowerWith(tt.x, tt.n).String()
			if tt.want != res {
				t.Fatalf("MaskLower bits dose not work correctly: %s != %s", res, tt.want)
			}
		})
	}
}
