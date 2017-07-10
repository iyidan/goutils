package mise

import (
	"crypto/md5"
	"encoding/hex"
)

// Md5String return 32 lower-case-letter hash
func Md5String(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
