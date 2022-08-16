package galidator

import (
	"fmt"
	"math"
)

// Determines the precision of a float number for print
func determinePrecision(number float64) string {
	for i := 0; ; i++ {
		ten := math.Pow10(i)
		if math.Floor(ten*number) == ten*number {
			return fmt.Sprint(i)
		}
	}
}
