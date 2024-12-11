package random

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	util "hieupc05.github/backend-server/internal/utils"
)

const alphabet = "abcdefghiklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())

}

// RandomInt generate a radom integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generate a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner generate a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney generates a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency generate a random currency code
func RandomCurrency() string {
	currencies := []string{util.EUR, util.USD, util.CAD}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

// radom email
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}
