package parser

import (
	"fmt"
)

type (
	Ast struct {
		Statements []Stmt `json:"statements"`
	}

	Variable struct {
		Id   uint16 `json:"id"`
		Name string `json:"name"`
		StmtBase
		IExp
	}

	Scope struct {
		Id   uint64    `json:"id"`
		Kind ScopeKind `json:"kind"`
		Body Ast       `json:"body"`
		StmtBase
	}

	ScopeKind int
)

const (
	RootScope ScopeKind = iota
	FuncScope
)

func (this Scope) GetTitle() string {
	//TODO implement me
	return "Scope"
}

func (this Scope) WriteMemASM() (string, error) {
	var ret string
	for _, stmt := range this.Body.Statements {
		asm, err := stmt.WriteMemASM()
		if err != nil {
			return "", err
		}
		ret += asm
	}
	return ret, nil
}

func (this Variable) GetTitle() string {
	//TODO implement me
	return "Variable"
}

func (this Variable) WriteMemASM() (ret string, err error) {

	if this.IExp.Count() == 1 {
		v, _ := this.IExp.Resolve()
		id := this.Id + this.Parser.GetVarsOffset()
		ret = fmt.Sprintf("CPY m%d, %d\nGET m%d\n", id, v, id)
	} else {
		ret, err = this.IExp.WriteMemASM()
		if err != nil {
			return "", err
		}
		ret += fmt.Sprintf("CPY m%d, m99\n", this.Id+this.Parser.GetVarsOffset())
	}
	return ret, nil
}

func NewVariable(id uint16, name string, value IExp, base StmtBase) *Variable {
	if value == nil {
		return nil
	}
	return &Variable{Id: id, Name: name, IExp: value, StmtBase: base}
}
