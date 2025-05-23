package parser

import (
	"errors"
	"exps-heredia/models"
	"exps-heredia/utils"
)

func (parser *Parser) ParseExpression(sub bool) (IExp, error) {
	token := parser.Get(0)
	if token == nil {
		return nil, utils.GetUnexpectedTokenNoPosErr(parser.Filename, "EOF")
	}

	var ret *ExpChain = nil
	var lastWasSig bool

	for i := 0; token != nil; i++ {

		if token.Kind == models.SPACE || token.Kind == models.TAB {
			parser.Consume(1)
			token = parser.Get(0)
			continue
		} else if token.Kind == models.R_PAREN {
			if sub == false {
				return nil, errors.New("unbalanced parenthesis")
			}
			break
		} else if token.Kind == models.L_PAREN {
			if ret != nil && !lastWasSig {
				return nil, errors.New("consecutive numbers")
			}
			lastWasSig = false

			// Start first value
			if ret == nil {
				ret = NewExpP(nil, nil, models.UNKNOW, parser.At(), parser)
			}

			parser.Consume(1)
			e, err := parser.ParseExpression(true)
			if err != nil {
				return nil, err
			}

			ret.AddTerm(e)

		} else if token.Kind == models.NUMBER || token.Kind == models.ID {
			if ret != nil && !lastWasSig {
				return nil, errors.New("consecutive numbers")
			}
			lastWasSig = false

			// Start first value
			if ret == nil {
				ret = NewExpP(nil, nil, models.UNKNOW, parser.At(), parser)
			}

			// Add the number to the list
			if token.Kind == models.ID {
				variable := parser.Variables[string(token.Value)]
				if variable == nil {
					return nil, errors.New("unknown id '" + string(token.Value) + "'")
				}
				ret.AddTerm(NewIdExp(variable.Id + parser.GetVarsOffset()))
			} else {
				ret.AddTerm(NewVExp(uint16(token.Value[0])))
			}

		} else if token.IsSignal() {
			if lastWasSig {
				return nil, errors.New("consecutive signals")
			}
			lastWasSig = true

			// If ret is nil, initialize it
			if ret == nil {
				ret = NewExpP(nil, nil, token.Kind, parser.At(), parser)
				parser.Consume(1)
				token = parser.Get(0)
				continue
			}

			// Chain without a signal
			if ret.Signal == models.UNKNOW {
				ret.Signal = token.Kind
				parser.Consume(1)
				token = parser.Get(0)
				continue
			}

			// Same signal, just go ahead
			if ret.Signal == token.Kind {
				parser.Consume(1)
				token = parser.Get(0)
				continue
			}

			// Actual signal has precedence compared to the last one
			if token.Kind.Weight() > ret.Signal.Weight() {
				ret, _ = ret.DeriveInclusiveExp(token.Kind)

				// Both actual and last signals have the same weight but they're different
			} else if token.Kind.Weight() == ret.Signal.Weight() {
				if ret.Father == nil {
					ret.Father = NewExpP([]IExp{ret}, nil, token.Kind, parser.At(), parser)
				}
				ret = ret.Father

				// Last signal has priority compared to actual one
			} else {

				// Goes up to reduce prescedence
				if ret.Father == nil {
					ret.Father = NewExpP([]IExp{ret}, nil, token.Kind, parser.At(), parser)
				}
				ret = ret.Father

				// If signals are different, generate a new father again
				if ret.Signal != token.Kind {
					ret.Father = NewExpP([]IExp{ret}, nil, token.Kind, parser.At(), parser)
					ret = ret.Father
				}
			}

		} else {

			if sub {
				return nil, errors.New("unbalanced parenthesis")
			}

			// If ret is nil, return 0
			if ret == nil {
				return NewVExp(0), nil
			}

			// If ret has a unique value, return just it
			res := ret.RootFather()
			if len(res.All) == 1 {
				return res.All[0], nil
			}
			return res, nil
		}

		parser.Consume(1)
		token = parser.Get(0)
	}

	// If ret is nil, return 0
	if ret == nil {
		return NewVExp(0), nil
	}

	// If ret has a unique value, return just it
	res := ret.RootFather()
	if len(res.All) == 1 {
		return res.All[0], nil
	}
	return res, nil
}
