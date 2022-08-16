package galidator

import (
	"fmt"
	"math"
	"reflect"
)

func determinePrecision(number float64) string {
	for i := 0; ; i++ {
		ten := math.Pow10(i)
		if math.Floor(ten*number) == ten*number {
			return fmt.Sprint(i)
		}
	}
}

func isNil(input interface{}) bool {
	inputValue := reflect.ValueOf(input)
	switch inputValue.Kind() {
	case reflect.Array, reflect.Struct, reflect.Map:
		return inputValue.IsNil()
	default:
		return false
	}
}
