package test

import "fmt"

// 这个文件用于测试重复代码检测

// ProcessUserData 处理用户数据
func ProcessUserData(name string, age int) {
	fmt.Printf("Processing user: %s, age: %d\n", name, age)
	if age < 18 {
		fmt.Println("User is a minor")
	} else {
		fmt.Println("User is an adult")
	}
}

// ProcessCustomerData 处理客户数据（与ProcessUserData重复）
func ProcessCustomerData(name string, age int) {
	fmt.Printf("Processing user: %s, age: %d\n", name, age)
	if age < 18 {
		fmt.Println("User is a minor")
	} else {
		fmt.Println("User is an adult")
	}
}

// ValidateEmail 验证邮箱
func ValidateEmail(email string) bool {
	if len(email) == 0 {
		return false
	}
	return true
}

// CheckEmail 检查邮箱（与ValidateEmail重复）
func CheckEmail(email string) bool {
	if len(email) == 0 {
		return false
	}
	return true
}

// UniqueFunction 这是一个独特的函数
func UniqueFunction() {
	fmt.Println("This is unique")
}
