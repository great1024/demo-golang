package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)
/*
 *https://www.whatsmydns.net/api/details?server=325&type=A&query=amazon.com
 *有用的网站，
 *没用的小程序
 *
 */
func main() {
	//reg := regexp.MustCompile("\\d{3}ms{1}")
	hostChose("amazon.com")
}
/**
 * 默认索引 20 行
 * 首先查询 包含所属域名的 行 有的话就修改没有的话 从最后一行追加
 */
func editHost(host string,delay Delay){
	f, err := os.Open("C:\\Windows\\System32\\drivers\\etc\\hosts")
	if err != nil {
		fmt.Println("read fail")
		//return ""
	}
	defer f.Close()
	var chunk []byte
	buf := make([]byte, 1024)
	for {
		//从file读取到buf中
		n, err := f.Read(buf)
		if err != nil && err != io.EOF{
			fmt.Println("read buf fail", err)
			//return ""
		}
		//说明读取结束
		if n == 0 {
			break
		}
		//读取到最终的缓冲区中
		chunk = append(chunk, buf[:n]...)
	}
	oldHosts := ConvertByte2String(chunk,"GB18030")
	hostLineList := strings.Split(oldHosts,"\n")
	 index :=-1
	for i := 0; i < len(hostLineList); i++ {
		if strings.Contains(hostLineList[i],host) {
			index = i
		}
	}
	if index > -1{
		hostLineList[index] = delay.ip +" " + host
	}else {
		hostLineList = append(hostLineList,"\n"+delay.ip+" "+host)
	}
	newHost := ""
	for _,s := range hostLineList {
		newHost += s
	}
	newHostBytes := []byte(newHost)
	//os.Remove("C:\\Windows\\System32\\drivers\\etc\\hosts")
	//_,err1 := f.Write(newHostBytes)
	fileName :="hosts_"+strconv.FormatInt(time.Now().Unix(),10)
	filePath := "./"+fileName
	err1 := ioutil.WriteFile(filePath, newHostBytes, os.ModeDevice)
	if err1 != nil {
		fmt.Println(err1)
	}
}


type Delay struct {
	ip string
	delayInt int    // 20
	delay string	// 20ms
	describe string // 描述
}
// 获取IP 地址
func networkDelay(ip string) Delay {
	reg := regexp.MustCompile("平均 = \\d{3}ms{1}")
	result := runCmd("ping "+ip)
	result = ConvertByte2String([]byte(result),"GB18030")
	//fmt.Println(result)
	networkDelays := reg.FindAllStringSubmatch(result,-1)
	var networkDelay = networkDelays[0][0]
	//fmt.Println(networkDelay)
	networkDelayString :=  networkDelay[9:12]
	networkDelayInt,err  :=  strconv.Atoi(networkDelayString)
	if err != nil{
		fmt.Println(err)
	}
	d := Delay{
		ip: ip,
		delayInt: networkDelayInt,
		delay: networkDelay,
		describe:result,
	}
	return d
}
func ipSearch(host string) []Delay{
	var delays []Delay
	set := make(map[string]int)
	result := runCmd("nslookup "+host)
	results := strings.Split(result, " ")
	ips := results[10:]
	for _, ip := range ips {
		ip = strings.Trim(ip," ")
		if ip != ""{
			ip = strings.ReplaceAll(ip," ","")
			_, inSet := set[ip]
			if ip != "" && len(ip) !=0 && !inSet{
				set[ip] = 0
				ip = strings.ReplaceAll(ip,"\r","")
				ip = strings.ReplaceAll(ip,"\n","")
				ip = strings.ReplaceAll(ip,"\t","")
				delays = append(delays,  networkDelay(ip))
			}
		}
	}
	return delays
}
func hostChose(domain string){
	ipList := ipSearch(domain)
	sort.Slice(ipList,func(i,j int) bool {
		return ipList[i].delayInt < ipList[j].delayInt
	})
	fmt.Println(ipList)
	if len(ipList) > 0 {
		editHost(domain,ipList[0])
	}else {
		fmt.Println("没有找到域名对应的IP地址！")
	}

}



func runCmd(cmdOrder string)string{
	list := strings.Split(cmdOrder," ")
	cmd := exec.Command(list[0],list[1:]...)
	//c := exec.Command("cmd", "/C", "nslookup", "amazon.com")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stderr.String()
	} else {
		return out.String()
	}
}
