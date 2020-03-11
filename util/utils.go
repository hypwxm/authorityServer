package util

import (
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
