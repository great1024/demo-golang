package main

import (
	"bufio"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/schollz/progressbar/v3"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const image = "图片链接"
const storeIndexValue = "店铺名"
const asinIndexValue = "asin"
const listingIndexValue = "listing"
const aboutIndexValue = "about"
const video = "视频"
const filePath = "./image"
const productNameIndexValue = "产品名"
const sellerSKUIndexValue = "msku"
const brandIndexValue = "品牌"
var bar *progressbar.ProgressBar
func main() {
	f, err := excelize.OpenFile("主图文件.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	//读取某个表单的所有数据
	var rows = f.GetRows("Sheet1")
	//  图片链接
	imageIndex := -1
	storeIndex := -1
	asinIndex := -1
	listingIndex := -1
	aboutIndex := -1
	productNameIndex := -1
	sellerSKUIndex := -1
	brandIndex := -1
	bar = progressbar.Default(int64(len(rows)),"下载进度:")
	// 数据解析 下载
	for i, row := range rows {
		if imageIndex < 0 && i == 0  {
			for i, value := range row {
				switch value {
				case image:
					imageIndex = i
				case storeIndexValue:
					storeIndex = i
				case asinIndexValue:
					asinIndex = i
				case listingIndexValue:
					listingIndex = i
				case aboutIndexValue:
					aboutIndex = i
				case productNameIndexValue:
					productNameIndex = i
				case sellerSKUIndexValue:
					sellerSKUIndex = i
				case brandIndexValue:
					brandIndex = i
				}
			}
		}
		if imageIndex >= 0 && i != 0  && len(row[imageIndex]) > 0{
			//fmt.Printf("\n\t图片链接:%s", row[imageIndex])
			downloadByUrl( row[imageIndex],row[storeIndex], row[asinIndex],row[listingIndex],row[aboutIndex],row[productNameIndex],row[sellerSKUIndex],row[brandIndex])
		}
		bar.Add(1)
	}
}
// 根据链接下载文件
func downloadByUrl(url string,store string,asin string,listing string,about string,productName string,sellerSKU string,brand string){
	store = strings.ReplaceAll(store," ","_")
	productName = strings.ReplaceAll(productName," ","_")
	productName = strings.ReplaceAll(productName,"/","_")
	timeUnixNano := time.Now().UnixNano()
	fileName := asin+strconv.FormatInt(timeUnixNano,36)+url[strings.LastIndex(url,"."):]
	//fmt.Printf("\n\t下载链接:%s", url)
	newFilePath := filePath+"/"+store+"/"+productName+"/"
	err:=os.Mkdir(filePath,os.ModePerm)
	err =os.Mkdir(filePath+"/"+store,os.ModePerm)
	err =os.Mkdir(filePath+"/"+store+"/"+productName,os.ModePerm)
	if err!=nil{
		//文件夹存在
		//fmt.Println("文件夹存在")
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
	reader := bufio.NewReaderSize(res.Body, 32 * 1024)
	file, err := os.Create(newFilePath+fileName)
	if err != nil {
		panic(err)
	}
	// 获得文件的writer对象
	writer := bufio.NewWriter(file)
	io.Copy(writer, reader)
	context := "品牌："+brand+"\n"+"ASIN: "+asin+"\n"+"msku: "+sellerSKU+"\n"+listing+"\n"+about
	ioutil.WriteFile(newFilePath+sellerSKU+"_"+"context.txt",[]byte(context), os.ModeDevice)
}
