package types

import (
	"fmt"
	"strconv"
	"strings"
)

type MalType interface {
}

type EnvType interface {
	Set(key string, val MalType)
	Find(key string) EnvType
	Get(key string) (MalType, error)
	New(binds, exprs []MalType) (EnvType, error)
}

func NewTypeError(expectedType string, actual MalType) error {
	return fmt.Errorf("unexpected type; expected %v; actual value: %v", expectedType, actual)
}

func RaiseTypeError(expectedType string, actual MalType) (MalType, error) {
	return nil, NewTypeError(expectedType, actual)
}

type MalList struct {
	Value []MalType
	// TODO: make these unexported
	StartStr string
	EndStr   string
}

func (ml MalList) String() string {
	vals := make([]string, len(ml.Value))
	for i, val := range ml.Value {
		vals[i] = fmt.Sprint(val)
	}
	return ml.StartStr + strings.Join(vals, " ") + ml.EndStr
}

func NewList(value []MalType) MalList {
	return MalList{Value: value, StartStr: "(", EndStr: ")"}
}

func NewListOf(values ...MalType) MalList {
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

func (ms MalSymbol) String() string {
	return ms.Value
}

func GetSymbol(val MalType) (MalSymbol, error) {
	if ms, ok := val.(MalSymbol); ok {
		return ms, nil
	} else {
		return MalSymbol{}, NewTypeError("symbol", val)
	}
}

type MalString struct {
	Value string
}

func (ms MalString) String() string {
	return ms.Value
}

type MalKeyword struct {
	Value string
}

func (mk MalKeyword) String() string {
	return mk.Value
}

type MalInt struct {
	Value int
}

func (mi MalInt) String() string {
	return strconv.Itoa(mi.Value)
}

type MalBool struct {
	Value bool
}

func (mb MalBool) String() string {
	return strconv.FormatBool(mb.Value)
}

var MalTrue = MalBool{Value: true}

var MalFalse = MalBool{Value: false}

type MalNil struct {
	Value interface{}
}

func (MalNil) String() string {
	return "nil"
}

func IsTruthy(val MalType) bool {
	switch val := val.(type) {
	case MalBool:
		return val.Value
	case MalNil:
		return false
	default:
		return true
	}
}

func GetFn(val MalType) (func([]MalType) (MalType, error), error) {
	if fn, ok := val.(func([]MalType) (MalType, error)); ok {
		return fn, nil
	} else {
		return nil, NewTypeError("function", val)
	}
}
