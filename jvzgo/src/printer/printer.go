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
		return joinStrings(strs, o.StartStr, o.EndStr)
	case MalMap:
		strs := make([]string, 0, len(o.Value)*2)
		for k, v := range o.Value {
			key := PrintStr(k, printReadably)
			val := PrintStr(v, printReadably)
			strs = append(strs, key, val)
		}
		return joinStrings(strs, "{", "}")
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
	case MalBool:
		if o.Value {
			return "true"
		} else {
			return "false"
		}
	case MalNil:
		return "nil"
	default:
		return fmt.Sprintf("%v", o)
	}
}

func joinStrings(strs []string, start, end string) string {
	if len(strs) == 0 {
		return start + end
	} else {
		return start + strings.Join(strs, " ") + end
	}
}
