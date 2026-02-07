package test

import "fmt"

// func add(a, b int) int {
//     return a + b
// }

func add(a ...int) int {
	sum := 0
	for _, v := range a {
		sum += v
	}
	return sum
}

// var oldConfig = "old"
var newConfig = "new"

// const MaxRetries = 3
const MaxRetries = 5

func process() {
	// TODO: 这是一个真正的注释，不应该被检测
	fmt.Println("Processing...")

	// if debug {
	//     fmt.Println("Debug mode")
	// }
}
