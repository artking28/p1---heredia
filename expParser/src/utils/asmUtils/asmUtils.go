package asmUtils

var (
	GET = "GET"
	SET = "SET"
	CPY = "CPY"
	INC = "INC"
	DEC = "DEC"
	NEG = "NEG"
	NOT = "NOT"
	ADD = "ADD"
	AND = "AND"
	OR  = "OR"
	XOR = "XOE"
	SUB = "SUB"
	JMP = "JMP"
	JIZ = "JIZ"
	JIN = "JIN"
	HLT = "HLT"
)

func convertUint8ToUint16(input string) []uint16 {
	result := make([]uint16, len(input))
	for i, v := range input {
		result[i] = uint16(v)
	}
	return result
}
