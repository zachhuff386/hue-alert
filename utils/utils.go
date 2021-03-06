package utils

import (
	"crypto/rand"
	"math/big"
	mathrand "math/rand"
)

var (
	chars = []rune(
		"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

func RandStr(n int) (str string) {
	strList := make([]rune, n)
	for i := range strList {
		strList[i] = chars[mathrand.Intn(len(chars))]
	}
	str = string(strList)
	return
}

func isEscapable(b byte) bool {
	return !('A' <= b && b <= 'Z' || 'a' <= b && b <= 'z' ||
		'0' <= b && b <= '9' || b == '-' || b == '.' || b == '_' || b == '~')
}

func Escape(s string) string {
	t := make([]byte, 0, 3*len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if isEscapable(c) {
			t = append(t, '%')
			t = append(t, "0123456789ABCDEF"[c>>4])
			t = append(t, "0123456789ABCDEF"[c&15])
		} else {
			t = append(t, s[i])
		}
	}
	return string(t)
}

func init() {
	n, err := rand.Int(rand.Reader, big.NewInt(9223372036854775806))
	if err != nil {
		panic(err)
	}

	mathrand.Seed(n.Int64())
}
