package test

import (
	"database/sql"
	"net/http"
	"os"
)

// 测试文件资源泄漏
func readFileWithoutDefer() error {
	file, err := os.Open("test.txt")
	if err != nil {
		return err
	}
	// 忘记 defer file.Close()

	// 读取文件内容
	buf := make([]byte, 1024)
	_, _ = file.Read(buf)
	return nil
}

// 测试HTTP资源泄漏
func httpRequestWithoutDefer() error {
	resp, err := http.Get("https://example.com")
	if err != nil {
		return err
	}
	// 忘记 defer resp.Body.Close()

	// 读取响应内容
	buf := make([]byte, 1024)
	_, _ = resp.Body.Read(buf)
	return nil
}

// 测试数据库连接泄漏
func databaseWithoutDefer() error {
	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		return err
	}
	// 忘记 defer db.Close()

	// 使用数据库连接
	_ = db.Ping()
	return nil
}

// 正确的示例：使用了defer
func readFileWithDefer() error {
	file, err := os.Open("test.txt")
	if err != nil {
		return err
	}
	defer file.Close() // 正确使用defer

	// 读取文件内容
	return nil
}
