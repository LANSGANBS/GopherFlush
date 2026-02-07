package test

import "fmt"

// 可导出的全局变量
var GlobalConfig string = "config"
var GlobalCounter int = 0

// 不可导出的全局变量
var internalCache map[string]string
var debugMode bool = true

func UseGlobalVars() {
	fmt.Println(GlobalConfig)
	GlobalCounter++
	internalCache = make(map[string]string)
	if debugMode {
		fmt.Println("Debug mode enabled")
	}
}
