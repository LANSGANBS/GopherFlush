package test

// 测试硬编码URL
const APIEndpoint = "https://api.example.com/v1"

// 测试硬编码API token
const APIToken = "sk-1234567890abcdef1234567890abcdef"

// 测试硬编码数据库连接
var DBConnection = "mysql://user:password@localhost:3306/mydb"

// 测试硬编码密钥
var SecretKey = "my-secret-key-12345"

// 正确的示例 - 使用占位符（不应该被检测）
const APIEndpointPlaceholder = "https://your-api-endpoint.com"

// 正确的示例 - 短字符串（不应该被检测）
const AppName = "MyApp"
