package utils

func Add8(x, y, carry byte) (sum, carryOut byte) {
	sum16 := uint16(x) + uint16(y) + uint16(carry)
	sum = uint8(sum16)
	carryOut = uint8(sum16 >> 8)
	return
}
