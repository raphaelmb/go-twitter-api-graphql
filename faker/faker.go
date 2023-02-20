package faker

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

var Password = "$2a$04$.vcyIl5L6cbcB0H8fSxWQO5HQvdsqx4uW/5uY.Bp32fGTj52n3rQy"

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandStringRunesLower(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes)/2)]
	}
	return string(b)
}

func RandInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func Username() string {
	return RandStringRunes(RandInt(2, 10))
}

func ID() string {
	return fmt.Sprintf("%s-%s-%s-%s", RandStringRunes(4), RandStringRunes(4), RandStringRunes(4), RandStringRunes(4))
}

func Email() string {
	return fmt.Sprintf("%s@example.com", RandStringRunesLower(RandInt(5, 10)))
}
