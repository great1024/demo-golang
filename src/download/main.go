package main

import (
	"bufio"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/schollz/progressbar/v3"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const image = "图片"
const video = "视频"
const filePath = "file"

var bar *progressbar.ProgressBar

func main() {
	f, err := excelize.OpenFile("下载链接.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	//读取某个表单的所有数据
	var rows = f.GetRows("Sheet1")
	//  图片链接
	imageIndex := -1
	// 视频链接
	videoIndex := -1
	bar = progressbar.Default(int64(len(rows)), "下载进度:")
	// 数据解析 下载
	for i, row := range rows {
		if imageIndex < 0 || videoIndex < 0 {
			for i, value := range row {
				switch value {
				case image:
					imageIndex = i
				case video:
					videoIndex = i
				}
			}
		}
		if imageIndex >= 0 && i != 0 && len(row[imageIndex]) > 0 {
			//fmt.Printf("\n\t图片链接:%s", row[imageIndex])
			parseAndDownloadUrls(row[imageIndex], image)
		}
		if videoIndex >= 0 && i != 0 && len(row[videoIndex]) > 0 {
			//fmt.Printf("\n\t视频链接:%s", row[videoIndex])
			parseAndDownloadUrls(row[videoIndex], video)
		}
		bar.Add(1)
	}
}

// 解析链接并下载
func parseAndDownloadUrls(urls string, urlType string) {

	// 去掉像素后缀 获取高清原图
	if strings.EqualFold(urlType, image) {
		urls = strings.ReplaceAll(urls, "._SY88", "")
	}
	if strings.Contains(urls, ",") {
		imageUrlArray := strings.Split(urls, ",")
		for _, url := range imageUrlArray {
			fileName := url[strings.LastIndexAny(url, "/"):len(url)]
			downloadByUrl(url, fileName)
		}
	} else {
		fileName := urls[strings.LastIndexAny(urls, "/"):len(urls)]
		downloadByUrl(urls, fileName)
	}
}

// 根据链接下载文件
func downloadByUrl(url string, fileName string) {
	//fmt.Printf("\n\t下载链接:%s", url)
	err := os.Mkdir(filePath, os.ModePerm)
	if err != nil {
		//文件存在
		//fmt.Println(err)
	}
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("A error occurred!")
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(res.Body)
	// 获得get请求响应的reader对象
	reader := bufio.NewReaderSize(res.Body, 32*1024)
	file, err := os.Create(filePath + fileName)
	if err != nil {
		panic(err)
	}
	// 获得文件的writer对象
	writer := bufio.NewWriter(file)
	io.Copy(writer, reader)
}
