package random

import (
	"math/rand"
	"time"
)

type random struct {
}

func NewRandom() *random {
	return &random{}
}

type Random interface {
	RandomString() (int, error)
}

func (r *random) RandomString() string {
	randomizer := rand.New(rand.NewSource(time.Now().Unix()))
	letters := []rune("qwertyuioplkjhgfdsazxcvbnmQWERTYUIOPLKJHGFDSAZXCVBNM1234567890")

	b := make([]rune, 10)
	for i := range b {
		b[i] = letters[randomizer.Intn(len(letters))]
	}

	return string(b)
}
