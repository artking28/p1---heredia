package parser

import (
	"exps-heredia/models"
	"exps-heredia/utils"
)

type CommentStmt struct {
	Value string `json:"value"`
	StmtBase
}

func NewCommentStmt(content string, pos utils.Pos, parser *Parser) *CommentStmt {
	return &CommentStmt{
		Value: content,
		StmtBase: StmtBase{
			Parser: parser,
			Title:  "CommentStmt",
			Pos:    pos,
		},
	}
}

func (this CommentStmt) WriteMemASM() (string, error) {
	return "", nil
}

func (parser *Parser) ParseComment() (*CommentStmt, error) {
	var comment string
	h0 := parser.Get(0)
	if h0 == nil {
		return nil, utils.GetUnexpectedTokenNoPosErr(parser.Filename, "EOF")
	}
	if h0.Kind != models.COMMENT_LINE {
		return nil, utils.GetUnexpectedTokenErr(parser.Filename, string(h0.Value), h0.Pos)
	}
	parser.Consume(1)
	for here := parser.Get(0); here != nil && here.Kind != models.BREAK_LINE; here = parser.Get(0) {
		comment += string(here.Value)
		parser.Consume(1)
	}
	return NewCommentStmt(comment, h0.Pos, parser), nil
}
