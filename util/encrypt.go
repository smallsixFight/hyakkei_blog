package util

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5Encrypt(password, salt string) (result string, err error) {
	h := md5.New()
	_, err = h.Write([]byte(password + salt))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
