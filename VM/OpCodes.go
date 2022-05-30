package VM

type OPCODE byte

var opcodeMap = map[string]byte{
	"ADD":         1,
	"SUB":         2,
	"MUL":         3,
	"DIV":         4,
	"GT":          5,
	"LT":          6,
	"SGT":         7,
	"SLT":         8,
	"EQ":          9,
	"AND":         10,
	"OR":          11,
	"NOT":         12,
	"XOR":         13,
	"ISZERO":      14,
	"POP":         16,
	"JUMP":        17,
	"JUMPI":       18,
	"SUBSCRIBE":   19,
	"UMSUBSCRIBE": 20,
	"FETCHDATA":   21,
	"ALLOCATE":    22,
}

const (
	ADD         OPCODE = 0x01
	SUB         OPCODE = 0x02
	MUL         OPCODE = 0x03
	DIV         OPCODE = 0x04
	GT          OPCODE = 0x05
	LT          OPCODE = 0x06
	SGT         OPCODE = 0x07
	SLT         OPCODE = 0x08
	EQ          OPCODE = 0x09
	AND         OPCODE = 0x0a
	OR          OPCODE = 0x0b
	NOT         OPCODE = 0x0c
	XOR         OPCODE = 0x0d
	ISZERO      OPCODE = 0x0e
	POP         OPCODE = 0x0f
	JUMP        OPCODE = 0x10
	JUMPI       OPCODE = 0x11
	SUBSCRIBE   OPCODE = 0x12
	UMSUBSCRIBE OPCODE = 0x13
	FETCHDATA   OPCODE = 0x14
	ALLOCATE    OPCODE = 0x15
	PUSH32      OPCODE = 51
	PUSHBOOl    OPCODE = 52
	PUSHTIME    OPCODE = 53
	PUSH64      OPCODE = 54
	PUSH256     OPCODE = 55
)
