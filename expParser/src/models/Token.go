package models

import (
	"exps-heredia/utils"
	"fmt"
	"strconv"
)

type Token struct {
	Pos    utils.Pos `json:"-"`
	Kind   TokenKind `json:"kind"`
	Value  []rune    `json:"value"`
	Repeat int       `json:"repeat"`
}

func NewToken(pos utils.Pos, kind TokenKind, repeat int, value ...rune) Token {
	return Token{Pos: pos, Kind: kind, Value: value, Repeat: repeat}
}

func (this Token) IsSignal() bool {
	return this.Kind == MUL ||
		this.Kind == MOD ||
		this.Kind == ADD ||
		this.Kind == SUB ||
		this.Kind == SHIFT_LEFT ||
		this.Kind == SHIFT_RIGHT ||
		this.Kind == AND_BIT ||
		this.Kind == OR_BIT
}

func ResolveTokenId(filename string, token Token) (Token, error) {
	if token.Kind != ID {
		return token, nil
	}
	value := string(token.Value)

	if tk := FindKeyword(value); tk != UNKNOW {
		return NewToken(token.Pos, tk, 1, token.Value...), nil
	}

	if n, err := strconv.ParseInt(value, 0, 64); err == nil {
		return NewToken(token.Pos, NUMBER, 1, []rune{rune(n)}...), nil
	}

	return token, nil
}

func FindKeyword(word string) TokenKind {
	switch word {
	case "def":
		return KEY_DEF
	case "func":
		return KEY_FUN
	default:
		return UNKNOW
	}
}

func (this *Token) String() string {
	s := this.Kind.String()
	v := string(this.Value)
	if this.Kind == BREAK_LINE {
		v = "\\n"
	} else if this.Kind == TAB {
		v = "\\t"
	} else if this.Kind == EOF {
		v = "\\0"
	} else if this.Kind == NUMBER {
		v = strconv.Itoa(int(this.Value[0]))
	}
	return fmt.Sprintf("Token{%s, \"%s\", %.2d}", s, v, this.Repeat)
}
