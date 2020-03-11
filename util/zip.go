package util

import (
	"archive/zip"
	"io"
	"os"
	"strings"
)

// ZipFiles compresses one or many files into a single zip archive file.
// 压缩多个文件到一个文件里面
// Param 1: 输出的zip文件的名字
// Param 2: 需要添加到zip文件里面的文件
// Param 3: 由于file是绝对路径，打包后可能不是想要的目录，oldform就是filename中需要被替换的掉的路径
// Param 4: 要替换成的路径
func ZipFiles(filename string, files []string, oldform, newform string) error {
	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// 把files添加到zip中
	for _, file := range files {

		zipfile, err := os.Open(file)
		if err != nil {
			return err
		}
		defer zipfile.Close()

		// 获取file的基础信息
		info, err := zipfile.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		//使用上面的FileInforHeader() 就可以把文件保存的路径替换成我们自己想要的了，如下面
		header.Name = strings.Replace(file, oldform, newform, -1)

		// 优化压缩
		// 更多参考see http://golang.org/pkg/archive/zip/#pkg-constants
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		if _, err = io.Copy(writer, zipfile); err != nil {
			return err
		}
	}
	return nil
}
