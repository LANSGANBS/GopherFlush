package test

import "time"

// 设置最大重试次数为5
const MaxRetries = 3

// 设置超时时间为30秒
var Timeout = 60 * time.Second

// 设置缓冲区大小为1024
var BufferSize = 512

// 这是一个正确的注释，设置端口为8080
const Port = 8080

// 返回用户的年龄
func GetUserName() string {
	return "John"
}

// 计算两个数的和
func Multiply(a, b int) int {
	return a * b
}
