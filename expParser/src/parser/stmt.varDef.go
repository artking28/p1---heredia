package parser

import (
	"exps-heredia/models"
	"exps-heredia/utils"
)

func NewVariableStmt(name string, pos utils.Pos, value IExp, parser *Parser) *Variable {
	base := StmtBase{
		Parser: parser,
		Title:  "VariableStmt",
		Pos:    pos,
	}
	ret := NewVariable(uint16(len(parser.Variables)+1), name, value, base)
	parser.Variables[name] = ret
	return ret
}

func (parser *Parser) ParseSingleVarDef() (ret *Variable, err error) {
	waitColon, nameTk := true, parser.Get(0)
	if nameTk == nil {
		return nil, utils.GetUnexpectedTokenNoPosErr(parser.Filename, "EOF")
	}
	if nameTk.Kind == models.KEY_DEF {
		parser.Consume(1)
		waitColon = false
	}
	nameTk, err = parser.HasNextConsume(OptionalSpaceMode, models.SPACE, models.ID)
	if nameTk == nil {
		return nil, utils.GetExpectedTokenErr(parser.Filename, "variable name", parser.At())
	}
	if waitColon {
		if _, err = parser.HasNextConsume(OptionalSpaceMode, models.SPACE, models.EQUAL); err != nil {
			return nil, utils.GetExpectedTokenErr(parser.Filename, "colon token", parser.At())
		}
	} else {
		if _, err = parser.HasNextConsume(OptionalSpaceMode, models.SPACE, models.ASSIGN); err != nil {
			return nil, utils.GetExpectedTokenErr(parser.Filename, "assign token", parser.At())
		}
	}
	parser.Consume(1)
	value, err := parser.ParseExpression(false)
	if err != nil {
		return nil, err
	}

	return NewVariableStmt(string(nameTk.Value), parser.At(), value, parser), nil
}
