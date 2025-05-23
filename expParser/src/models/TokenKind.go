package models

import (
	"exps-heredia/utils/asmUtils"
	"fmt"
)

type TokenKind int

const (
	EOF TokenKind = (iota + 1) * 100
	UNKNOW
	BREAK_LINE
	TAB
	SPACE
	ID
	EQUAL
	NUMBER
	L_PAREN
	R_PAREN
	L_BRACE
	R_BRACE
	KEY_FUN
	KEY_DEF
	SLASH
	COMMENT_LINE
	MUL
	MOD
	SUB
	ADD
	AND_BIT
	OR_BIT
	SHIFT_LEFT
	SHIFT_RIGHT
	ASSIGN
	ASSIGN_MUL
	ASSIGN_MOD
	ASSIGN_SUB
	ASSIGN_ADD
	ASSIGN_AND_BIT
	ASSIGN_OR_BIT
	ASSIGN_SHIFT_LEFT
	ASSIGN_SHIFT_RIGHT
	GREATER_THAN
	LOWER_THAN
)

func (this *TokenKind) String() (s string) {
	switch *this {
	case EOF:
		return "EOF"
	case BREAK_LINE:
		return "BREAK_LINE"
	case TAB:
		return "TAB"
	case SPACE:
		return "SPACE"
	case ID:
		return "ID"
	case NUMBER:
		return "NUMBER"
	case KEY_FUN:
		return "FUN"
	case KEY_DEF:
		return "DEF"
	case L_PAREN:
		return "L_PAREN"
	case R_PAREN:
		return "R_PAREN"
	case L_BRACE:
		return "L_BRACE"
	case R_BRACE:
		return "R_BRACE"
	case SHIFT_LEFT:
		return "SHIFT_LEFT"
	case SHIFT_RIGHT:
		return "SHIFT_RIGHT"
	case ASSIGN:
		return "ASSIGN"
	case MUL:
		return "MUL"
	case MOD:
		return "MOD"
	case ADD:
		return "ADD"
	case SUB:
		return "SUB"
	case AND_BIT:
		return "AND_BIT"
	case OR_BIT:
		return "OR_BIT"
	default:
		s = fmt.Sprintf("UNKNOWN(%d)", *this)
	}
	return s
}

func (this TokenKind) Weight() uint8 {
	switch this {
	case MUL, MOD:
		return 254
	case ADD, SUB:
		return 253
	case SHIFT_LEFT, SHIFT_RIGHT:
		return 252
	case AND_BIT:
		return 251
	case OR_BIT:
		return 250
	case UNKNOW:
		return 1
	default:
		return 0
	}
}

func (this *TokenKind) GetSymbol() string {
	switch *this {
	case ADD:
		return asmUtils.ADD
	case SUB:
		return asmUtils.SUB
	case AND_BIT:
		return asmUtils.AND
	case OR_BIT:
		return asmUtils.OR
	default:
		return "<NONE>"
	}
}

func CombineTokens(tk0, tk1 Token) (TokenKind, []rune) {

	if tk0.Kind == ADD && tk1.Kind == ASSIGN {
		return ASSIGN_ADD, []rune("+=")
	} else if tk0.Kind == SUB && tk1.Kind == ASSIGN {
		return ASSIGN_SUB, []rune("-=")
	} else if tk0.Kind == MOD && tk1.Kind == ASSIGN {
		return ASSIGN_MOD, []rune("%=")
	} else if tk0.Kind == MUL && tk1.Kind == ASSIGN {
		return ASSIGN_MUL, []rune("*=")
	} else if tk0.Kind == AND_BIT && tk1.Kind == ASSIGN {
		return ASSIGN_AND_BIT, []rune("&=")
	} else if tk0.Kind == OR_BIT && tk1.Kind == ASSIGN {
		return ASSIGN_OR_BIT, []rune("|=")
	} else if tk0.Kind == GREATER_THAN && tk1.Kind == GREATER_THAN {
		return SHIFT_RIGHT, []rune(">>")
	} else if tk0.Kind == SHIFT_RIGHT && tk1.Kind == ASSIGN {
		return ASSIGN_SHIFT_RIGHT, []rune(">>=")
	} else if tk0.Kind == LOWER_THAN && tk1.Kind == LOWER_THAN {
		return SHIFT_LEFT, []rune("<<")
	} else if tk0.Kind == SHIFT_LEFT && tk1.Kind == ASSIGN {
		return ASSIGN_SHIFT_LEFT, []rune("<<=")
	} else if tk0.Kind == SLASH && tk1.Kind == SLASH {
		return COMMENT_LINE, []rune("//")
	} else {
		return UNKNOW, []rune("")
	}
}
