package printer

import (
	"fmt"
	"strconv"
	"strings"
	. "types"
)

func PrintStr(obj MalType, printReadably bool) string {
	switch o := obj.(type) {
	case MalList:
		strs := make([]string, len(o.Value))
		for i, val := range o.Value {
			strs[i] = PrintStr(val, printReadably)
		}
		return o.Surround(strings.Join(strs, " "))
	case MalMap:
		strs := make([]string, 0, len(o.Value)*2)
		for k, v := range o.Value {
			key := PrintStr(k, printReadably)
			val := PrintStr(v, printReadably)
			strs = append(strs, key, val)
		}
		return joinStrings(strs, "{", "}")
	case *MalAtom:
		return "(atom " + PrintStr(o.Value, printReadably) + ")"
	case MalInt:
		return strconv.Itoa(o.Value)
	case MalSymbol:
		return o.Value
	case MalString:
		if printReadably {
			return strconv.Quote(o.Value)
		} else {
			return o.Value
		}
	case MalKeyword:
		return ":" + o.Value
	case MalBool:
		if o.Value {
			return "true"
		} else {
			return "false"
		}
	case MalNil:
		return "nil"
	case MalFn:
		return "#<function>"
	case MalFunc:
		if o.IsMacro() {
			return "#<macro>"
		}
		return "#<function>"
	case func([]MalType) (MalType, error):
		return "#<function>"
	case MalError:
		return PrintStr(o.Value, printReadably)
	default:
		return fmt.Sprintf("#<unknown: %v>", o)
	}
}

func joinStrings(strs []string, start, end string) string {
	if len(strs) == 0 {
		return start + end
	} else {
		return start + strings.Join(strs, " ") + end
	}
}
