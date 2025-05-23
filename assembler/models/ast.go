package models

import (
    "ASM/neanderExecutor"
    "encoding/json"
    "fmt"
)

type (
    Pos struct {
        Line   int64 `json:"line"`
        Column int64 `json:"column"`
    }

    Ast struct {
        Statements []Stmt `json:"statements"`
    }

    Stmt interface {
        WriteMemASM() ([]uint16, []uint16, error)
        GetTitle() string
    }

    StmtBase struct {
        Parser *Parser `json:"-"`
        Title  string  `json:"title"`
        Pos    Pos     `json:"pos"`
    }

    CommentStmt struct {
        Value string `json:"value"`
        StmtBase
    }

    LabelDeclStmt struct {
        LabelName string `json:"labelName"`
        StmtBase
    }

    JumpStmt struct {
        TargetLabelName string        `json:"TargetLabelName"`
        JumpKind        TokenKindEnum `json:"jumpKind"`
        StmtBase
    }

    PureInstructionStmt struct {
        Code TokenKindEnum `json:"code"`
        StmtBase
    }

    SingleInstructionStmt struct {
        PureInstructionStmt
        Left Token `json:"left"`
    }

    DoubleInstructionStmt struct {
        SingleInstructionStmt
        Right Token `json:"right"`
    }
)

func (this Ast) Inspect() {
    str, err := json.MarshalIndent(this, "", "   ")
    if err != nil {
        panic(err.Error())
    }

    fmt.Printf("%s\n", string(str))
}

func NewCommentStmt(content string, pos Pos, parser *Parser) CommentStmt {
    return CommentStmt{
        Value: content,
        StmtBase: StmtBase{
            Parser: parser,
            Title:  "CommentStmt",
            Pos:    pos,
        },
    }
}

func NewLabelDeclStmt(labelName string, pos Pos, parser *Parser) LabelDeclStmt {
    return LabelDeclStmt{
        LabelName: labelName,
        StmtBase: StmtBase{
            Parser: parser,
            Title:  "LabelDeclStmt",
            Pos:    pos,
        },
    }
}

func NewJumpStmt(targetLabelName string, jumpKind TokenKindEnum, pos Pos, parser *Parser) JumpStmt {
    return JumpStmt{
        TargetLabelName: targetLabelName,
        JumpKind:        jumpKind,
        StmtBase: StmtBase{
            Parser: parser,
            Title:  "JumpStmt",
            Pos:    pos,
        },
    }
}

func NewPureInstructionStmt(code TokenKindEnum, pos Pos, parser *Parser) PureInstructionStmt {
    return PureInstructionStmt{
        Code: code,
        StmtBase: StmtBase{
            Parser: parser,
            Title:  "PureInstructionStmt",
            Pos:    pos,
        },
    }
}

func NewSingleInstructionStmt(code TokenKindEnum, pos Pos, left Token, parser *Parser) SingleInstructionStmt {
    s := SingleInstructionStmt{
        PureInstructionStmt: NewPureInstructionStmt(code, pos, parser),
        Left:                left,
    }
    s.StmtBase.Title = "SingleInstructionStmt"
    return s
}

func (this SingleInstructionStmt) GetLeftASUint16() uint16 {
    return uint16(this.Left.Value[0])
}

func NewDoubleInstructionStmt(code TokenKindEnum, pos Pos, left, right Token, parser *Parser) DoubleInstructionStmt {
    d := DoubleInstructionStmt{
        SingleInstructionStmt: NewSingleInstructionStmt(code, pos, left, parser),
        Right:                 right,
    }
    d.StmtBase.Title = "DoubleInstructionStmt"
    return d
}

func (this DoubleInstructionStmt) GetRightASUint16() uint16 {
    return uint16(this.Right.Value[0])
}

func (this CommentStmt) GetTitle() string {
    return this.Title
}

func (this LabelDeclStmt) GetTitle() string {
    return this.Title
}

func (this JumpStmt) GetTitle() string {
    return this.Title
}

func (this PureInstructionStmt) GetTitle() string {
    return this.Title
}

func (this SingleInstructionStmt) GetTitle() string {
    return this.Title
}

func (this DoubleInstructionStmt) GetTitle() string {
    return this.Title
}

func (this CommentStmt) WriteMemASM() (ret []uint16, reviewOffset []uint16, err error) {
    return []uint16{}, nil, nil
}

func (this LabelDeclStmt) WriteMemASM() (ret []uint16, reviewOffset []uint16, err error) {
    return []uint16{}, nil, nil
}

