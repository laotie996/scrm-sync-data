package crypto

import (
	"crypto/sha1"
	"encoding/hex"
)

type SHA1 struct {
}

func (*SHA1) Crypt(v1 []byte, v2 []byte) string {
	s := sha1.New()
	s.Write(v1)
	return hex.EncodeToString(s.Sum(v2))
}
