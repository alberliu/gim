## 功能描述

1. 登录（手机号 + 验证码，首次登录自动注册，无单独注册页）
2. 会话列表（显示未读消息数角标）
3. 好友列表、新朋友页（好友申请/同意/忽略）、添加好友
4. 私聊、群聊
5. 消息内容使用 Markdown 格式（UTF-8 字符串，发送和渲染均按 Markdown 处理）
6. 支持主题色设置（明/暗模式切换 + 预设颜色主题）

## 技术选型

- 构建工具：Vite 8.0.10
- 框架：Vue 3.5.33
- 开发语言：TypeScript 6.0.3
- UI 框架：Naive UI 2.44.1
- 服务端通信协议：gRPC-web，通信库：@connectrpc/connect-web 2.1.1

## 消息浏览器本地持久化

使用 IndexedDB 持久化以下数据：

1. 会话列表（从本地消息记录推断，有收发消息记录的联系人/群组即出现在列表中）
2. 消息列表

## 后端通信地址

### gRPC-web 通信地址

```
127.0.0.1:8080
```

### WebSocket 通信地址

```
ws://127.0.0.1:8003/ws
```
### 文件上传下载

```
# 文件上传
http://127.0.0.1:8081/upload
# 文件下载
http://127.0.0.1:8081/file
```

## 后端通信协议

所有 protobuf 通信协议在 `pkg/proto` 目录。

### 鉴权

登录后获得 `token`，后续所有 gRPC-web 请求在 Header 中携带：

```
token: <token>
```

## 后端待新增接口

1. **群组创建外部接口**：在 business 层新增 `GroupExtService`，至少包含创建群组（参数：群名称、初始成员列表）。
2. **文件上传接口**：新增 HTTP 文件上传接口，用于头像上传及聊天内嵌图片，返回可访问的文件 URL。

## 流程

### WebSocket 建立连接流程

1. 调用 `UserExtService.SignIn` 获取 `user_id`、`device_id`、`token`
2. 建立 WebSocket 连接：`ws://127.0.0.1:8003/ws?user_id=xxx&device_id=xxx&token=xxx`
3. WebSocket 消息使用 protobuf 二进制编码，帧结构为 `connect.Packet`
4. 使用本地最大 seq 同步离线消息（无本地记录时传 0）
5. 连接建立后每 **5 分钟**发送一次心跳包（`PC_HEARTBEAT`），服务端会原包回显

### 消息接收流程

1. 收到 `PC_MESSAGE` 推送后，将 `Packet.content` 反序列化为 `connect.Message`
2. 根据 `Message.command` 判断消息类型：
   - `1`（MC_USER_MESSAGE）：私聊消息，content 反序列化为 `UserMessagePush`
   - `2`（MC_GROUP_MESSAGE）：群聊消息，content 反序列化为 `GroupMessagePush`
   - `110`（ADD_FRIEND）：好友申请推送，content 反序列化为 `AddFriendPush`
   - `111`（AGREE_ADD_FRIEND）：同意好友推送，content 反序列化为 `AgreeAddFriendPush`
   - `120`（UPDATE_GROUP）：群组更新推送，content 反序列化为 `UpdateGroupPush`
3. 使用本地最大 seq 调用 `MessageExtService.Sync` 同步消息（`has_more` 为 true 时循环拉取）
4. Sync 拿到消息后，调用 `MessageACKRequest` 向服务端确认（`device_ack` 为本次同步的最大 seq）

### 好友申请流程

1. 搜索用户（`UserExtService.SearchUser`，支持手机号或昵称搜索），发送好友申请（`FriendExtService.Add`，可填写备注和描述）
2. 收到 `AddFriendPush` 推送时，在"新朋友"页面显示申请记录（好友列表入口显示红点）
3. 对方在"新朋友"页面点击"同意"（可填备注），调用 `FriendExtService.Agree`；或点击"忽略"（仅前端隐藏，不调接口）
4. 同意后收到 `AgreeAddFriendPush` 推送，双方好友关系建立
5. 所有申请均已处理（同意或忽略）后，好友列表入口红点消失

### 创建群组流程

1. 从好友列表多选成员，填写群名称
2. 调用后端 `GroupExtService.Create` 创建群组（设备类型使用 `DT_WEB = 5`）
3. 创建成功后直接跳转到该群聊会话

## 功能细节

### 整体布局

三栏结构（参考微信 PC 版）：
- **左侧**：导航图标栏（会话、好友、设置等入口）
- **中间**：列表栏（会话列表 / 好友列表 / 新朋友等）
- **右侧**：内容区（聊天窗口 / 设置页等）

### 会话列表

- 从 IndexedDB 本地消息记录推断，有收发记录的联系人或群组自动出现
- 按最新消息时间降序排列
- 每个会话显示最后一条消息摘要和时间
- 显示未读消息数角标（具体数字）
- 进入会话时未读数立即清零

### 消息时间戳

- 相邻两条消息时间差超过 **5 分钟**时，在消息之间插入时间戳
- 当天显示"时:分"，跨天显示"月/日 时:分"

### 消息输入框

- 所见即所得（WYSIWYG）Markdown 编辑器
- **Enter** 发送消息，**Shift+Enter** 换行
- 支持图片插入（三种方式均需支持）：
  1. Ctrl+V 粘贴截图或复制的图片
  2. 点击工具栏"插入图片"按钮，弹出文件选择器
  3. 拖拽图片文件到输入框
- 图片插入后自动调用 `/upload` 接口上传，成功后将占位符替换为服务端返回的图片 URL，嵌入 Markdown

### 消息气泡

- 自己发送的消息靠右显示，对方消息靠左显示
- 群聊中每条消息气泡上方显示发送者昵称（自己的消息不显示昵称）
- 点击任意用户的头像或昵称，弹出该用户的资料卡片（含昵称、头像、添加好友按钮等）

### 消息历史

- 打开聊天窗口只加载 IndexedDB 本地记录，不主动向服务端拉取历史
- 消息列表展示本地已有记录，离线期间的消息通过 WebSocket 连接后 Sync 补全

### 头像

- 支持上传本地图片作为头像，调用 `/upload` 接口，获取 URL 后通过 `UpdateUser` 更新 `avatar_url`
- 群聊头像：自动生成九宫格（由成员头像拼合，类似微信）

### 用户资料

- 个人设置页可修改：头像、昵称
- 登录设备类型固定为 `DT_WEB`（枚举值 5）

### 主题设置

- 明/暗模式切换

### 好友功能

- 添加好友时可填写备注（remarks）和描述（description）
- 同意好友申请时也可填写备注（remarks）

### 群聊

- 支持创建群组（从好友列表多选成员 + 填写群名称）
- 支持收发群消息
- 不支持群成员管理（添加/移除成员、修改群信息等）

### 验证码

- 测试阶段使用固定验证码：`000000`

## 验收

可以使用 chrome-devtools 工具调试测试，保证功能可用。