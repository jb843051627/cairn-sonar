package codec

import (
	"crypto/sha256"
	"encoding/hex"
)

func Digest(parts ...string) string {
	h := sha256.New()
	for _, part := range parts {
		h.Write([]byte(part))
	}
	return hex.EncodeToString(h.Sum(nil))
}
