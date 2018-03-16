package printer

import (
	"fmt"
	"strconv"
	"types"
)

func PrintStr(obj types.MalType, printReadably bool) string {
	switch o := obj.(type) {
	case types.MalList:
		listSize := len(o.Value)
		if listSize == 0 {
			return "()"
		}
		list := make([]string, 0, listSize)
		for _, val := range o.Value {
			printed := PrintStr(val, printReadably)
			list = append(list, printed)
		}
		size := listSize + 1
		for _, s := range list {
			size += len(s)
		}
		b := make([]byte, size)
		b[0] = '('
		pos := copy(b[1:], list[0]) + 1
		for _, s := range list[1:] {
			b[pos] = ' '
			pos++
			pos += copy(b[pos:], s)
		}
		b[pos] = ')'
		return string(b)
	case types.MalInt:
		return strconv.Itoa(o.Value)
	case types.MalSymbol:
		return o.Value
	case types.MalString:
		if printReadably {
			return strconv.Quote(o.Value)
		} else {
			return o.Value
		}
	case types.MalBool:
		if o.Value {
			return "true"
		} else {
			return "false"
		}
	case types.MalNil:
		return "nil"
	default:
		return fmt.Sprintf("%v", o)
	}
}
