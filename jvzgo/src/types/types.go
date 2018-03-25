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
	Value    []MalType
	startStr string
	endStr   string
}

func (ml MalList) String() string {
	vals := make([]string, len(ml.Value))
	for i, val := range ml.Value {
		vals[i] = fmt.Sprint(val)
	}
	return ml.startStr + strings.Join(vals, " ") + ml.endStr
}

func (ml MalList) Surround(str string) string {
	return ml.startStr + str + ml.endStr
}

func (ml MalList) New(value []MalType) MalList {
	return MalList{Value: value, startStr: ml.startStr, endStr: ml.endStr}
}

func NewList(value []MalType) MalList {
	return MalList{Value: value, startStr: "(", endStr: ")"}
}

func NewListOf(values ...MalType) MalList {
	return NewList(values)
}

func IsList(val MalType) bool {
	switch list := val.(type) {
	case MalList:
		return list.startStr == "(" && list.endStr == ")"
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
	return MalList{Value: value, startStr: "[", endStr: "]"}
}

func NewVecOf(values ...MalType) MalList {
	return NewVec(values)
}

func IsVec(val MalType) bool {
	switch list := val.(type) {
	case MalList:
		return list.startStr == "[" && list.endStr == "]"
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

type MalFn struct {
	fn   func([]MalType) (MalType, error)
	meta MalType
}

func (MalFn) String() string {
	return "#<function>"
}

type MalFunc struct {
	eval    func(MalType, EnvType) (MalType, error)
	binds   []MalType
	expr    MalType
	env     EnvType
	meta    MalType
	isMacro bool
}

func NewFunc(eval func(MalType, EnvType) (MalType, error), binds []MalType, expr MalType, env EnvType) MalFunc {
	return MalFunc{eval: eval, binds: binds, expr: expr, env: env, meta: MalNil{}}
}

func (mf MalFunc) String() string {
	if mf.isMacro {
		return "#<macro>"
	}
	return "#<function>"
}

func (mf MalFunc) Fn() func([]MalType) (MalType, error) {
	return func(args []MalType) (MalType, error) {
		inner, err := mf.env.New(mf.binds, args)
		if err != nil {
			return nil, err
		}
		return mf.eval(mf.expr, inner)
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
	case MalFn:
		return fn.fn, nil
	case MalFunc:
		return fn.Fn(), nil
	case func([]MalType) (MalType, error):
		return fn, nil
	default:
		return nil, NewTypeError("function", val)
	}
}

func IsFn(val MalType) bool {
	switch val.(type) {
	case MalFn:
		return true
	case MalFunc:
		return true
	case func([]MalType) (MalType, error):
		return true
	default:
		return false
	}
}

func IsMacro(val MalType) bool {
	if fn, ok := val.(MalFunc); ok {
		return fn.isMacro
	}
	return false
}

func GetMeta(val MalType) (MalType, error) {
	switch fn := val.(type) {
	case MalFn:
		return fn.meta, nil
	case MalFunc:
		return fn.meta, nil
	case func([]MalType) (MalType, error):
		return MalNil{}, nil
	default:
		return RaiseTypeError("function", val)
	}
}

func WithMeta(val MalType, meta MalType) (MalType, error) {
	switch fn := val.(type) {
	case MalFn:
		return MalFn{fn: fn.fn, meta: meta}, nil
	case MalFunc:
		return MalFunc{eval: fn.eval, binds: fn.binds, expr: fn.expr, env: fn.env, meta: meta, isMacro: fn.isMacro}, nil
	case func([]MalType) (MalType, error):
		return MalFn{fn: fn, meta: meta}, nil
	default:
		return RaiseTypeError("function", val)
	}
}
