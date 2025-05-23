package compiler

import (
    "ASM/models"
)

func ParseAll(parser *models.Parser) error {

    //models.StartHeap(128)
    tk := parser.Get(0)
    for ; tk != nil && tk.Kind != models.TOKEN_EOF; tk = parser.Get(0) {
        switch tk.Kind {
        case models.TOKEN_SLASH:
            err := ParseComment(parser)
            if err != nil {
                return err
            }
            break
        case models.TOKEN_HASHTAG:
            err := ParseLabelDecl(parser)
            if err != nil {
                return err
            }
            break
        case models.TOKEN_INC, models.TOKEN_DEC, models.TOKEN_NEG, models.TOKEN_NOT, models.TOKEN_HLT:
            err := ParsePureInstruction(parser)
            if err != nil {
                return err
            }
            break

        case models.TOKEN_GET, models.TOKEN_SET, models.TOKEN_ADD, models.TOKEN_AND, models.TOKEN_OR, models.TOKEN_XOR, models.TOKEN_SUB:
            err := ParseSingleInstruction(parser)
            if err != nil {
                return err
            }
            break

        case models.TOKEN_JMP, models.TOKEN_JIZ, models.TOKEN_JIN:
            err := ParseJumpInstruction(parser)
            if err != nil {
                return err
            }
            break

        case models.TOKEN_CPY:
            err := ParseDoubleInstruction(parser)
            if err != nil {
                return err
            }
            break
        default:
            break
        }
        parser.Consume(1)
    }

    return nil
}

func ParseComment(parser *models.Parser) error {
    var comment string
    h0 := parser.Get(0)
    if h0 == nil {
        return models.GetUnexpectedTokenNoPosErr(parser.Filename, "EOF")
    }
    if h0.Kind != models.TOKEN_SLASH || (h0.Kind == models.TOKEN_SLASH && h0.Repeat < 2) {
        return models.GetUnexpectedTokenErr(parser.Filename, string(h0.Value), h0.Pos)
    }
    parser.Consume(2)
    for here := parser.Get(0); here != nil && here.Kind != models.TOKEN_BREAK_LINE; here = parser.Get(0) {
        comment += string(here.Value)
        parser.Consume(1)
    }
    parser.Inject(models.NewCommentStmt(comment, h0.Pos, parser))
    return nil
}

func ParseLabelDecl(parser *models.Parser) error {
    h0 := parser.Get(0)
    if h0 == nil {
        return models.GetUnexpectedTokenNoPosErr(parser.Filename, "EOF")
    }
    if h0.Kind != models.TOKEN_HASHTAG {
        return models.GetUnexpectedTokenErr(parser.Filename, string(h0.Value), h0.Pos)
    }
    parser.Consume(1)
    h1 := parser.HasNextConsume(models.NoSpaceMode, models.TOKEN_ID)
    if h1 == nil {
        return models.GetExpectedTokenErr(parser.Filename, "valid identifier", h0.Pos)
    }
    parser.Inject(models.NewLabelDeclStmt(string(h1.Value), h0.Pos, parser))
    return nil
}

func ParseJumpInstruction(parser *models.Parser) error {
    h0 := parser.Get(0)
    if h0 == nil {
        return models.GetUnexpectedTokenNoPosErr(parser.Filename, "EOF")
    }
    parser.Consume(1)
    if parser.HasNextConsume(models.MandatorySpaceMode, models.TOKEN_HASHTAG) == nil {
        return models.GetExpectedTokenErr(parser.Filename, "label", h0.Pos)
    }
    h1 := parser.HasNextConsume(models.NoSpaceMode, models.TOKEN_ID)
    if h1 == nil {
        return models.GetExpectedTokenErr(parser.Filename, "label", h0.Pos)
    }
    parser.Inject(models.NewJumpStmt(string(h1.Value), h0.Kind, h0.Pos, parser))
    return nil
}

func ParsePureInstruction(parser *models.Parser) error {
    h0 := parser.Get(0)
    if h0 == nil {
        return models.GetUnexpectedTokenNoPosErr(parser.Filename, "EOF")
    }
    parser.Inject(models.NewPureInstructionStmt(h0.Kind, h0.Pos, parser))
    parser.Consume(1)
    return nil
}

func ParseSingleInstruction(parser *models.Parser) error {
    h0 := parser.Get(0)
    if h0 == nil {
        return models.GetUnexpectedTokenNoPosErr(parser.Filename, "EOF")
    }
    parser.Consume(1)
    var h1 *models.Token
    phrase := "memory address"
    if h0.Kind == models.TOKEN_GET {
        h1 = parser.HasNextConsume(models.MandatorySpaceMode, models.TOKEN_MEM, models.TOKEN_NUMBER)
        phrase += " or number literal"
    } else {
        h1 = parser.HasNextConsume(models.MandatorySpaceMode, models.TOKEN_MEM)
    }
    if h1 == nil {
        return models.GetExpectedTokenErr(parser.Filename, phrase, h0.Pos)
    }
    parser.Inject(models.NewSingleInstructionStmt(h0.Kind, h0.Pos, *h1, parser))
    return nil
}

// ParseDoubleInstruction parses: Instruction space memAddress space comma space (memAddress|number)
func ParseDoubleInstruction(parser *models.Parser) error {
    h0 := parser.Get(0)
    if h0 == nil {
        return models.GetUnexpectedTokenNoPosErr(parser.Filename, "EOF")
    }
    parser.Consume(1)
    h1 := parser.HasNextConsume(models.MandatorySpaceMode, models.TOKEN_MEM)
    if h1 == nil {
        return models.GetExpectedTokenErr(parser.Filename, "memory address", h0.Pos)
    }
    if k := parser.HasNextConsume(models.OptionalSpaceMode, models.TOKEN_COMMA); k == nil {
        return models.GetExpectedTokenErr(parser.Filename, "comma", h0.Pos)
    }
    h2 := parser.HasNextConsume(models.OptionalSpaceMode, models.TOKEN_MEM, models.TOKEN_NUMBER)
    if h2 == nil {
        return models.GetExpectedTokenErr(parser.Filename, "memory address or number literal", h0.Pos)
    }

    parser.Inject(models.NewDoubleInstructionStmt(h0.Kind, h0.Pos, *h1, *h2, parser))
    return nil
}
