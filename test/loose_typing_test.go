package test

// 测试宽泛类型 - 参数使用 interface{}
func StoreUserInfo(data interface{}) error {
	// 存储用户信息
	return nil
}

// 测试宽泛类型 - 返回值使用 interface{}
func GetConfig(key string) interface{} {
	return nil
}

// 测试宽泛类型 - 使用 any (Go 1.18+)
func ProcessData(input any) any {
	return input
}

// 正确的示例 - 使用具体类型
func SaveUser(name string, age int) error {
	return nil
}

// 正确的示例 - 使用明确的接口
type Reader interface {
	Read(p []byte) (n int, err error)
}

func ReadData(r Reader) ([]byte, error) {
	buf := make([]byte, 1024)
	n, err := r.Read(buf)
	return buf[:n], err
}
