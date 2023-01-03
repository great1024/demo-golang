package main

import (
	"golang.org/x/text/encoding/simplifiedchinese"
)

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

//func main() {
//	command := "ping"
//	params := []string{"127.0.0.1","-t"}
//	cmd := exec.Command(command, params...)
//	stdout, err := cmd.StdoutPipe()
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	cmd.Start()
//	in := bufio.NewScanner(stdout)
//	for in.Scan() {
//		cmdRe:=ConvertByte2String(in.Bytes(),"GB18030")
//		fmt.Println(cmdRe)
//	}
//	cmd.Wait()
//}

func ConvertByte2String(byte []byte, charset Charset) string {
	var str string
	switch charset {
	case GB18030:
		var decodeBytes,_=simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str= string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}
	return str
}



