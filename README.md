# GopherFlush

AI生成代码质量检测工具 - 帮助你清理代码中的"屎"

## 📋 项目简介

GopherFlush 是一个专门用于检测 AI 生成代码质量问题的工具。它可以识别常见的代码异味，如过大的文件和函数、全局变量滥用、重复代码等，帮助你保持代码库的健康和可维护性。

**支持的语言**：目前仅支持 **Go 语言**（使用 Go AST 解析器）

## ✨ 功能特性

### 已实现的检测规则

**"一坨大的"检测**

- 文件大小检测：识别超过 800 行的文件
- 函数大小检测：识别超过 200 行的函数
- 严重程度分级：低、中等、严重、极其严重（带颜色标记）

### 核心功能

- ✅ 交互式命令行界面
- ✅ 可扩展的规则系统（注册中心模式）
- ✅ 多种输出格式（控制台、JSON）
- ✅ 彩色输出（严重程度可视化）
- ✅ 详细的检测报告
- ✅ 灵活的配置系统

## 🚀 安装

### 前置要求

- Go 1.19 或更高版本

### 方法一：从源码编译

1. 克隆仓库

```bash
git clone <repository-url>
cd GopherFlush
```

2. 下载依赖

```bash
go mod tidy
```

3. 编译项目

```bash
# Windows
go build -o bin/gopherflush.exe ./cmd/gopherflush

# Linux/Mac
go build -o bin/gopherflush ./cmd/gopherflush
```

### 方法二：直接安装

```bash
go install ./cmd/gopherflush
```

安装后，`gopherflush` 命令将被添加到你的 `$GOPATH/bin` 目录中。

## 🎯 启动和使用

### 启动方式

**方式一：交互模式（推荐）**

直接运行工具进入交互式命令行：

```bash
# Windows
.\bin\gopherflush.exe

# Linux/Mac
./bin/gopherflush

# 或者如果已安装到 $GOPATH/bin
gopherflush
```

**方式二：一次性检测模式**

使用命令行参数直接执行检测：

```bash
gopherflush -path ./src -rules file-size,function-size
```

### 交互模式命令

进入交互模式后，可以使用以下命令：

#### `/help` - 查看帮助

```
gopherflush> /help
```

显示所有可用命令及其用法。

#### `/show` - 显示信息

```
gopherflush> /show flags      # 显示所有命令行参数
gopherflush> /show rules      # 显示所有检测规则
gopherflush> /show config     # 显示配置文件格式
```

#### `/run` - 运行检测

```
gopherflush> /run                                    # 检测当前目录
gopherflush> /run ./src                              # 检测指定目录
gopherflush> /run ./src --rules=file-size            # 只运行特定规则
gopherflush> /run . --rules=file-size,function-size  # 运行多个规则
```

#### `/exit` - 退出程序

```
gopherflush> /exit
```

### 一次性模式参数

```bash
gopherflush [选项]

选项:
  -path string
        要检测的代码路径 (默认 ".")
  -rules string
        要运行的规则（逗号分隔，留空表示运行所有规则）
  -config string
        配置文件路径
  -format string
        输出格式（console/json） (默认 "console")
  -output string
        详细报告输出路径 (默认 "gopherflush-report.json")
  -interactive
        启动交互模式
```

## 📖 使用示例

### 示例 1：快速开始

```bash
# 启动交互模式
gopherflush

# 在交互模式中运行检测
gopherflush> /run ./src
```

### 示例 2：检测特定规则

```bash
# 只检测文件大小
gopherflush -path ./src -rules file-size

# 检测文件和函数大小
gopherflush -path ./src -rules file-size,function-size
```

### 示例 3：生成 JSON 报告

```bash
gopherflush -path ./src -format json -output report.json
```

### 示例 4：使用自定义配置

```bash
gopherflush -path ./src -config custom.yaml
```

## 🎨 输出说明

### 严重程度颜色标记

- 🔴 **极其严重**（红色）：需要立即处理的严重问题
- 🟠 **严重**（橙色）：应该尽快处理的问题
- 🟡 **中等**（黄色）：建议处理的问题
- 🟢 **低**（绿色）：可以考虑优化的问题

### 报告内容

**控制台输出**：

- 检测文件总数
- 问题总数
- 按严重程度统计
- 按文件分组的详细信息

**JSON 报告**：

- 完整的违规记录
- 文件路径、行号、列号
- 违规描述和修复建议
- 保存位置：`gopherflush-report.json`（默认）

