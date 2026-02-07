# GopherFlush 架构文档

## 项目结构

```
gopherflush/
├── cmd/gopherflush/          # 命令行入口
├── internal/                 # 内部实现
│   ├── analyzer/            # 代码分析器
│   ├── rules/               # 规则引擎
│   │   └── builtin/        # 内置规则
│   ├── reporter/            # 报告生成器
│   └── config/              # 配置管理
├── pkg/types/               # 公共类型定义
└── configs/                 # 配置文件
```

## 核心组件

### 1. Analyzer (分析器)
- 负责遍历和解析Go源代码
- 使用Go AST进行代码分析
- 协调规则执行

### 2. Rules (规则引擎)
- 可扩展的规则系统
- 每个规则实现Rule接口
- 通过Registry进行管理

### 3. Reporter (报告器)
- 支持多种输出格式
- Console: 控制台输出
- JSON: JSON文件输出

### 4. Config (配置管理)
- YAML配置文件
- 支持规则启用/禁用
- 可配置阈值参数

## 扩展指南

### 添加新规则

1. 在 `internal/rules/builtin/` 创建新规则文件
2. 实现 `Rule` 接口
3. 在主程序中注册规则
