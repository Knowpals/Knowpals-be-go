package tool

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateVerifyCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}
