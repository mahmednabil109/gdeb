package VM

type OPCODE byte

var opcodeMap = map[string]byte{
	"ADD":       1,
	"SUB":       2,
	"MUL":       3,
	"DIV":       4,
	"GT":        5,
	"LT":        6,
	"SGT":       7,
	"SLT":       8,
	"EQ":        9,
	"AND":       10,
	"OR":        11,
	"NOT":       12,
	"XOR":       13,
	"ISZERO":    14,
	"PUSH":      15,
	"POP":       16,
	"MLOAD":     17,
	"MSTORE":    18,
	"JUMP":      19,
	"JUMPI":     20,
	"SUBSCRIBE": 21,
	"ALLOCATE":  22,
}

const (
	ADD       OPCODE = 0x01
	SUB       OPCODE = 0x02
	MUL       OPCODE = 0x03
	DIV       OPCODE = 0x04
	GT        OPCODE = 0x05
	LT        OPCODE = 0x06
	SGT       OPCODE = 0x07
	SLT       OPCODE = 0x08
	EQ        OPCODE = 0x09
	ISZERO    OPCODE = 0x0a
	AND       OPCODE = 0x0b
	OR        OPCODE = 0x0c
	NOT       OPCODE = 0x0d
	XOR       OPCODE = 0x0e
	PUSH      OPCODE = 0x0f
	POP       OPCODE = 0x10
	MLOAD     OPCODE = 0x11
	MSTORE    OPCODE = 0x12
	JUMP      OPCODE = 0x13
	JUMPI     OPCODE = 0x14
	SUBSCRIBE OPCODE = 0x15
	ALLOCATE  OPCODE = 0x16
)
