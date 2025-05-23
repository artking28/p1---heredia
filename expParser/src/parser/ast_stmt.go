package parser

import (
	"exps-heredia/utils"
)

type (
	Stmt interface {
		WriteMemASM() (string, error)
		GetTitle() string
	}

	StmtBase struct {
		Parser *Parser   `json:"-"`
		Title  string    `json:"title"`
		Pos    utils.Pos `json:"pos"`
	}
)

func (this StmtBase) GetTitle() string {
	return this.Title
}
