package utils

import "fmt"

type (
	ErrCode uint

	ErrLabel string

	Err struct {
		Code  ErrCode
		Label ErrLabel
		Msg   string
	}
)

func (this Err) Error() string {
	return fmt.Sprintf("error %.4d | %s\n", this.Code, this.Msg)
}

const (
	NilPointerExceptionErrCode ErrCode = iota
	FileNotFoundErrCode
	EmptyFileErrCode
	DivideByZeroErrCode
	InvalidTokenPerSubsetErrCode
	InvalidArgumentErrCode
	UnexpectedTokenErrCode
	UnexpectedForLoopStatementInRootErrCode
	UnexpectedIfStatementInRootErrCode
	TooManyValuesErrCode
	ConsecutiveOperatorsErrCode
	ConsecutiveValuesErrCode
	MismatchedTypesErrCode
	UnkownVariableErrCode

	NilPointerExceptionErrLabel              ErrLabel = "error.nil.pointer"
	FileNotFoundErrLabel                     ErrLabel = "error.file.not.found"
	EmptyFileErrLabel                        ErrLabel = "error.empty.file"
	DivideByZeroErrLabel                     ErrLabel = "error.divide.by.zero"
	InvalidTokenPerSubsetErrLabel            ErrLabel = "error.invalid.token.per.subset"
	InvalidArgumentErrLabel                  ErrLabel = "error.invalid.argument"
	UnexpectedTokenErrLabel                  ErrLabel = "error.unexpected.token"
	UnexpectedForLoopStatementInRootErrLabel ErrLabel = "error.unexpected.for.in.root"
	UnexpectedIfStatementInRootErrLabel      ErrLabel = "error.unexpected.if.in.root"
	TooManyValuesErrLabel                    ErrLabel = "error.too.many.values"
	ConsecutiveOperatorsErrLabel             ErrLabel = "error.consecutive.operators"
	ConsecutiveValuesErrLabel                ErrLabel = "error.consecutive.values"
	MismatchedTypesErrLabel                  ErrLabel = "error.mismatch.type.values"
	UnkownVariableErrLabel                   ErrLabel = "error.mismatch.variable"
)

func GetNilPointerExceptionErr() Err {
	return Err{
		Code:  NilPointerExceptionErrCode,
		Label: NilPointerExceptionErrLabel,
		Msg:   "",
	}
}

func GetDivideByZeroErr() Err {
	return Err{
		Code:  DivideByZeroErrCode,
		Label: DivideByZeroErrLabel,
		Msg:   "",
	}
}

func GetFileNotFoundErr() Err {
	return Err{
		Code:  FileNotFoundErrCode,
		Label: FileNotFoundErrLabel,
		Msg:   "",
	}
}

func GetInvalidArgumentErr() Err {
	return Err{
		Code:  InvalidArgumentErrCode,
		Label: InvalidArgumentErrLabel,
		Msg:   "",
	}
}

func GetEmptyFileErr(filename string) Err {
	return Err{
		Code:  EmptyFileErrCode,
		Label: EmptyFileErrLabel,
		Msg:   fmt.Sprintf("The file '%s' is empty.", filename),
	}
}

func GetUnexpectedTokenNoPosErr(filename string, word string) Err {
	return Err{
		Code:  UnexpectedTokenErrCode,
		Label: UnexpectedTokenErrLabel,
		Msg:   fmt.Sprintf("Unexpected token '%s' in the file '%s'.", word, filename),
	}
}

func GetUnexpectedForLoopStatementInRoot(filename string, pos Pos) Err {
	return Err{
		Code:  UnexpectedForLoopStatementInRootErrCode,
		Label: UnexpectedForLoopStatementInRootErrLabel,
		Msg:   fmt.Sprintf("Unexpected for token in the root. File '%s' at line %d, column %d.", filename, pos.Line, pos.Column),
	}
}

