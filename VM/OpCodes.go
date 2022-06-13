package VM

type OPCODE byte

var opcodeStringToInt = map[string]byte{
	"Halt":                0,
	"ADD":                 1,
	"SUB":                 2,
	"MUL":                 3,
	"DIV":                 4,
	"GT":                  10,
	"LT":                  11,
	"GT_EQ":               12,
	"LT_EQ":               13,
	"SGT":                 14,
	"SLT":                 15,
	"EQ":                  16,
	"AND":                 17,
	"OR":                  18,
	"NOT":                 19,
	"XOR":                 20,
	"ISZERO":              21,
	"ISNEGATIVE":          22,
	"POP":                 30,
	"PUSH32":              31,
	"PUSHBOOl":            32,
	"PUSHTIME":            33,
	"PUSH64":              34,
	"PUSH256":             35,
	"PUSHSTRING":          36,
	"LOAD":                37,
	"STORE":               38,
	"JUMP":                40,
	"JUMPI":               41,
	"AddToPC":             42,
	"SUBSCRIBE":           50,
	"EXECUTEPERIODICALLY": 51,
	"WAIT_FOR_DATA":       52,
	"TRANSFER":            53,
}

var OpcodeIntToString = map[byte]string{
	0:  "Halt",
	1:  "ADD",
	2:  "SUB",
	3:  "MUL",
	4:  "DIV",
	10: "GT",
	11: "LT",
	12: "GT_EQ",
	13: "LT_EQ",
	14: "SGT",
	15: "SLT",
	16: "EQ",
	17: "AND",
	18: "OR",
	19: "NOT",
	20: "XOR",
	21: "ISZERO",
	22: "ISNEGATIVE",
	30: "POP",
	31: "PUSH32",
	32: "PUSHBOOl",
	33: "PUSHTIME",
	34: "PUSH64",
	35: "PUSH256",
	36: "PUSHSTRING",
	37: "LOAD",
	38: "STORE",
	40: "JUMP",
	41: "JUMPI",
	42: "AddToPC",
	50: "SUBSCRIBE",
	51: "EXECUTEPERIODICALLY",
	52: "WAIT_FOR_DATA",
	53: "TRANSFER",
}

const (
	Halt                OPCODE = 0
	ADD                 OPCODE = 1
	SUB                 OPCODE = 2
	MUL                 OPCODE = 3
	DIV                 OPCODE = 4
	GT                  OPCODE = 10
	LT                  OPCODE = 11
	GT_EQ               OPCODE = 12
	LT_EQ               OPCODE = 13
	SGT                 OPCODE = 14
	SLT                 OPCODE = 15
	EQ                  OPCODE = 16
	AND                 OPCODE = 17
	OR                  OPCODE = 18
	NOT                 OPCODE = 19
	XOR                 OPCODE = 20
	ISZERO              OPCODE = 21
	ISNEGATIVE          OPCODE = 22
	POP                 OPCODE = 30
	PUSH32              OPCODE = 31
	PUSHBOOl            OPCODE = 32
	PUSHTIME            OPCODE = 33
	PUSH64              OPCODE = 34
	PUSH256             OPCODE = 35
	PUSHSTRING          OPCODE = 36
	LOAD                OPCODE = 37
	STORE               OPCODE = 38
	JUMP                OPCODE = 40
	JUMPI               OPCODE = 41
	AddToPC             OPCODE = 42
	SUBSCRIBE           OPCODE = 50
	EXECUTEPERIODICALLY OPCODE = 51
	WAIT_FOR_DATA       OPCODE = 52
	TRANSFER            OPCODE = 53
)
