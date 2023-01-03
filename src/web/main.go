package main

import (
	"errors"
	"fmt"
)

func main() {
	//gotoFunc()
	//errFunc()
	//iotaFunc()
	//可变参
	//paramFunc(1,2,3)
	//数组变长
	arrFunc()
	fmt.Print(
		"666")
}

func gotoFunc(){
	Here:
		a := 1
		fmt.Print(a)
		goto Here
}

func errFunc(){
	err := errors.New("创建了一个错误")
	if err != nil {
		fmt.Print(err)
	}
}

// 数组改变长度 是否需要传地址测试
func arrFunc(){
	var arr = []int{1,2,3}

	fmt.Printf("调用前数组的内容为：%v",arr)
	arrAdd(&arr)
	fmt.Printf("调用后数组的内容为：%v",arr)
}

func arrAdd(arr *[]int) {
	//arr[3]
}
//iota

const (
	x = iota
	y = iota
	z = iota
	w
)

const v = iota

const h, i, j = iota,iota,iota

const (
	a = iota
	b = 'B'
	c = iota
	d,e,f = iota,iota,iota
	g = iota
)

func iotaFunc() {
	//可以看出 iota的值只与枚举常量定义顺序有关
	fmt.Print(a,b,c,d,e,f,g,g,h,i,j,x,y,z,w,v)
}

func paramFunc(a ...int){
	fmt.Print(len(a))
}
