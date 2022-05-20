package utils

import "crypto/sha1"

func SHA1OF(id string) []byte {
	h := sha1.New()
	h.Write([]byte(id))
	return h.Sum(nil)
}
