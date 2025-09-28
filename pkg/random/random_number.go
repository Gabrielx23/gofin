package random

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateRandomNumber(digits int) string {
	rand.Seed(time.Now().UnixNano())

	max := 1
	for i := 0; i < digits; i++ {
		max *= 10
	}

	format := fmt.Sprintf("%%0%dd", digits)
	return fmt.Sprintf(format, rand.Intn(max))
}
