package util

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
)

func RandomString(n int) string {
	replacer := strings.NewReplacer("/", "", "+", "")
	str := ""

	for len(str) < n {
		b := make([]byte, n-len(str))
		_, err := rand.Read(b)
		if err != nil {
			panic(err)
		}
		str = str + base64.RawStdEncoding.EncodeToString(b)
		str = replacer.Replace(str)
	}

	return str[:n]
}
