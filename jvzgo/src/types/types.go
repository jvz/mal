package types

type MalType interface {
}

type MalList struct {
	Value []MalType
}

func NewList(values ...MalType) MalList {
	return MalList{Value: values}
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
