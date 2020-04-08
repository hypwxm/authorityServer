package util

import (
	"crypto/md5"
	"encoding/hex"
	uuid "github.com/satori/go.uuid"
	"strings"
	"time"
)

func FormatQuota(s string) string {
	s = strings.Replace(s, "'", "\\'", -1)
	s = strings.Replace(s, "(", "-", -1)
	s = strings.Replace(s, ")", "-", -1)
	return s
}

func GetCurrentMS() int64 {
	return time.Now().UnixNano() / 1e6
}

func GetUuid() string {
	return uuid.NewV4().String()
}


func SignPwd(pwd string, salt string) string {
	m := md5.New()
	m.Write([]byte(pwd))
	m.Write([]byte(salt))
	return hex.EncodeToString(m.Sum(nil))
}
