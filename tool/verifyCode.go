package tool

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func GenerateRandomCode(n int) string {
	rand.Seed(time.Now().UnixNano())
	max := int(math.Pow10(n)) - 1
	min := int(math.Pow10(n - 1))
	code := rand.Intn(max-min+1) + min
	return fmt.Sprintf("%0*d", n, code)
}
