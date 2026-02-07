package test

import "database/sql"

// 测试缺少验证 - 直接使用变量构建SQL（SQL注入风险）
func GetUserByID(userID string) error {
	db, _ := sql.Open("mysql", "connection")
	defer db.Close()

	// 危险：直接使用变量，没有验证
	query := "SELECT * FROM users WHERE id = " + userID
	_, err := db.Exec(query)
	return err
}

// 测试缺少验证 - 字符串拼接构建SQL
func DeleteUser(username string) error {
	db, _ := sql.Open("mysql", "connection")
	defer db.Close()

	// 危险：使用字符串拼接构建SQL
	_, err := db.Exec("DELETE FROM users WHERE name = '" + username + "'")
	return err
}

// 正确的示例 - 使用参数化查询
func GetUserSafe(userID int) error {
	db, _ := sql.Open("mysql", "connection")
	defer db.Close()

	// 安全：使用参数化查询
	_, err := db.Exec("SELECT * FROM users WHERE id = ?", userID)
	return err
}
