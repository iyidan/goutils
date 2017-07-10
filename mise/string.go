package mise

import (
	"math/rand"
	"time"
)

const (
	letterBytes   = "1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	letterIdxBits = 6                    // 6 bits to represent a letter index, 62letters=111110b
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits, = 111111b
	letterIdxMax  = 63 / letterIdxBits   // of letter indices fitting in 63 bits
)

var randStringSrc = rand.NewSource(time.Now().UnixNano())

// RandStr random a n byte string, which chars in [0-9a-zA-Z]
// see http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
func RandStr(n int) string {
	b := make([]byte, n)
	// A randStringSrc.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, randStringSrc.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randStringSrc.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// RandStrCustom rand str with given letters
func RandStrCustom(n int, letters string) string {
	b := make([]byte, n)
	// A randStringSrc.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, randStringSrc.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randStringSrc.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letters) {
			b[i] = letters[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