## 🔧 配置文件

配置文件使用 YAML 格式，示例：

```yaml
rules:
  # 文件大小检测
  file_size:
    enabled: true
    max_lines: 800

  # 函数大小检测
  function_size:
    enabled: true
    max_lines: 200

  # 全局变量检测
  global_vars:
    enabled: true

  # 重复代码检测
  duplicates:
    enabled: true

# 输出配置
output:
  format: console # console 或 json
  path: "" # JSON输出路径

# 排除目录
exclude:
  - vendor/
  - node_modules/
  - .git/
  - "**/*_test.go"
```

## 📏 检测规则详解

### 1. 文件大小检测 (file-size)

检测文件是否超过指定行数（默认 800 行）。

**严重程度分级**：

- 800-1000 行：低
- 1000-1500 行：中等
- 1500-2000 行：严重
- 2000 行以上：极其严重

**建议**：将大文件拆分为多个较小的文件，提高代码可维护性。

### 2. 函数大小检测 (function-size)

检测函数是否超过指定行数（默认 200 行）。

**严重程度分级**：

- 200-300 行：低
- 300-500 行：中等
- 500-800 行：严重
- 800 行以上：极其严重

**建议**：将大函数拆分为多个较小的函数，提高代码可读性和可维护性。

### 3. 全局变量检测 (global-vars)

检测全局变量的使用和滥用。

**严重程度分级**：

- 可导出的全局变量（大写开头）：严重
- 不可导出的全局变量（小写开头）：中等

**建议**：考虑将全局变量改为局部变量、函数参数或使用依赖注入模式。

### 4. 重复代码检测 (duplicates)

检测重复的变量或函数声明（开发中）。

## 🏗️ 项目架构

```
GopherFlush/
├── cmd/gopherflush/          # 命令行入口
│   └── main.go
├── internal/                 # 内部实现
│   ├── analyzer/            # 代码分析器
│   │   ├── analyzer.go      # 主分析器
│   │   └── parser.go        # Go AST 解析器
│   ├── cli/                 # 命令行交互
│   │   ├── command.go       # 命令接口
│   │   ├── registry.go      # 命令注册中心
│   │   ├── repl.go          # 交互式命令行
│   │   └── builtin/         # 内置命令
│   │       ├── help.go      # /help 命令
│   │       ├── show.go      # /show 命令
│   │       ├── run.go       # /run 命令
│   │       └── exit.go      # /exit 命令
│   ├── rules/               # 规则引擎
│   │   ├── rule.go          # 规则接口
│   │   ├── registry.go      # 规则注册表
│   │   └── builtin/         # 内置规则
│   │       ├── file_size.go
│   │       ├── function_size.go
│   │       ├── global_vars.go
│   │       └── duplicates.go
│   ├── reporter/            # 报告生成器
│   │   ├── reporter.go      # 报告器接口
│   │   ├── console.go       # 控制台输出
│   │   ├── json.go          # JSON 输出
│   │   └── color.go         # 颜色支持
│   └── config/              # 配置管理
│       ├── config.go        # 配置结构
│       └── loader.go        # 配置加载器
├── pkg/types/               # 公共类型定义
│   └── violation.go         # 违规记录类型
├── configs/                 # 配置文件
│   └── default.yaml
├── docs/                    # 文档
├── test/                    # 测试文件
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## 🔌 扩展指南

### 添加新的检测规则

1. 在 `internal/rules/builtin/` 创建新规则文件
2. 实现 `Rule` 接口：

```go
type Rule interface {
    Name() string
    Description() string
    Check(file *ast.File, fset *token.FileSet, filePath string) []*types.Violation
    Enabled() bool
}
```

3. 在 `main.go` 的 `initializeRules()` 函数中注册规则

### 添加新的命令

1. 在 `internal/cli/builtin/` 创建新命令文件
2. 实现 `Command` 接口：

```go
type Command interface {
    Name() string
    Description() string
    Usage() string
    Execute(args []string) error
}
```

3. 在 `main.go` 的 `startInteractiveMode()` 函数中注册命令

## 🛠️ 开发

### 编译

```bash
make build
```

### 测试

```bash
make test
```

### 清理

```bash
make clean
```

## 📝 许可证

[添加许可证信息]

## 🤝 贡献

欢迎贡献！请随时提交 Issue 或 Pull Request。

---

**GopherFlush** - 让你的代码保持清新！🚽✨
