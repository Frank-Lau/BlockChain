package main

import (
	"strings"
	"fmt"
	"bytes"
)

func main(){
	strsA := []string{"hello","world","golang"}

	strRes :=strings.Join(strsA,"=")
	fmt.Printf("strRes:%s\n",strRes)


	//func Join(s [][]byte, sep []byte) []byte {将二位切片用特定符号链接并返回一个以为切片
	joinRes := bytes.Join([][]byte{[]byte("hello"),[]byte("world"),[]byte("golang")},[]byte{})

	fmt.Printf("joinRes:%s\n",joinRes)
}