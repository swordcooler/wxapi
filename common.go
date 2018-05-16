package wx

import (
	"crypto/md5"
	"fmt"
	"io"
	"math/rand"
	"sort"
	"time"
)

const letterBytes = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func RandomString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
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

func GenerateSign(secret string, params map[string]string) string {
	paramList := make([]string, 0)
	for k, _ := range params {
		paramList = append(paramList, k)
	}

	sort.Strings(paramList)

	var paramStr string
	for k, v := range paramList {
		paramStr += fmt.Sprint("%s=%s", v, params[v])
		if k != len(paramList) {
			paramStr += "&"
		}
	}

	stringSignTemp := fmt.Sprintf("%s&key=%s", paramStr, secret)
	md5w := md5.New()
	io.WriteString(md5w, stringSignTemp)
	sign := fmt.Sprintf("%x", md5w.Sum(nil))

	return sign
}
