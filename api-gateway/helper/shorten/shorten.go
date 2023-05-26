package shorten

import (
	"hash/fnv"
	"math/big"
	"strings"
)

const (
	base62Alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

type shorten struct{}

func New() Shorten {
	return new(shorten)
}

type Shorten interface {
	Encode(str string) string
	Decode(str string) *big.Int
	EnforceHTTP(url string) string
}

func (s *shorten) generateUniqueID(str string) *big.Int {
	shorten := fnv.New64a()
	shorten.Write([]byte(str))
	sum := shorten.Sum(nil)

	identifier := new(big.Int)
	identifier.SetBytes(sum)

	return identifier
}

func (s *shorten) Encode(str string) string {
	number := s.generateUniqueID(str)

	base := big.NewInt(62)
	result := ""

	zero := big.NewInt(0)
	rem := &big.Int{}

	for number.Cmp(zero) > 0 {
		number.DivMod(number, base, rem)
		result = string(base62Alphabet[rem.Int64()]) + result
	}

	padLen := 6 - len(result)
	if padLen > 0 {
		result = strings.Repeat("0", padLen) + result
	}

	return result
}

func (s *shorten) Decode(str string) *big.Int {
	base := big.NewInt(62)
	identifier := big.NewInt(0)

	for _, char := range str {
		value := strings.IndexRune(base62Alphabet, char)
		identifier.Mul(identifier, base)
		identifier.Add(identifier, big.NewInt(int64(value)))
	}

	return identifier
}

func (s *shorten) EnforceHTTP(url string) string {
	if url[:4] != "http" {
		return "http://" + url
	}
	return url
}