func (this JumpStmt) WriteMemASM() (ret []uint16, reviewOffset []uint16, err error) {
    switch this.JumpKind {
    case TOKEN_JMP:
        ret = append(ret, neander.JMP, this.Parser.labels[this.TargetLabelName])
        break
    case TOKEN_JIZ:
        ret = append(ret, neander.JZ, this.Parser.labels[this.TargetLabelName])
        break
    case TOKEN_JIN:
        ret = append(ret, neander.JN, this.Parser.labels[this.TargetLabelName])
        break
    default:
        //TODO implement me
        panic("implement me switch default branch in JumpStmt WriteMemASM implementation")
    }
    return ret, nil, nil
}

func (this PureInstructionStmt) WriteMemASM() (ret []uint16, reviewOffset []uint16, err error) {
    switch this.Code {
    case TOKEN_INC:
        ret = append(ret, neander.ADD, OneValue)
        break
    case TOKEN_DEC:
        ret = append(ret, neander.ADD, MinusOneValue)
        break
    case TOKEN_NEG:
        ret = append(ret, neander.NOT, neander.ADD, OneValue)
        break
    case TOKEN_NOT:
        ret = append(ret, neander.NOT)
        break
    case TOKEN_HLT:
        ret = append(ret, neander.HLT)
        break
    default:
        //TODO implement me
        panic("implement me switch default branch in PureInstructionStmt WriteMemASM implementation")
    }
    return ret, nil, nil
}

func (this SingleInstructionStmt) WriteMemASM() (ret []uint16, reviewOffset []uint16, err error) {
    switch this.Code {
    case TOKEN_GET:
        if this.Left.Kind == TOKEN_NUMBER {
            memAddr := this.Parser.AllocNum(int16(this.GetLeftASUint16()))
            ret = append(ret, neander.LDA, memAddr)
            break
        }
        ret = append(ret, neander.LDA, this.GetLeftASUint16())
        reviewOffset = append(reviewOffset, uint16(len(ret)))
        break
    case TOKEN_SET:
        ret = append(ret, neander.STA, this.GetLeftASUint16())
        reviewOffset = append(reviewOffset, uint16(len(ret)))
        break
    case TOKEN_ADD:
        ret = append(ret, neander.ADD, this.GetLeftASUint16())
        reviewOffset = append(reviewOffset, uint16(len(ret)))
        break
    case TOKEN_AND:
        ret = append(ret, neander.AND, this.GetLeftASUint16())
        reviewOffset = append(reviewOffset, uint16(len(ret)))
        break
    case TOKEN_OR:
        ret = append(ret, neander.OR, this.GetLeftASUint16())
        reviewOffset = append(reviewOffset, uint16(len(ret)))
        break
    case TOKEN_XOR:
        ret = append(ret, neander.STA, AcCache1Addr)
        ret = append(ret, neander.OR, this.GetLeftASUint16())
        reviewOffset = append(reviewOffset, uint16(len(ret)))
        ret = append(ret, neander.STA, AcCache0Addr)
        ret = append(ret, neander.LDA, AcCache1Addr)
        ret = append(ret, neander.AND, this.GetLeftASUint16())
        reviewOffset = append(reviewOffset, uint16(len(ret)))
        ret = append(ret, neander.NOT)
        ret = append(ret, neander.AND, AcCache0Addr)
        break
    case TOKEN_SUB:
        ret = append(ret, neander.STA, AcCache0Addr)
        ret = append(ret, neander.LDA, this.GetLeftASUint16())
        reviewOffset = append(reviewOffset, uint16(len(ret)))
        ret = append(ret, neander.NOT)
        ret = append(ret, neander.ADD, OneValue)
        ret = append(ret, neander.ADD, AcCache0Addr)
    default:
        //TODO implement me
        panic("implement me switch default branch in SingleInstructionStmt WriteMemASM implementation")
    }
    return ret, reviewOffset, nil
}

func (this DoubleInstructionStmt) WriteMemASM() (ret []uint16, reviewOffset []uint16, err error) {
    switch this.Code {
    case TOKEN_CPY:
        ret = append(ret, neander.STA, AcCache0Addr)
        if this.Right.Kind == TOKEN_NUMBER {
            memAddr := this.Parser.AllocNum(int16(this.GetRightASUint16()))
            ret = append(ret, neander.LDA, memAddr)
        } else {
            ret = append(ret, neander.LDA, this.GetRightASUint16())
            reviewOffset = append(reviewOffset, uint16(len(ret)))
        }
        ret = append(ret, neander.STA, this.GetLeftASUint16())
        reviewOffset = append(reviewOffset, uint16(len(ret)))
        ret = append(ret, neander.LDA, AcCache0Addr)
        break
    default:
        //TODO implement me
        panic("implement me switch default branch in DoubleInstructionStmt WriteMemASM implementation")
    }
    return ret, reviewOffset, nil
}
