package parser

import (
	"exps-heredia/models"
	"exps-heredia/utils"
	"exps-heredia/utils/asmUtils"
	"fmt"
)

type AssignStmt struct {
	VariableId uint16 `json:"variableId"`
	ScopeId    uint64 `json:"scopeId"`
	Expression IExp   `json:"expression"`
	StmtBase
}

func NewAssignStmt(id uint16, exp IExp, pos utils.Pos, parser *Parser) *AssignStmt {
	return &AssignStmt{
		VariableId: id,
		Expression: exp,
		StmtBase: StmtBase{
			Parser: parser,
			Title:  "AssignStmt",
			Pos:    pos,
		}}
}

func (this AssignStmt) WriteMemASM() (ret string, err error) {
	if this.Expression.Count() == 1 {
		v, _ := this.Expression.Resolve()
		ret += fmt.Sprintf("CPY m%d, %d\n", this.VariableId+this.Parser.GetVarsOffset(), v)
	} else {
		ret, err = this.Expression.WriteMemASM()
		if err != nil {
			return "", err
		}
		ret += fmt.Sprintf("%s m%d, m99\n", asmUtils.CPY, this.VariableId+this.Parser.GetVarsOffset())
	}
	return ret, nil
}

func (parser *Parser) ParseArgAssign(kind models.TokenKind) (ret *AssignStmt, err error) {
	h0 := parser.Get(0)
	if h0 == nil {
		return nil, utils.GetUnexpectedTokenNoPosErr(parser.Filename, "EOF")
	}

	parser.Consume(1)
	if _, err = parser.HasNextConsume(OptionalSpaceMode, models.SPACE, kind); err != nil {
		return nil, utils.GetExpectedTokenErrOr(parser.Filename, "assignment", err.Error(), h0.Pos)
	}

	parser.Consume(1)
	assignValue, err := parser.ParseExpression(false)
	if err != nil {
		return nil, err
	}

	variable := parser.Variables[string(h0.Value)]
	var exp IExp
	v := NewIdExp(variable.Id + parser.GetVarsOffset())
	switch kind {
	case models.ASSIGN_MUL:
		exp = NewExpP([]IExp{v, assignValue}, nil, models.MUL, h0.Pos, parser)
		break
	case models.ASSIGN_MOD:
		exp = NewExpP([]IExp{v, assignValue}, nil, models.MOD, h0.Pos, parser)
		break
	case models.ASSIGN_ADD:
		exp = NewExpP([]IExp{v, assignValue}, nil, models.ADD, h0.Pos, parser)
		break
	case models.ASSIGN_SUB:
		exp = NewExpP([]IExp{v, assignValue}, nil, models.SUB, h0.Pos, parser)
		break
	case models.ASSIGN_AND_BIT:
		exp = NewExpP([]IExp{v, assignValue}, nil, models.AND_BIT, h0.Pos, parser)
		break
	case models.ASSIGN_OR_BIT:
		exp = NewExpP([]IExp{v, assignValue}, nil, models.OR_BIT, h0.Pos, parser)
		break
	case models.ASSIGN_SHIFT_RIGHT:
		exp = NewExpP([]IExp{v, assignValue}, nil, models.SHIFT_RIGHT, h0.Pos, parser)
		break
	case models.ASSIGN_SHIFT_LEFT:
		exp = NewExpP([]IExp{v, assignValue}, nil, models.SHIFT_LEFT, h0.Pos, parser)
		break
	case models.ASSIGN:
		exp = assignValue
		break
	default:
		panic("implement me | ParseArgAssign switch case")
	}

	if variable == nil {
		return nil, utils.GetUnkownVariableErr(parser.Filename, string(h0.Value), h0.Pos)
	}
	return NewAssignStmt(variable.Id, exp, h0.Pos, parser), nil
}
