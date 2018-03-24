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
	return NewList(values)
}

func IsList(val MalType) bool {
	switch list := val.(type) {
	case MalList:
		return list.StartStr == "(" && list.EndStr == ")"
	default:
		return false
	}
}

func GetList(val MalType) (MalList, error) {
	if IsList(val) {
		return val.(MalList), nil
	}
	return NewListOf(), NewTypeError("list", val)
}

func NewVec(value []MalType) MalList {
	return MalList{Value: value, StartStr: "[", EndStr: "]"}
}

func NewVecOf(values ...MalType) MalList {
	return NewVec(values)
}

func IsVec(val MalType) bool {
	switch list := val.(type) {
	case MalList:
		return list.StartStr == "[" && list.EndStr == "]"
	default:
		return false
	}
}

func GetVec(val MalType) (MalList, error) {
	if IsVec(val) {
		return val.(MalList), nil
	}
	return NewVecOf(), NewTypeError("vector", val)
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

func GetMap(val MalType) (MalMap, error) {
	if mm, ok := val.(MalMap); ok {
		return mm, nil
	}
	return MalMap{}, NewTypeError("map", val)
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
	}
	return MalSymbol{}, NewTypeError("symbol", val)
}

type MalString struct {
	Value string
}

func (ms MalString) String() string {
	return ms.Value
}

func GetString(val MalType) (MalString, error) {
	if ms, ok := val.(MalString); ok {
		return ms, nil
	}
	return MalString{}, NewTypeError("string", val)
}

type MalKeyword struct {
	Value string
}

func (mk MalKeyword) String() string {
	return mk.Value
}

func GetKeyword(val MalType) (MalKeyword, error) {
	if mk, ok := val.(MalKeyword); ok {
		return mk, nil
	}
	return MalKeyword{}, NewTypeError("keyword", val)
}

type MalInt struct {
	Value int
}

func (mi MalInt) String() string {
	return strconv.Itoa(mi.Value)
}

func GetInt(val MalType) (MalInt, error) {
	if mi, ok := val.(MalInt); ok {
		return mi, nil
	}
	return MalInt{}, NewTypeError("int", val)
}

type MalBool struct {
	Value bool
}

func (mb MalBool) String() string {
	return strconv.FormatBool(mb.Value)
}

func GetBool(val MalType) (MalBool, error) {
	if mb, ok := val.(MalBool); ok {
		return mb, nil
	}
	return MalBool{}, NewTypeError("bool", val)
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

func IsNil(val MalType) bool {
	_, ok := val.(MalNil)
	return ok
}

func GetFn(val MalType) (func([]MalType) (MalType, error), error) {
	if fn, ok := val.(func([]MalType) (MalType, error)); ok {
		return fn, nil
	} else {
		return nil, NewTypeError("function", val)
	}
}
