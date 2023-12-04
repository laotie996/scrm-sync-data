package crypto

import (
	"encoding/base64"
)

type Base64 struct {
	Coder *base64.Encoding
}

func (*Base64) Init() *Base64 {
	table := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	b := &Base64{
		Coder: base64.NewEncoding(table),
	}
	return b
}

func (b *Base64) Encode(src []byte) []byte {
	return []byte(b.Coder.EncodeToString(src))
}

func (b *Base64) Decode(src []byte) ([]byte, error) {
	return b.Coder.DecodeString(string(src))
}
