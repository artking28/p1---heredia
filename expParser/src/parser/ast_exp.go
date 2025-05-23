package parser

import (
	"exps-heredia/models"
	"exps-heredia/utils"
	"exps-heredia/utils/asmUtils"
	"fmt"
	"strings"
)

type (
	IExp interface {
		Stmt
		Resolve() (uint16, error)
		Count() int
		String() string
	}

	ExpChain struct {
		StmtBase
		All    []IExp           `json:"all"`
		Father *ExpChain        `json:"father"`
		Signal models.TokenKind `json:"oper"`
	}

	VExp uint16

	IdExp uint16
)

func (this IdExp) WriteMemASM() (string, error) {
	return fmt.Sprintf("%s %s\n", asmUtils.GET, this.String()), nil
}

func (this IdExp) Count() int {
	return 1
}

func (this IdExp) String() string {
	return fmt.Sprintf("m%d", this)
}

func (this VExp) String() string {
	return fmt.Sprintf("%d", this)
}

func (this VExp) WriteMemASM() (string, error) {
	return fmt.Sprintf("%s %d\n", asmUtils.GET, this), nil
}

func (this ExpChain) String() (ret string) {
	if this.Count() == 1 {
		return this.All[0].String()
	}
	ret = fmt.Sprintf("( %s", this.All[0].String())
	for _, v := range this.All[1:] {
		ret += fmt.Sprintf(" %s %s", this.Signal.GetSymbol(), v.String())
	}
	return fmt.Sprintf("%s )", ret)
}

func (this ExpChain) WriteMemASM() (string, error) {
	if len(this.All) == 0 {
		return "", nil
	}

	var sb strings.Builder

	s, err := this.All[0].WriteMemASM()
	if err != nil {
		return "", nil
	}

	// Resolve o primeiro
	sb.WriteString(s)
	sb.WriteString("SET m99\n")

	for _, exp := range this.All[1:] {
		s, err = exp.WriteMemASM()
		if err != nil {
			return "", nil
		}
		sb.WriteString(s)
		sb.WriteString(fmt.Sprintf("%s m99\n", this.Signal.GetSymbol()))
		sb.WriteString("SET m99\n")
	}

	return sb.String(), nil
}

func NewExp(values []IExp, father *ExpChain, operator models.TokenKind, base StmtBase) *ExpChain {
	return &ExpChain{
		All:      values,
		Signal:   operator,
		Father:   father,
		StmtBase: base,
	}
}

func NewExpP(exps []IExp, father *ExpChain, oper models.TokenKind, pos utils.Pos, parser *Parser) *ExpChain {
	return NewExp(exps, father, oper, StmtBase{
		Parser: parser,
		Title:  "Exp",
		Pos:    pos,
	})
}

func (this *ExpChain) DeriveInclusiveExp(signal models.TokenKind) (*ExpChain, error) {
	last := this.All[len(this.All)-1]
	e := NewExp([]IExp{last}, this, signal, this.StmtBase)
	this.All[len(this.All)-1] = e
	return e, nil
}

func (this *ExpChain) AddTerm(term IExp) {
	this.All = append(this.All, term)
}

func (this *ExpChain) RootFather() *ExpChain {
	if this.Father != nil {
		return this.Father.RootFather()
	}
	return this
}

func NewIdExp(value uint16) *IdExp {
	n := IdExp(value)
	return &n
}

func NewVExp(value uint16) *VExp {
	n := VExp(value)
	return &n
}

func (this VExp) GetTitle() string {
	return "VExp"
}

func (this IdExp) GetTitle() string {
	return "IdExp"
}

func (this IdExp) Resolve() (uint16, error) {
	return 0, nil
}

func (this VExp) Resolve() (uint16, error) {
	return uint16(this), nil
}

func (this VExp) Count() int {
	return 1
}

func (this ExpChain) Resolve() (uint16, error) {
	//TODO implement me
	panic("implement me | Exp@Resolve")
}

func (this ExpChain) Count() int {
	return len(this.All)
}
