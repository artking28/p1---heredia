package parser

import (
	"errors"
	"exps-heredia/lexer"
	"exps-heredia/models"
	"exps-heredia/utils"
)

type (
	Parser struct {
		Filename   string               `json:"filename"`
		OutputFile string               `json:"outputFile"`
		Tokens     []models.Token       `json:"tokens"`
		Variables  map[string]*Variable `json:"variables"`
		Cursor     int                  `json:"cursor"`
	}
)

func NewParser(filename, output string) (*Parser, error) {

	tokens, err := lexer.Tokenize(filename)
	if err != nil {
		return nil, err
	}

	return &Parser{
		Filename:   filename,
		OutputFile: output,
		Tokens:     tokens,
		Variables:  map[string]*Variable{},
		Cursor:     0,
	}, nil
}

func (this *Parser) GetVarsOffset() uint16 {
	return 100
}

func (this *Parser) At() utils.Pos {
	return this.Tokens[this.Cursor].Pos
}

func (this *Parser) Get(n int) *models.Token {
	if this.Cursor+n >= len(this.Tokens) {
		return nil
	}
	return &this.Tokens[this.Cursor+n]
}

func (this *Parser) Consume(n int) {
	if this.Cursor+n >= len(this.Tokens) {
		return
	}
	this.Cursor += n
}

func (this *Parser) GetFirstAfter(afterOf ...models.TokenKind) (*models.Token, error) {
	all := map[models.TokenKind]bool{}
	for _, t := range afterOf {
		all[t] = true
	}

	token := this.Get(1)
	for i := 1; token != nil; i++ {
		if all[token.Kind] {
			token = this.Get(i)
			continue
		}
		return token, nil
	}
	return nil, errors.New("no token has been found")
}

const (
	NoSpaceMode = iota
	OptionalSpaceMode
	MandatorySpaceMode
)

func (this *Parser) HasNextConsume(spaceMode int, fillOf models.TokenKind, kinds ...models.TokenKind) (*models.Token, error) {
	if spaceMode < NoSpaceMode || spaceMode > MandatorySpaceMode {
		panic("invalid argument in function 'HasNextConsume', unknown space mode")
	}
	if len(kinds) <= 0 {
		panic("invalid argument in function 'HasNextConsume', kinds is null or empty")
	}
	for hasSpace := false; ; {
		token := this.Get(0)
		if token == nil {
			// Fim dos tokens sem encontrar um tipo esperado
			return nil, errors.New("no token has been found")
		}

		for _, kind := range kinds {
			if token.Kind == kind {
				// Se espaços eram obrigatórios mas não foram encontrados, falha
				if spaceMode == MandatorySpaceMode && !hasSpace {
					return nil, errors.New("rule expects spaces but none has been found")
				}
				this.Consume(1)
				return token, nil
			}
		}

		if token.Kind == fillOf {
			// Se espaços não eram permitidos
			if spaceMode == NoSpaceMode {
				return nil, errors.New("space(s) has been found when it actually isn't allowed here")
			}
			hasSpace = true
			this.Consume(1)
			continue
		}

		// Se espaços eram obrigatórios e encontrou outro token, falha
		if spaceMode == MandatorySpaceMode {
			return nil, errors.New("rule expects spaces but none has been found")
		}

		return nil, errors.New("unknown error") // Qualquer outro caso não esperado falha
	}
}

func (parser *Parser) ParseScope(scopeType ScopeKind) (ret Scope, err error) {

	ret.StmtBase = StmtBase{
		Parser: parser,
		Title:  "Scope",
		Pos:    parser.At(),
	}
	tk := parser.Get(0)
	for tk != nil && tk.Kind != models.EOF {

		// Parses some statement on root context of the file
		switch tk.Kind {

		// Parses a comment section
		case models.COMMENT_LINE:
			c, e := parser.ParseComment()
			err = errors.Join(e)
			ret.Body.Statements = append(ret.Body.Statements, c)
			break

		// Parses a function
		case models.KEY_FUN:
			f, e := parser.ParseFunction()
			err = errors.Join(e)
			ret.Body.Statements = append(ret.Body.Statements, f)
			break

		// Parses a global variable
		case models.KEY_DEF:

			svd, e := parser.ParseSingleVarDef()
			err = errors.Join(e)
			ret.Body.Statements = append(ret.Body.Statements, svd)
			break

		// Parses a variable definition, assigment or function call
		case models.ID:
			if scopeType == RootScope {
				return ret, utils.GetUnexpectedTokenNoPosErr(parser.Filename, string(tk.Value))
			}

			t, e0 := parser.GetFirstAfter(models.SPACE, models.ID)
			if e0 != nil {
				err = errors.Join(err, utils.GetUnexpectedTokenNoPosErr(parser.Filename, "EOF"))
				break
			}
			// Parses assignments
			if t.Kind == models.ASSIGN ||
					t.Kind == models.ASSIGN_ADD ||
					t.Kind == models.ASSIGN_SUB ||
					t.Kind == models.ASSIGN_MUL ||
					t.Kind == models.ASSIGN_MOD ||
					t.Kind == models.ASSIGN_AND_BIT ||
					t.Kind == models.ASSIGN_OR_BIT ||
					t.Kind == models.ASSIGN_SHIFT_RIGHT ||
					t.Kind == models.ASSIGN_SHIFT_LEFT {

				assignStmt, e := parser.ParseArgAssign(t.Kind)
				err = errors.Join(err, e)
				ret.Body.Statements = append(ret.Body.Statements, assignStmt)
				break

				// Parses function call
			} else if t.Kind == models.L_PAREN {
				fc, e := parser.ParseFuncCall()
				err = errors.Join(e)
				ret.Body.Statements = append(ret.Body.Statements, fc)
				break

				// Error
			}
			err = errors.Join(err, utils.GetExpectedTokenErr(parser.Filename, "some token to create a variable definition, an assigment or function call", tk.Pos))
			break

		// Ends any kind of AST structure calling the scope parse
		case models.R_BRACE:
			return ret, err

		// Default handler
		default:
			break
		}

		if err != nil {
			return ret, err
		}

		// Advances the parser cursor and update latest token
		parser.Consume(1)
		tk = parser.Get(0)
	}
	return ret, err
}
