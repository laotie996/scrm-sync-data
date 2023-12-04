//Package crypto_services
/*
 * @Time    : 2021年05月06日 09:28:55
 * @Author  : user
 * @Project : SEGI
 * @File    : rsa_service.go
 * @Software: GoLand
 * @Describe:
 */
package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

type RSA struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

var privateKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAvXxtRMlwKLvS4vq1Ws3wXFxWhATDF+0Mb8FZP27Q12F/vdex
loKpg/xHSbZ8BvTuLAcdCp2hsewWMFxYT8GO8/ocTnww+hmT46bKL+IYjfT08L/Q
CS48lNDEHw/AJQWHuz0eTNnPAdqGyVDoyuTrDG/1p675hl1JFQ2wag4W9TvzesEL
raczm9gV+G5JO5T60VtIiT8qv7Xtija7B32S150OYzcqB9Y1FIOgZwi//JSqVcCZ
bA9dwdR4z/c4ZeavGYRjw3Nxmpr391RnhPJsUaZePVO0+0vjvc4xKo+uTgjnxaYQ
APLcIaW7a5aT9vPu3dVQqH24oDLwapK5o+pMTwIDAQABAoIBAEJjwB0jvupmtILW
eIgyM48Iwz9KM9jEM2FMiyqQdxefj5RCaHRln2MhFxEaoOYHFxPSAjjf9fFS2Itl
L4PyO1X8QcK1/KsEHy7aF2kjfJrwQom/ivJMgulZI/4GFLpj76SIy379qWvq1uLe
OFIuVrRq3dU8lyyerqZzv2XJPf35wsqYpjAY+MHal4wVUcBrIErGxiRgPeFHt8/+
3hIp8lQj7VXv9UGOaJE53uEjYlzEYB49Bn9VL4zGZcBJLtNJ5Ga1PR+s3fzfI9L2
PsLgwpRfr0laVO8eFx8/2blfb0+hnwiq8atJN2EEToYIO1b9UySaipnGkN9ijqA1
wxRZAAECgYEA6I+8ojE4M58WWwEkVa637XSRwc2f59pPQfD8ZGGyyDXhj2IHyuNf
2KOCrBNuEpApL+i9T+SA+7oiQAOv3gskYPNzfNXizuKKnptk7arego2lTS8MmQ1N
+EGoQ4kL4T66+sYA8pbHJYKpop2t3b6X8kRdHkM7VdNhsSkwPLzzzE8CgYEA0JVP
osKz+sKtsP3KRspfgheQMEtN5OTv8fN1icPISef3o41Vgn1d902XD7gejTlT2uaz
dAZjmyhAXShsu9tuW5wdNlec7OgJNL7CFLgfq4raKTFPzVBnDQcvkoJDST8C14uE
zUB8iP9Qef8LdWjOCbvPvxzj5z+/AefmUeOBgAECgYAzvrekBvuQUNdyjEB/aA64
oPVFf/ESb3qvG2WtGCflcEet1YwSUqEi3M7YZsJZEwZ3CHwG6LriR7oTMD7TUvtX
LxQFNLnwemuJet3sG9bCia28DLnq3BD7KfC8hiPEjGaiIahYUcAh0n7YFel3q063
swpdD2yGXjAgcR+whAJi9QKBgFUE1z33cNcAbT3cZIuwR3fGVP5qY2iuLGOJodoy
lDwTsivXGwiiQx/fX3iHyuQzfsuWB4V+aTFAGeQe6xszWOs8WefhlFQ4TDBUpNU/
K6GHal8p+2YrpKV5zVlBgV2ksmrbMpla9Kf+sXXYkHod6wWqqQY0J20F3bxQIuA0
BIABAoGAefizN7wiMojxMi0UQBQHeQEngnC/w5b2qcGCG2RL6zijxSxQHSUTQuyB
XrxIVZMgoJLMv+M6WXYzzdTBrLm2NP830OEwvjqvxS79dZ48jcEd8PTIplQI3bxE
kQNCQDwuF9yxX7pYo5ruMJaoWkLM5M6CiXYuOqBH2qAm9HwVB2w=
-----END RSA PRIVATE KEY-----`)

var publicKey = []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvXxtRMlwKLvS4vq1Ws3w
XFxWhATDF+0Mb8FZP27Q12F/vdexloKpg/xHSbZ8BvTuLAcdCp2hsewWMFxYT8GO
8/ocTnww+hmT46bKL+IYjfT08L/QCS48lNDEHw/AJQWHuz0eTNnPAdqGyVDoyuTr
DG/1p675hl1JFQ2wag4W9TvzesELraczm9gV+G5JO5T60VtIiT8qv7Xtija7B32S
150OYzcqB9Y1FIOgZwi//JSqVcCZbA9dwdR4z/c4ZeavGYRjw3Nxmpr391RnhPJs
UaZePVO0+0vjvc4xKo+uTgjnxaYQAPLcIaW7a5aT9vPu3dVQqH24oDLwapK5o+pM
TwIDAQAB
-----END PUBLIC KEY-----`)

func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf)
	}
	return chunks
}

func (r *RSA) Init() (*RSA, error) {
	return r.Load(privateKey, publicKey)
}

func (r *RSA) Load(privateKey []byte, publicKey []byte) (*RSA, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	block, _ = pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error")
	}
	pri, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	} else {
		return &RSA{
			publicKey:  pub,
			privateKey: pri,
		}, nil
	}
}

// Encrypt rsa 加密
func (r *RSA) Encrypt(origin []byte) ([]byte, error) {
	var buffer bytes.Buffer
	partLen := r.publicKey.N.BitLen()/8 - 11
	chunks := split(origin, partLen)
	for _, chunk := range chunks {
		encryptBytes, err := rsa.EncryptPKCS1v15(rand.Reader, r.publicKey, chunk)
		if err != nil {
			return nil, err
		}
		buffer.Write(encryptBytes)
	}
	return buffer.Bytes(), nil
}

// Decrypt rsa解密
func (r *RSA) Decrypt(cipher []byte) ([]byte, error) {
	var buffer bytes.Buffer
	partLen := r.publicKey.N.BitLen() / 8
	chunks := split(cipher, partLen)
	for _, chunk := range chunks {
		decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, r.privateKey, chunk)
		if err != nil {
			return nil, err
		}
		buffer.Write(decrypted)
	}
	return buffer.Bytes(), nil
}
