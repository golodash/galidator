package galidator

import (
	"reflect"
	"runtime"
	"strings"
)

func getFunctionName(temp interface{}) string {
	strs := strings.Split((runtime.FuncForPC(reflect.ValueOf(temp).Pointer()).Name()), ".")
	return strs[len(strs)-1]
}
