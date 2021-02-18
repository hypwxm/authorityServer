package aliyunOss

import (
	"math/rand"
	"strconv"
	"time"
)

func CreateFilename(userID string) string {
	now := time.Now().Format("20060102150405")
	rand.Seed(time.Now().UnixNano())
	rn := rand.Intn(1000000) + 1000000
	return now + "b" + userID + "a" + strconv.Itoa(rn)

}
