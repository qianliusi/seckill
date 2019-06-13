package utils

import (
	"crypto/md5"
	"fmt"
)

func Md5(data string) string {
	md5Str := md5.Sum([]byte(data))
	return encodeHex(md5Str[:])
}

func encodeHex(data []byte) string {
	var buf []byte
	for i := 0; i < len(data); i++ {
		if (data[i] & 0xff) < 0x10 {
			buf = append(buf, '0')
		}
		buf = append(buf, fmt.Sprintf("%x", data[i]&0xff)...)
	}
	return string(buf)
}
