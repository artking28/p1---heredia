package models

import (
	"fmt"
	"strconv"
	"strings"
)

type (
	TokenKindEnum int

	Token struct {
		Pos    Pos           `json:"-"`
		Kind   TokenKindEnum `json:"kind"`
		Value  []rune        `json:"value"`
		Repeat int           `json:"repeat"`
	}
)

func NewToken(pos Pos, kind TokenKindEnum, repeat int, value ...rune) Token {
	return Token{Pos: pos, Kind: kind, Value: value, Repeat: repeat}
}

func (this Token) String(complete bool) string {
	var s string
	switch this.Kind {
	case TOKEN_SPACE:
		s = "TOKEN_SPACE"
		break
	case TOKEN_BREAK_LINE:
		s = "TOKEN_BREAK_LINE"
		break
	case TOKEN_TAB:
		s = "TOKEN_TAB"
		break
	case TOKEN_ID:
		s = "TOKEN_ID"
		break
	case TOKEN_NUMBER:
		s = "TOKEN_NUMBER"
		break
	case TOKEN_COMMA:
		s = "TOKEN_COMMA"
		break
	case TOKEN_COLON:
		s = "TOKEN_COLON"
		break
	case TOKEN_MEM:
		s = "TOKEN_MEM"
		break
	case TOKEN_SLASH:
		s = "TOKEN_SLASH"
		break
	case TOKEN_HASHTAG:
		s = "TOKEN_HASHTAG"
		break
	case TOKEN_EOF:
		s = "TOKEN_EOF"
		break
	case TOKEN_GET:
		s = "TOKEN_GET"
		break
	case TOKEN_SET:
		s = "TOKEN_SET"
		break
	case TOKEN_CPY:
		s = "TOKEN_CPY"
		break
	case TOKEN_INC:
		s = "TOKEN_INC"
		break
	case TOKEN_DEC:
		s = "TOKEN_DEC"
		break
	case TOKEN_NEG:
		s = "TOKEN_NEG"
		break
	case TOKEN_NOT:
		s = "TOKEN_NOT"
		break
	case TOKEN_ADD:
		s = "TOKEN_ADD"
		break
	case TOKEN_AND:
		s = "TOKEN_AND"
		break
	case TOKEN_OR:
		s = "TOKEN_OR"
		break
	case TOKEN_XOR:
		s = "TOKEN_XOR"
		break
	case TOKEN_SUB:
		s = "TOKEN_SUB"
		break
	case TOKEN_JMP:
		s = "TOKEN_JMP"
		break
	case TOKEN_JIZ:
		s = "TOKEN_JIZ"
		break
	case TOKEN_JIN:
		s = "TOKEN_JIN"
		break
	case TOKEN_HLT:
		s = "TOKEN_HLT"
		break
	}
	if complete {
		v := string(this.Value)
		if this.Kind == TOKEN_BREAK_LINE {
			v = "\\n"
		} else if this.Kind == TOKEN_TAB {
			v = "\\t"
		} else if this.Kind == TOKEN_EOF {
			v = "\\0"
		} else if this.Kind == TOKEN_MEM || this.Kind == TOKEN_NUMBER {
			v = strconv.Itoa(int(this.Value[0]))
		}
		return fmt.Sprintf("Token{%s, \"%s\", %.2d}", s, v, this.Repeat)
	}
	return s
}

func AppendToken(tokens *[]Token, token Token) {
	if tokens == nil {
		tokens = &[]Token{}
	}
	count := len(*tokens)
	if count > 0 && (*tokens)[count-1].Kind == token.Kind && string((*tokens)[count-1].Value) == string(token.Value) {
		(*tokens)[count-1].Repeat = (*tokens)[count-1].Repeat + 1
		return
	}
	*tokens = append(*tokens, token)
}

