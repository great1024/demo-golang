package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)



func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	excelParse()
}

type Vertex struct {
	value string
}

var m map[string]Vertex

func excelParse(){
	filePath := "target.xlsx"
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		println(err.Error())
		return
	}
	var cell = f.GetCellValue("Sheet1", "B2")
	if err != nil {
		println(err.Error())
		return
	}
	println(cell)
	lable := make(map[string]Vertex)
	// Get all the rows in the Sheet1.
	var rows = f.GetRows("Sheet1")
	for _, row := range rows {
		lable[row[0]] = Vertex{
			row[1],
		}
	}
	for i, row := range rows {
		for k, v  := range lable {
			if strings.Contains(row[3],k) && i > 0 {
				fmt.Print(i, "\t")
				cell  := "E"
				cellValue := ""
				if strings.EqualFold(row[4],"")||len(row[4]) ==0 {
					cellValue = v.value
				}else if !strings.Contains(row[4],v.value){
					cellValue = v.value + "," + row[4]
				}else if strings.Contains(row[4],v.value){
					cellValue = row[4]
				}
				f.SetCellValue("sheet1",cell+ strconv.Itoa(i+1), cellValue)
			}
		}

	}

	_ = f.Save()

}

func textParse(){
	// 使用 io/ioutil.ReadFile 方法一次性将文件读取到内存中
	filePath := "email.txt"
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		// log.Fatal(err)
		fmt.Printf("读取文件出错：%v", err)
	}
	emailList := strings.Fields(string(content))
	var emails string
	for _,val := range emailList {
		emails +="'"+val+"',"
	}
	var filename = "output.txt"
	var f *os.File
	var err1 error
	if checkFileIsExist(filename) { //如果文件存在
		f, err1 = os.OpenFile(filename, os.O_APPEND, 0666) //打开文件
		fmt.Println("文件存在")
	} else {
		f, err1 = os.Create(filename) //创建文件
		fmt.Println("文件不存在")
	}

	check(err1)
	sql := "sql：update br_customer set manager1_id = (select id from br_user where username like '用户名')  where email in ("+emails[:len(emails)-1]+")"
	n, err1 := io.WriteString(f, sql) //写入文件(字符串)
	check(err1)
	fmt.Printf("写入 %d 个字节n", n)
	_ = f.Close()
}


