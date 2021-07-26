package baseController

import (
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/hypwxm/authorityServer/base_control/aliyunOss"
	"github.com/hypwxm/authorityServer/config"
	"github.com/hypwxm/authorityServer/logger"
	"github.com/hypwxm/authorityServer/util/response"

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
		fullfilename = strings.ReplaceAll(fullfilename, "\\", "/")

		// fileurl, err := aliyunOss.UploadFileStream(formFile.File, filepath.Join(action, filename))
		// if err != nil {
		// 	sender.Fail(err.Error())
		// 	return
		// }
		// sender.Success(fileurl)

		c.StoreFormFile(formFile, fullfilename)
		storePath := filepath.Join(basePath, action, filename)
		storePath = strings.ReplaceAll(storePath, "\\", "/")
		sender.Success(config.Config.UPLOADERHOST + storePath)
	})()
	c.SendJson(200, sender)
}

// 单文件上传
func uploadToAliyunOss(c rider.Context) {
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

		filename := createFileName(c) + filepath.Ext(formFile.Name)
		fileurl, err := aliyunOss.UploadFileStream(formFile.File, filepath.Join(action, filename))
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success(fileurl)
	})()
	c.SendJson(200, sender)
}

func createFileName(c rider.Context) string {
	var userId = c.GetLocals("userID").(string)
	nowDay := time.Now().Format("20060102")
	nowTime := time.Now().Format("150405")

	rand.Seed(time.Now().UnixNano())
	rn := rand.Intn(100000000000) + 100000000000
	return filepath.Join(userId, nowDay, nowTime+"a"+strconv.Itoa(rn))
}
