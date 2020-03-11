package baseController

import (
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"worldbar/logger"
	"worldbar/util/response"

	"worldbar/config"

	"github.com/hypwxm/rider"
	"github.com/hypwxm/rider/utils/file"
	"github.com/sirupsen/logrus"
)

var cwd string

const basePath = "/assets/files"

func init() {
	cwd, _ = os.Getwd()
}

// 单文件上传
func upload(c rider.Context) {
	var err error
	sender := response.NewSender()
	(func() {
		var formFile *rider.UploadFile
		formFile, err = c.FormFile("file")
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{
				"event": "上传文件-读取文件",
			}).Error(err)
			sender.Fail("系统错误")
			return
		}
		var action string = c.QueryDefault("action", "upload")
		fullDirpath := filepath.Join(cwd, basePath, action)
		if !file.IsDir(fullDirpath) {
			err = os.MkdirAll(fullDirpath, os.ModePerm)
			if err != nil {
				logger.Logger.WithFields(logrus.Fields{
					"event": "上传文件-创建目录",
				}).Error(err)
				sender.Fail("系统错误")
				return
			}
		}
		filename := createFileName(c) + filepath.Ext(formFile.Name)
		fullfilename := filepath.Join(fullDirpath, filename)
		c.StoreFormFile(formFile, fullfilename)
		storePath := filepath.Join(basePath, action, filename)
		sender.Success(config.Config.UPLOADERHOST + storePath)
	})()
	c.SendJson(200, sender)
}

func createFileName(c rider.Context) string {
	var uploadUser string
	/* var ok bool
	if uploadUser, ok = c.Jwt().Get("user").(string); !ok {
		uploadUser = "cu"
	} */
	now := time.Now().Format("20060102150405")
	rand.Seed(time.Now().UnixNano())
	rn := rand.Intn(1000000) + 1000000
	return now + "b" + uploadUser + "a" + strconv.Itoa(rn)
}
