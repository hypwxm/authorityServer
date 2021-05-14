package service2

import (
	"crypto/md5"
	"encoding/hex"
)

func SignPwd(pwd string, salt string) string {
	m := md5.New()
	m.Write([]byte(pwd))
	m.Write([]byte(salt))
	return hex.EncodeToString(m.Sum(nil))
}
