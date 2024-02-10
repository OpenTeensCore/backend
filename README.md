# OpenTeens 社区后端开发文档 📖

## 项目介绍 🚀

OpenTeens社区后端项目是由OpenTeensCore团队开发的，旨在为前端提供必要的API接口和技术支持。本项目基于Go语言和Gin框架，以提供高效、稳定和易扩展的服务。

## 开始之前 🛠️

确保您的开发环境中安装了以下软件：

- Go (版本1.15或更高)
- Git
- 任何Go语言支持的IDE（推荐使用Visual Studio Code或GoLand）

## 如何下载项目 📦

1. 打开终端或命令提示符。
2. 运行以下命令来克隆项目仓库：

```bash
git clone https://github.com/OpenTeensCore/openteens-backend.git
```

3. 进入项目目录：

```bash
cd openteens-backend
```

## 启动项目 🚀

在项目根目录下，执行以下命令来启动项目：

```bash
go run main.go
```

这将启动服务器，默认监听在`localhost:8080`。您可以通过访问`http://localhost:8080`来测试是否成功运行。

## 进行开发 👨‍💻👩‍💻

### 设置开发环境

建议使用Visual Studio Code或GoLand作为开发IDE，它们提供了良好的Go语言支持和便捷的调试工具。

### 目录结构

- `/docs` - 存放文档
- `/controller` - 存放路由处理函数
- `/services` - 业务逻辑处理（被controller层调用）【待重构】
- `/model` - 数据模型
- `/middleware` - 中间件
- `/utils` - 工具函数
- `/dao` - 数据库相关配置
- `/router` - 路由配置

### 添加新的API

1. 在`/controller`目录和`/service`目录下创建相应的文件。
2. 定义路由和处理函数。
3. 在`/router`中注册新的路由。

### 编写代码

请遵循[Go代码评审标准](https://github.com/golang/go/wiki/CodeReviewComments)来保持代码质量。

## Git提交规范 📝

**经过讨论，commit消息使用**

创建分支的时候请遵循分支分类+名称的命名规则，如`feat/sqlite`。如果有很多修改可以用名字/分类-名称的方式，如`lanbinshijie/sqlite`。这些分支名经过测试都是合法的。

为了提高项目的可维护性，请遵循以下Git提交信息规范：

- 功能添加：`feat: 添加了新的登录功能`
- 问题修复：`fix: 修复了用户认证bug`
- 文档更新：`docs: 更新了README文件`
- 性能优化：`perf: 优化了数据库查询效率`
- 代码重构：`refactor: 重构了用户服务模块` 或 `ref: 重构了用户服务模块`
- 测试代码：`test: 添加了新的API测试用例`

## 贡献指南 🤝

欢迎通过Pull Requests或Issue来贡献您的智慧。在提交PR之前，请确保您的代码符合上述开发和Git提交规范。
