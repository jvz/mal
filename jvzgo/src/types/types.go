package types

import "fmt"

type MalType interface {
}

func NewTypeError(expectedType string, actual MalType) error {
	return fmt.Errorf("unexpected type; expected %v; actual value: %v", expectedType, actual)
}

func RaiseTypeError(expectedType string, actual MalType) (MalType, error) {
	return nil, NewTypeError(expectedType, actual)
}

type MalList struct {
	Value    []MalType
	StartStr string
	EndStr   string
}

func NewList(values ...MalType) MalList {
	return MalList{Value: values, StartStr: "(", EndStr: ")"}
}

func IsList(val MalType) bool {
	switch list := val.(type) {
	case MalList:
		return list.StartStr == "(" && list.EndStr == ")"
	default:
		return false
	}
}

func NewVec(values ...MalType) MalList {
	return MalList{Value: values, StartStr: "[", EndStr: "]"}
}

func IsVec(val MalType) bool {
	switch list := val.(type) {
	case MalList:
		return list.StartStr == "[" && list.EndStr == "]"
	default:
		return false
	}
}

func GetSlice(val MalType) ([]MalType, error) {
	list, ok := val.(MalList)
	if !ok {
		return nil, fmt.Errorf("provided value is not sliceable: %v", val)
	}
	return list.Value, nil
}

type MalMap struct {
	Value map[MalType]MalType
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
