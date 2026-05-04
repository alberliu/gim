# GIM

[![Go](https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![GitHub](https://img.shields.io/github/stars/alberliu/gim?style=social)](https://github.com/alberliu/gim)

GIM 是一个使用 Go 语言开发的高性能即时通讯服务器，可作为业务系统的 IM 组件快速集成。

## 特性

- **多协议支持** - 同时支持 TCP 和 WebSocket 接入
- **多设备同步** - 单用户多设备同时在线，消息实时同步
- **离线消息** - 完整的离线消息存储与同步机制
- **多场景支持** - 支持单聊、群聊、房间聊天等多种场景
- **水平扩展** - 支持服务水平扩展，轻松应对高并发
- **容器化部署** - 支持 Docker Compose 和 Kubernetes 部署
- **领域驱动设计** - 代码结构清晰，易于维护和扩展

## 技术栈

| 组件 | 技术 |
|------|----|
| 语言 | Go |
| 数据库 | MySQL |
| 缓存 | Redis |
| 通讯框架 | gRPC |
| 传输协议 | Protocol Buffers |
| ORM | GORM |

## 快速开始

### Docker Compose 单机部署

```bash
# 一键部署
./deploy_compose.sh
```

### Kubernetes 部署

```bash
./deploy_k8s.sh
```

## 项目结构

项目结构遵循 [golang-standards/project-layout](https://github.com/golang-standards/project-layout)

```
gim/
├── cmd/           # 服务启动入口
├── config/        # 服务配置文件
├── deploy/        # 部署配置（Docker、K8s）
├── docs/          # 项目文档
├── internal/      # 私有代码（不对外暴露）
├── pkg/           # 公共代码（可被外部引用）
├── sql/           # 数据库脚本
├── test/          # 测试代码
└── web/           # Web IM 客户端（Vue 3 + Vite）
```

## 服务架构

GIM 由三个核心服务组成：

### 1. Connect 服务
长连接管理服务，负责：
- 维持客户端 TCP/WebSocket 长连接
- 心跳检测与连接保活
- TCP 拆包粘包处理
- 消息编解码

### 2. Logic 服务
业务逻辑服务，负责：
- 设备信息管理
- 好友关系管理
- 群组信息管理
- 消息转发与路由

### 3. Business 服务
业务扩展服务，提供：
- 用户注册登录
- 基础鉴权功能
- 可替换为自己的业务服务

> 如需接入自有业务系统，只需实现 `user.int.proto` 中的 `UserIntService.Auth` 接口即可。

## 客户端接入

### 接入流程

```
1. 登录获取凭证
   └─► 调用 business.UserExtService.SignIn
   └─► 获取 device_id, user_id, token

2. 建立长连接
   └─► Web 端：WebSocket 连接
   └─► APP 端：TCP 长连接
   └─► 发送 SignInInput 完成长连接登录

3. 同步离线消息
   └─► 调用 logic.MessageExtService.Sync

4. 收发消息
   └─► 发送：调用 MessageIntService.PushToUsers
   └─► 接收：通过长连接实时接收
```

### 消息同步机制

每个用户维护一个自增序列号（seq），用于消息同步：

1. **发送消息时**：获取发送者的 seq 并递增，将消息持久化到发送者的消息列表
2. **离线同步时**：客户端携带本地最大 seq，服务端返回大于该 seq 的所有消息

## Web 客户端

仓库内置一个Web IM 客户端，位于 `web/` 目录，基于 Vue 3 + TypeScript + Vite + Naive UI，通过 gRPC-Web 与 WebSocket 与后端通信，使用 IndexedDB 做本地消息持久化。

### 主要功能

- 手机号 + 验证码登录
- 会话列表、好友列表
- 私聊、群聊
- Markdown 消息，支持粘贴 / 拖拽 / 选择图片上传
- 明 / 暗主题切换

### 本地启动

```bash
cd web
npm install
npm run dev      # 启动开发服务，默认 http://127.0.0.1:5173
```

更详细的功能说明与协议交互见 [`web.md`](web.md)。

## API 接口

主要 Proto 文件位于 `pkg/protocol/proto/` 目录：

```
business/
├── user.ext.proto      # 用户外部接口（登录注册等）
├── user.int.proto      # 用户内部接口（鉴权）
├── friend.ext.proto    # 好友接口
└── message.ext.proto   # 消息接口

logic/
├── message.ext.proto   # 消息外部接口
├── message.int.proto   # 消息内部接口
├── device.int.proto    # 设备接口
├── group.int.proto     # 群组接口
└── room.int.proto      # 房间接口

connect/
├── connect.ext.proto   # 连接外部接口
└── connect.int.proto   # 连接内部接口
```

## 许可证

本项目基于 [MIT License](LICENSE) 开源。

## 链接

- GitHub: https://github.com/alberliu/gim