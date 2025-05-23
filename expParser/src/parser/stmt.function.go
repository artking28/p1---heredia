package parser

import (
	"errors"
	"exps-heredia/models"
	"exps-heredia/utils"
)

type FuncStmt struct {
	Name string `json:"name"`
	Body Scope  `json:"body"`
	StmtBase
}

func (this FuncStmt) WriteMemASM() (string, error) {
	if this.Name != "main" {
		return "", errors.New("only main function can be created")
	}
	return this.Body.WriteMemASM()
}

func NewFuncStmt(name string, body Scope, pos utils.Pos, parser *Parser) *FuncStmt {
	return &FuncStmt{
		Name: name,
		Body: body,
		StmtBase: StmtBase{
			Parser: parser,
			Title:  "FuncStmt",
			Pos:    pos,
		},
	}
}

func (parser *Parser) ParseFunction() (ret *FuncStmt, err error) {
	h0 := parser.Get(0)
	if h0 == nil {
		return nil, utils.GetUnexpectedTokenNoPosErr(parser.Filename, "EOF")
	}
	parser.Consume(1)
	nameTk, err := parser.HasNextConsume(MandatorySpaceMode, models.SPACE, models.ID)
	if nameTk == nil {
		return nil, utils.GetExpectedTokenErr(parser.Filename, "function name", h0.Pos)
	}
	if _, err = parser.HasNextConsume(OptionalSpaceMode, models.SPACE, models.L_PAREN); err != nil {
		return nil, utils.GetExpectedTokenErrOr(parser.Filename, "left parenthesis", err.Error(), h0.Pos)
	}
	if _, err = parser.HasNextConsume(OptionalSpaceMode, models.SPACE, models.R_PAREN); err != nil {
		return nil, utils.GetExpectedTokenErrOr(parser.Filename, "right parenthesis", err.Error(), h0.Pos)
	}
	if _, err = parser.HasNextConsume(OptionalSpaceMode, models.SPACE, models.L_BRACE); err != nil {
		return nil, utils.GetExpectedTokenErrOr(parser.Filename, "left brace", err.Error(), h0.Pos)
	}
	ast, err := parser.ParseScope(FuncScope)
	if err != nil {
		return nil, err
	}
	if _, err = parser.HasNextConsume(OptionalSpaceMode, models.SPACE, models.R_BRACE); err != nil {
		return nil, utils.GetExpectedTokenErrOr(parser.Filename, "right brace", err.Error(), h0.Pos)
	}
	return NewFuncStmt(string(nameTk.Value), ast, h0.Pos, parser), nil
}
