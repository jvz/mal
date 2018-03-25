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

type MalError struct {
	Value MalType
}

func (e MalError) Error() string {
	return fmt.Sprint(e.Value)
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

func IsMap(val MalType) bool {
	_, ok := val.(MalMap)
	return ok
}

func CopyMap(val MalMap) MalMap {
	m := make(map[MalType]MalType)
	for k, v := range val.Value {
		m[k] = v
	}
	return MalMap{Value: m}
}

type MalAtom struct {
	Value MalType
}

func (ma *MalAtom) String() string {
	return fmt.Sprint(ma.Value)
}

func (ma *MalAtom) Set(val MalType) {
	ma.Value = val
}

func NewAtom(val MalType) *MalAtom {
	return &MalAtom{Value: val}
}

func GetAtom(val MalType) (*MalAtom, error) {
	if ma, ok := val.(*MalAtom); ok {
		return ma, nil
	}
	return nil, NewTypeError("atom", val)
}

func IsAtom(val MalType) bool {
	_, ok := val.(*MalAtom)
	return ok
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

func IsSymbol(val MalType) bool {
	_, ok := val.(MalSymbol)
	return ok
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

func IsString(val MalType) bool {
	_, ok := val.(MalString)
	return ok
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

func IsKeyword(val MalType) bool {
	_, ok := val.(MalKeyword)
	return ok
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

func IsInt(val MalType) bool {
	_, ok := val.(MalInt)
	return ok
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

func IsBool(val MalType) bool {
	_, ok := val.(MalBool)
	return ok
}

var MalTrue = MalBool{Value: true}

func IsTrue(val MalType) bool {
	b, ok := val.(MalBool)
	return ok && b.Value
}

var MalFalse = MalBool{Value: false}

func IsFalse(val MalType) bool {
	b, ok := val.(MalBool)
	return ok && !b.Value
}

type MalNil struct {
	Value interface{}
}

func (MalNil) String() string {
	return "nil"
}

func IsNil(val MalType) bool {
	_, ok := val.(MalNil)
	return ok
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

type MalFunc struct {
	Eval    func(MalType, EnvType) (MalType, error)
	Binds   []MalType
	Expr    MalType
	Env     EnvType
	isMacro bool
}

func (mf MalFunc) Fn() func([]MalType) (MalType, error) {
	return func(args []MalType) (MalType, error) {
		inner, err := mf.Env.New(mf.Binds, args)
		if err != nil {
			return nil, err
		}
		return mf.Eval(mf.Expr, inner)
	}
}

func (mf *MalFunc) IsMacro() bool {
	return mf.isMacro
}

func (mf *MalFunc) SetMacro(b bool) {
	mf.isMacro = b
}

func GetFn(val MalType) (func([]MalType) (MalType, error), error) {
	switch fn := val.(type) {
	case MalFunc:
		return fn.Fn(), nil
	case func([]MalType) (MalType, error):
		return fn, nil
	default:
		return nil, NewTypeError("function", val)
	}
}
