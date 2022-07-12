package sign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

type Signer struct {
	Key string
}

func NewSigner(key string) *Signer {
	return &Signer{
		Key: key,
	}
}

func (e *Signer) Sign(src string) []byte {
	h := hmac.New(sha256.New, []byte(e.Key))
	h.Write([]byte(src))
	return h.Sum(nil)
}

func (e *Signer) Verify(sign []byte, value string) bool {
	return hmac.Equal(e.Sign(value), sign)
}

func (e *Signer) EncodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func (e *Signer) DecodeBase64(sEnc string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(sEnc)
}