func GetUnexpectedIfStatementInRoot(filename string, pos Pos) Err {
	return Err{
		Code:  UnexpectedIfStatementInRootErrCode,
		Label: UnexpectedIfStatementInRootErrLabel,
		Msg:   fmt.Sprintf("Unexpected if token in the root. File '%s' at line %d, column %d.", filename, pos.Line, pos.Column),
	}
}

func GetInvalidTokenPerSubset(filename string, word string, pos Pos) Err {
	return Err{
		Code:  InvalidTokenPerSubsetErrCode,
		Label: InvalidTokenPerSubsetErrLabel,
		Msg:   fmt.Sprintf("The specified set do not allows the token '%s', take it off from '%s' at line %d, column %d.", word, filename, pos.Line, pos.Column),
	}
}

func GetUnexpectedTokenErr(filename string, word string, pos Pos) Err {
	return Err{
		Code:  UnexpectedTokenErrCode,
		Label: UnexpectedTokenErrLabel,
		Msg:   fmt.Sprintf("Unexpected token '%s' in the file '%s' at line %d, column %d.", word, filename, pos.Line, pos.Column),
	}
}

func GetUnkownVariableErr(filename string, varName string, pos Pos) Err {
	return Err{
		Code:  UnkownVariableErrCode,
		Label: UnkownVariableErrLabel,
		Msg:   fmt.Sprintf("Unexpected variable '%s' in the file '%s' at line %d, column %d.", varName, filename, pos.Line, pos.Column),
	}
}

func GetExpectedSomeTokenErr(filename string, pos Pos) Err {
	return Err{
		Code:  UnexpectedTokenErrCode,
		Label: UnexpectedTokenErrLabel,
		Msg:   fmt.Sprintf("Missing some token in the file '%s' at line %d, column %d.", filename, pos.Line, pos.Column),
	}
}

func GetExpectedTokenErr(filename string, phrase string, pos Pos) Err {
	return Err{
		Code:  UnexpectedTokenErrCode,
		Label: UnexpectedTokenErrLabel,
		Msg:   fmt.Sprintf("Missing %s in the file '%s' at line %d, column %d", phrase, filename, pos.Line, pos.Column),
	}
}

func GetExpectedTokenErrOr(filename string, phrase, add string, pos Pos) Err {
	return Err{
		Code:  UnexpectedTokenErrCode,
		Label: UnexpectedTokenErrLabel,
		Msg:   fmt.Sprintf("Missing %s in the file '%s' at line %d, column %d, or %s", phrase, filename, pos.Line, pos.Column, add),
	}
}

func GetTooManyValuesErr(filename string, line int64) Err {
	return Err{
		Code:  TooManyValuesErrCode,
		Label: TooManyValuesErrLabel,
		Msg:   fmt.Sprintf("Too many values associated in the file '%s' at line %d", filename, line),
	}
}

func GetConsecutiveOperatorsErr(filename string, pos Pos) Err {
	return Err{
		Code:  ConsecutiveOperatorsErrCode,
		Label: ConsecutiveOperatorsErrLabel,
		Msg:   fmt.Sprintf("Found consecutive operators in the file '%s' at line %d", filename, pos.Line),
	}
}

func GetConsecutiveValuesErr(filename string, pos Pos) Err {
	return Err{
		Code:  ConsecutiveValuesErrCode,
		Label: ConsecutiveValuesErrLabel,
		Msg:   fmt.Sprintf("Found consecutive values in the file '%s' at line %d", filename, pos.Line),
	}
}

func GetMismatchedTypesErr(filename string, expected, found string, pos Pos) Err {
	return Err{
		Code:  MismatchedTypesErrCode,
		Label: MismatchedTypesErrLabel,
		Msg:   fmt.Sprintf("Mismatched types: expected %s but found %s in file '%s' at line %d, column %d", expected, found, filename, pos.Line, pos.Column),
	}
}
