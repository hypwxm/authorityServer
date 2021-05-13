package util

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
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
	return strings.ReplaceAll(uuid.NewV4().String(), "-", "")
}

func SignPwd(pwd string, salt string) string {
	m := md5.New()
	m.Write([]byte(pwd))
	m.Write([]byte(salt))
	return hex.EncodeToString(m.Sum(nil))
}

func ValidatePwd(pwd string) bool {
	if len(pwd) < 6 {
		return false
	}
	return true
}

// 数组去重
func ArrayStringDuplicateRemoval(l []string) []string {
	var arr = make([]string, len(l))
loop:
	for _, v := range l {
		for _, had := range arr {
			if had == v {
				continue loop
			}
		}
		arr = append(arr, v)
	}
	return arr
}
