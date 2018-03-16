package types

type MalType interface {
}

type MalError struct {
	Value MalType
}

type MalList struct {
	Value []MalType
}

type MalSymbol struct {
	Value string
}

type MalString struct {
	Value string
}

type MalInt struct {
	Value int
}

type MalBool struct {
	Value bool
}

type MalNil struct {
}

var MalTrueVal MalType = MalBool{Value: true}
var MalFalseVal MalType = MalBool{Value: false}
var MalNilVal MalType = MalNil{}
