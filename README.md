# 📈 TDX 通达信股票数据查询系统

> 基于通达信协议的股票数据获取库 + Web 可视化界面 + RESTful API

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Docker](https://img.shields.io/badge/Docker-支持-2496ED?style=flat&logo=docker)](https://www.docker.com)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

**感谢源作者 [injoyai](https://github.com/injoyai/tdx)，请支持原作者！**

---

## ✨ 功能特性

| 分类               | 功能                                                                      |
| ------------------ | ------------------------------------------------------------------------- |
| **📊 核心功能**    | 实时行情（五档盘口）、K 线数据（10 种周期）、分时数据、股票搜索、批量查询 |
| **🌐 Web 界面**    | 现代化 UI、ECharts 图表、智能搜索、实时刷新                               |
| **🔌 RESTful API** | 32 个接口、完整文档、多语言示例、高性能                                   |
| **🐳 Docker 部署** | 开箱即用、国内镜像加速、跨平台支持                                        |

---

## 🚀 快速开始

### 方式一：Docker 部署（推荐）⭐

```bash
# 克隆项目
git clone https://github.com/oficcejo/tdx-api.git
cd tdx-api

# 启动服务（已配置国内镜像加速）
docker-compose up -d

# 默认访问 http://localhost:8080
```

**一键启动脚本：**

- Windows: 双击 `docker-start.bat`
- Linux/Mac: `chmod +x docker-start.sh && ./docker-start.sh`

### 方式二：源码运行

```bash
# 前置要求: Go 1.22+

# 1. 下载依赖
go mod download

# 2. 进入web目录并运行
cd web
go run .

# 3. 默认访问 http://localhost:8080
```

> ⚠️ **注意**: 必须使用 `go run .` 编译所有 Go 文件，不能使用 `go run server.go`

可选环境变量：

- `ENV_TDX_API_HOST`：监听地址，默认 `localhost`
- `ENV_TDX_API_PORT`：监听端口，默认 `8080`

示例：

```bash
cd web
ENV_TDX_API_HOST=127.0.0.1 ENV_TDX_API_PORT=9090 go run .
```

---

## API 文档

README 不再维护接口清单，避免与实现和示例漂移。

- 完整接口定义、请求参数、返回结构、错误码： [API\_接口文档.md](API_接口文档.md)
- 接入流程、接口组合建议、超时重试、排坑： [API_集成指南.md](API_集成指南.md)
- 文档分工说明： [API_完成总结.md](API_完成总结.md)
- 如果你在 Codex 环境里查接口，可结合仓库内 skill： [skills/tdx-api-docs/SKILL.md](skills/tdx-api-docs/SKILL.md)

---

## 使用示例

### API 调用

```bash
# 获取实时行情
curl "http://localhost:8080/api/quote?code=000001"

# 获取日K线
curl "http://localhost:8080/api/kline?code=000001&type=day"

# 搜索股票
curl "http://localhost:8080/api/search?keyword=平安"

# 健康检查
curl "http://localhost:8080/api/health"
```

### Go 库使用

```go
import "github.com/injoyai/tdx"

// 连接服务器
c, _ := tdx.DialDefault(tdx.WithDebug(false))

// 获取行情
quotes, _ := c.GetQuote("000001", "600519")

// 获取日K线
kline, _ := c.GetKlineDayAll("000001")
```

---

## Docker 配置说明

### 国内镜像加速

Docker 配置已使用国内镜像源，加速构建：

| 组件        | 镜像源                                             |
| ----------- | -------------------------------------------------- |
| Go 基础镜像 | `registry.cn-hangzhou.aliyuncs.com/library/golang` |
| Alpine 镜像 | `registry.cn-hangzhou.aliyuncs.com/library/alpine` |
| Alpine APK  | `mirrors.aliyun.com`                               |
| Go Proxy    | `goproxy.cn` + `mirrors.aliyun.com/goproxy`        |

### 常用命令

```bash
docker-compose up -d       # 启动服务
docker-compose logs -f     # 查看日志
docker-compose stop        # 停止服务
docker-compose restart     # 重启服务
docker-compose down        # 完全清理
```

**详细部署文档**: [DOCKER_DEPLOY.md](DOCKER_DEPLOY.md)

---

## 📊 支持的数据类型

| 数据类型               | 方法                | 说明                         |
| ---------------------- | ------------------- | ---------------------------- |
| 五档行情               | `GetQuote`          | 实时买卖五档、最新价、成交量 |
| 1/5/15/30/60 分钟 K 线 | `GetKlineXXXAll`    | 分钟级 K 线数据              |
| 日/周/月 K 线          | `GetKlineDayAll` 等 | 中长期 K 线数据              |
| 分时数据               | `GetMinute`         | 当日每分钟价格               |
| 分时成交               | `GetTrade`          | 逐笔成交记录                 |
| 股票列表               | `GetCodeAll`        | 全市场代码                   |

---

## 📁 项目结构

```
tdx-api/
├── client.go              # TDX客户端核心
├── protocol/              # 通达信协议实现
├── web/                   # Web应用
│   ├── server.go          # 主服务器
│   ├── server_api_extended.go  # 扩展API
│   ├── tasks.go           # 任务管理
│   └── static/            # 前端文件
├── extend/                # 扩展功能
├── Dockerfile             # Docker镜像（国内源）
├── docker-compose.yml     # Docker编排
└── docs/                  # 文档
```

---

## Skill 使用

仓库内提供了一个面向 API 文档问答的 skill：

- Skill 路径：`skills/tdx-api-docs/SKILL.md`
- 适用场景：查询接口是否已实现、参数含义、返回结构、字段单位、HTTP/WebSocket 调用示例、接入方式建议
- 校验顺序：优先看 `API_接口文档.md`，涉及接入策略时再看 `API_集成指南.md`，最后回查 `web/server.go`、`web/server_api_extended.go`、`web/server_ws.go`

在支持 skill 的 Codex 环境里，可以直接这样使用：

```text
用 tdx-api-docs 看一下 /api/kline-history 的返回结构
帮我确认 /ws/quote 是否已经实现
解释一下 /api/trade-history/full 的字段含义
```

这个 skill 的目标是让 API 对接时优先以当前代码行为为准，减少文档和实现不一致带来的误判。

---

## 相关资源

| 资源        | 链接                                          |
| ----------- | --------------------------------------------- |
| 原项目      | [injoyai/tdx](https://github.com/injoyai/tdx) |
| API 文档    | [API\_接口文档.md](API_接口文档.md)           |
| 集成指南    | [API_集成指南.md](API_集成指南.md)           |
| 文档说明    | [API_完成总结.md](API_完成总结.md)           |
| Docker 部署 | [DOCKER_DEPLOY.md](DOCKER_DEPLOY.md)          |
| Python 示例 | [API\_使用示例.py](API_使用示例.py)           |
| API Skill   | [skills/tdx-api-docs/SKILL.md](skills/tdx-api-docs/SKILL.md) |

### 通达信服务器

系统自动连接最快的服务器：

| IP             | 地区       |
| -------------- | ---------- |
| 124.71.187.122 | 上海(华为) |
| 122.51.120.217 | 上海(腾讯) |
| 121.36.54.217  | 北京(华为) |
| 124.71.85.110  | 广州(华为) |

---

## ⚠️ 免责声明

1. 本项目仅供学习和研究使用
2. 数据来源于通达信公共服务器，可能存在延迟
3. 不构成任何投资建议，投资有风险

---

## 📄 许可证

MIT License - 详见 [LICENSE](LICENSE)

---

**如果这个项目对您有帮助，请点个 Star ⭐ 支持一下！**
