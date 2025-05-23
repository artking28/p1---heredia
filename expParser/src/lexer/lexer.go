package lexer

import (
	"errors"
	"exps-heredia/models"
	"exps-heredia/utils"
	"os"
	"strings"
)

func Tokenize(filename string) ([]models.Token, error) {

	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if len(bytes) <= 0 {
		return nil, utils.GetEmptyFileErr(filename)
	}

	var ret []models.Token
	column, line := 0, 1
	runes := []rune(string(bytes))
	buffer := strings.Builder{}
	for _, run := range runes {
		column++
		if 'a' <= run && run <= 'z' ||
			'A' <= run && run <= 'Z' ||
			'0' <= run && run <= '9' {
			buffer.WriteRune(run)
			continue
		}
		if buffer.Len() > 0 {
			var tk models.Token
			pos := utils.Pos{Line: int64(line), Column: int64(column - buffer.Len())}
			tk = models.NewToken(pos, models.ID, 1, []rune(buffer.String())...)
			tk, err = models.ResolveTokenId(filename, tk)
			if err != nil {
				return nil, err
			}
			err = errors.Join(err, AppendToken(filename, &ret, tk))
			buffer.Reset()
		}
		switch run {
		case '\n':
			line += 1
			column = 0
			pos := utils.Pos{Line: int64(line), Column: int64(column)}
			tk := models.NewToken(pos, models.BREAK_LINE, 1, run)
			err = errors.Join(err, AppendToken(filename, &ret, tk))
			break
		case '\t':
			pos := utils.Pos{Line: int64(line), Column: int64(column)}
			tk := models.NewToken(pos, models.TAB, 1, run)
			err = errors.Join(err, AppendToken(filename, &ret, tk))
			break
		case '/':
			pos := utils.Pos{Line: int64(line), Column: int64(column)}
			tk := models.NewToken(pos, models.COMMENT_LINE, 1, run)
			err = errors.Join(err, AppendToken(filename, &ret, tk))
			break
		case '(':
			pos := utils.Pos{Line: int64(line), Column: int64(column)}
			tk := models.NewToken(pos, models.L_PAREN, 1, run)
			err = errors.Join(err, AppendToken(filename, &ret, tk))
			break
		case ')':
			pos := utils.Pos{Line: int64(line), Column: int64(column)}
			tk := models.NewToken(pos, models.R_PAREN, 1, run)
			err = errors.Join(err, AppendToken(filename, &ret, tk))
			break
		case '{':
			pos := utils.Pos{Line: int64(line), Column: int64(column)}
			tk := models.NewToken(pos, models.L_BRACE, 1, run)
			err = errors.Join(err, AppendToken(filename, &ret, tk))
			break
		case '}':
			pos := utils.Pos{Line: int64(line), Column: int64(column)}
			tk := models.NewToken(pos, models.R_BRACE, 1, run)
			err = errors.Join(err, AppendToken(filename, &ret, tk))
			break
		case '=':
			pos := utils.Pos{Line: int64(line), Column: int64(column)}
			tk := models.NewToken(pos, models.ASSIGN, 1, run)
			err = errors.Join(err, AppendToken(filename, &ret, tk))
			break
		case '+':
			pos := utils.Pos{Line: int64(line), Column: int64(column)}
			tk := models.NewToken(pos, models.ADD, 1, run)
			err = errors.Join(err, AppendToken(filename, &ret, tk))
			break
		case '-':
			pos := utils.Pos{Line: int64(line), Column: int64(column)}
			tk := models.NewToken(pos, models.SUB, 1, run)
			err = errors.Join(err, AppendToken(filename, &ret, tk))
			break
		//case '*':
		//	pos := utils.Pos{Line: int64(line), Column: int64(column)}
		//	tk := models.NewToken(pos, models.MUL, 1, run)
		//	err = errors.Join(err, AppendToken(filename, &ret, tk))
		//	break
		//case '<':
		//	pos := utils.Pos{Line: int64(line), Column: int64(column)}
		//	tk := models.NewToken(pos, models.LOWER_THAN, 1, run)
		//	err = errors.Join(err, AppendToken(filename, &ret, tk))
		//	break
		//case '>':
		//	pos := utils.Pos{Line: int64(line), Column: int64(column)}
		//	tk := models.NewToken(pos, models.GREATER_THAN, 1, run)
		//	err = errors.Join(err, AppendToken(filename, &ret, tk))
		//	break
		case ' ':
			pos := utils.Pos{Line: int64(line), Column: int64(column)}
			tk := models.NewToken(pos, models.SPACE, 1, run)
			err = errors.Join(err, AppendToken(filename, &ret, tk))
			break
		case '|':
			pos := utils.Pos{Line: int64(line), Column: int64(column)}
			tk := models.NewToken(pos, models.OR_BIT, 1, run)
			err = errors.Join(err, AppendToken(filename, &ret, tk))
			break
		case '&':
			pos := utils.Pos{Line: int64(line), Column: int64(column)}
			tk := models.NewToken(pos, models.AND_BIT, 1, run)
			err = errors.Join(err, AppendToken(filename, &ret, tk))
			break
		//case '%':
		//	pos := utils.Pos{Line: int64(line), Column: int64(column)}
		//	tk := models.NewToken(pos, models.MOD, 1, run)
		//	err = errors.Join(err, AppendToken(filename, &ret, tk))
		//	break
		default:
			err = utils.GetUnexpectedTokenErr(filename, string(run), utils.Pos{Line: int64(line), Column: int64(column)})
		}
		if err != nil {
			return nil, err
		}
	}

	pos := utils.Pos{Line: int64(line), Column: int64(column)}
	tk := models.NewToken(pos, models.EOF, 1, '0')
	_ = AppendToken(filename, &ret, tk)
	return ret, nil
}

func AppendToken(filename string, tokens *[]models.Token, token models.Token) error {
	if tokens == nil {
		tokens = &[]models.Token{}
	}
	count := len(*tokens)
	if count > 0 {
		last := (*tokens)[count-1]
		if last.Kind == token.Kind && string(last.Value) == string(token.Value) {
			(*tokens)[count-1].Repeat = last.Repeat + 1
			last = (*tokens)[count-1]
		}
		if c, v := models.CombineTokens(last, token); c != models.UNKNOW {
			(*tokens)[count-1].Kind = c
			(*tokens)[count-1].Value = v
			return nil
		}
	}
	*tokens = append(*tokens, token)
	return nil
}