func ResolveTokenId(filename string, token Token) (Token, error) {
	if token.Kind != TOKEN_ID {
		return token, nil
	}
	value := string(token.Value)
	count := len(value)
	if value[0] == 'm' {
		mem, err := strconv.ParseInt(value[1:], 10, 64)
		if err == nil {
			return NewToken(token.Pos, TOKEN_MEM, 1, rune(mem)), nil
		}
	} else if value[0] == '#' {
		return NewToken(token.Pos, TOKEN_LABEL, 1, []rune(value[1:len(value)-1])...), nil
	} else if strings.ToUpper(value) == ("GET") {
		return NewToken(token.Pos, TOKEN_GET, 1, []rune(value)...), nil
	} else if strings.ToUpper(value) == ("SET") {
		return NewToken(token.Pos, TOKEN_SET, 1, []rune(value)...), nil
	} else if strings.ToUpper(value) == ("CPY") {
		return NewToken(token.Pos, TOKEN_CPY, 1, []rune(value)...), nil
	} else if strings.ToUpper(value) == ("INC") {
		return NewToken(token.Pos, TOKEN_INC, 1, []rune(value)...), nil
	} else if strings.ToUpper(value) == ("DEC") {
		return NewToken(token.Pos, TOKEN_DEC, 1, []rune(value)...), nil
	} else if strings.ToUpper(value) == ("NEG") {
		return NewToken(token.Pos, TOKEN_NEG, 1, []rune(value)...), nil
	} else if strings.ToUpper(value) == ("NOT") {
		return NewToken(token.Pos, TOKEN_NOT, 1, []rune(value)...), nil
	} else if strings.ToUpper(value) == ("ADD") {
		return NewToken(token.Pos, TOKEN_ADD, 1, []rune(value)...), nil
	} else if strings.ToUpper(value) == ("AND") {
		return NewToken(token.Pos, TOKEN_AND, 1, []rune(value)...), nil
	} else if strings.ToUpper(value) == ("OR") {
		return NewToken(token.Pos, TOKEN_OR, 1, []rune(value)...), nil
	} else if strings.ToUpper(value) == ("XOR") {
		return NewToken(token.Pos, TOKEN_XOR, 1, []rune(value)...), nil
	} else if strings.ToUpper(value) == ("SUB") {
		return NewToken(token.Pos, TOKEN_SUB, 1, []rune(value)...), nil
	} else if strings.ToUpper(value) == ("JMP") {
		return NewToken(token.Pos, TOKEN_JMP, 1, []rune(value)...), nil
	} else if strings.ToUpper(value) == ("JIZ") {
		return NewToken(token.Pos, TOKEN_JIZ, 1, []rune(value)...), nil
	} else if strings.ToUpper(value) == ("JIN") {
		return NewToken(token.Pos, TOKEN_JIN, 1, []rune(value)...), nil
	} else if strings.ToUpper(value) == ("HLT") {
		return NewToken(token.Pos, TOKEN_HLT, 1, []rune(value)...), nil
	} else {
		num, err := strconv.ParseInt(value[:count], 0, 64)
		if err == nil {
			return NewToken(token.Pos, TOKEN_NUMBER, 1, rune(num)), nil
		}
		//return NewToken(token.Pos, TOKEN_NUMBER, 1, rune(num)), GetUnexpectedTokenErr(filename, string(token.Value), token.Pos)
	}

	return token, nil
}

const (

	// #########################
	//       Normal tokens
	// #########################
	TOKEN_SPACE TokenKindEnum = iota
	TOKEN_BREAK_LINE
	TOKEN_TAB
	TOKEN_ID
	TOKEN_NUMBER
	TOKEN_COMMA
	TOKEN_COLON
	TOKEN_MEM
	TOKEN_LABEL
	TOKEN_SLASH
	TOKEN_HASHTAG
	TOKEN_EOF

	// #########################
	//         MINMONICS
	// #########################

	// Memory manipulations
	TOKEN_GET
	TOKEN_SET
	TOKEN_CPY

	// Simple operations
	TOKEN_INC
	TOKEN_DEC
	TOKEN_NEG
	TOKEN_NOT

	// Operations
	TOKEN_ADD
	TOKEN_AND
	TOKEN_OR
	TOKEN_XOR

	// Loops and comparations
	TOKEN_SUB
	TOKEN_JMP
	TOKEN_JIZ
	TOKEN_JIN

	// Runtime actions
	TOKEN_HLT
)
