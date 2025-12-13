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
|------|------|
| 语言 | Go 1.18+ |
| 数据库 | MySQL |
| 缓存 | Redis |
| 通讯框架 | gRPC |
| 传输协议 | Protocol Buffers |
| ORM | GORM |

## 快速开始

### Docker Compose 部署（推荐）

```bash
# 克隆项目
git clone https://github.com/alberliu/gim.git
cd gim

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
└── test/          # 测试代码
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
- 可根据业务需求自定义扩展

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
   └─► 参数 seq：客户端已收到消息的最大序列号
   └─► 首次登录或换设备登录时 seq 传 0

4. 收发消息
   └─► 发送：调用 MessageIntService.PushToUsers
   └─► 接收：通过长连接实时接收
```

### 消息同步机制

每个用户维护一个自增序列号（seq），用于消息同步：

1. **发送消息时**：获取发送者的 seq 并递增，将消息持久化到发送者的消息列表
2. **接收消息时**：获取接收者的 seq 并递增，将消息持久化到接收者的消息列表
3. **离线同步时**：客户端携带本地最大 seq，服务端返回大于该 seq 的所有消息

## 消息模型

### 读扩散 vs 写扩散

| 对比项 | 读扩散 | 写扩散 |
|--------|--------|--------|
| **原理** | 消息写入会话，成员各自同步 | 消息写入每个成员的消息列表 |
| **优点** | 写入次数少，节省存储空间 | 同步简单，只需维护一个 seq |
| **缺点** | 同步时需要处理多个会话 | 群成员多时写入压力大 |

### GIM 的选择

| 场景 | 模型 | 说明 |
|------|------|------|
| **群聊** | 写扩散 | 成员信息持久化，支持完整离线同步 |
| **房间** | 读扩散 | 消息暂存 Redis，不同步离线消息 |

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