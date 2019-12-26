### 简要介绍
gim是一个即时通讯服务器，代码全部使用golang完成。主要功能  
1.支持tcp，websocket接入  
2.离线消息同步  
3.多业务接入  
4.单用户多设备同时在线    
5.单聊，群聊，以及超大群聊天场景  
6.支持服务水平扩展
### 使用技术：
数据库：Mysql+Redis  
组件：grpc+jsoniter+zap  
### 安装部署
1.首先安装MySQL，Redis  
2.创建数据库gim，执行sql/create_table.sql，完成初始化表的创建（数据库包含提供测试的一些初始数据）   
3.修改conf下配置文件，使之和你本地配置一致  
4.分别切换到app的tcp_conn,ws_conn,logic目录下，执行go run main.go,启动TCP连接层服务器,WebSocket连接层服务器,逻辑层服务器  
### 业务服务器如何接入
1.首先生成私钥和公钥  
2.在app表里根据你的私钥添加一条app记录    
3.将app_id和公钥保存到业务服务器  
4.将用户通过LogicClientExtServer.AddUser接口添加到IM服务器  
5.通过LogicClientExtServer.RegisterDevice接口注册设备，获取设备id(device_id)  
6.将app_id，user_id,device_id用公钥通过公钥加密，生成token,相应库的代码在public/util/aes.go  
7.接下来使用这个token，app就可以和IM服务器交互
### rpc接口简介
项目所有的proto协议在gim/public/proto/目录下  
1.tcp.proto  
长连接通讯协议  
2.logic_client.ext.proto  
对客户端（Android设备，IOS设备）提供的rpc协议  
3.logic_server.ext.proto    
对业务服务器提供的rpc协议  
4.logic.int.proto  
对conn服务层提供的rpc协议  
5.conn.int.proto  
对logic服务层提供的rpc协议  
### 项目目录介绍
```bash
├─ app # 服务启动入口
│   ├── tcp_conn # TCP连接层启动入口
|   ├── ws_conn # WebSocket连接层启动入口
│   └── logic   # 逻辑层启动入口
├─ conf # 配置
├─ tcp_conn # TCP连接层服务代码
├─ ws_conn # WebSocket连接层服务代码
├─ ligic # 逻辑层服务代码
├─ public # 连接层和逻辑层公共代码
├─ sql # 数据库建表语句
├─ test # 测试脚本
├─ docs # 项目文档
```
### TCP拆包粘包
遵循TLV的协议格式，一个消息包分为三部分，消息类型（两个字节），消息包内容长度（两个字节），消息内容。  
这里为了减少内存分配，拆出来的包的内存复用读缓存区内存。  
**拆包流程：**  
1.首先从系统缓存区读取字节流到buffer  
2.根据包头的length字段，检查报的value字段的长度是否大于等于length  
3.如果大于，返回一个完整包（此包内存复用），重复步骤2  
4.如果小于，将buffer的有效字节前移，重复步骤1  
### 服务简介
1.tcp_conn  
维持与客户端的TCP长连接，心跳，以及TCP拆包粘包，消息编解码  
1.ws_conn  
维持与客户端的WebSocket长连接，心跳，消息编解码  
2.logic  
消息转发逻辑，设备信息，用户信息，群组信息的操作  
### 离线消息同步
用户的消息维护一个自增的序列号，当客户端TCP连接断开重新建立连接时，首先要做TCP长连接的登录，然后用客户端本地已经同步的最大的序列号做消息同步，这样就可以保证离线消息的不丢失。  
### 单用户多设备支持
当用户发送消息时，除了将消息发送目的用户  
在DB中，每个用户只维护一个自己的消息列表，但是用户的每个设备各自维护自己的同步序列号，设备使用自己的同步序列号在消息列表中做消息同步  
### 消息转发逻辑
单聊和普通群组采用写扩散，超级大群使用读扩散。  
读扩散和写扩散的选型。  
首先解释一下，什么是读扩散，什么是写扩散  
#### 读扩散
**简介**：群组成员发送消息时，也是先建立一个会话，都将这个消息写入这个会话中，同步离线消息时，需要同步这个会话的未同步消息  
**优点**：每个消息只需要写入数据库一次就行，减少数据库访问次数，节省数据库空间  
**缺点**：一个用户有n个群组，客户端每次同步消息时，要上传n个序列号，服务器要对这n个群组分别做消息同步  
#### 写扩散
**简介**：就是每个用户维持一个消息列表，当有其他用户给这个用户发送消息时，给这个用户的消息列表插入一条消息即可  
**优点**：每个用户只需要维护一个序列号和消息列表  
**缺点**：一个群组有多少人，就要插入多少条消息，当群组成员很多时，DB的压力会增大
### 群组简介
#### 普通群组：
1.支持离线消息同步    
2.群组成员越多，DB压力越大
#### 超大群组：
1.DB压力不会随着群组成员的人数的增加而增加  
2.不支持离线消息同步
### 核心流程时序图
#### 长连接登录
![eaf3a08af9c64bbd.png](http://www.wailian.work/images/2019/10/26/eaf3a08af9c64bbd.png)
#### 离线消息同步
![ef9c9452e65be3ced63573164fec7ed5.png](http://s1.wailian.download/2019/12/25/ef9c9452e65be3ced63573164fec7ed5.png)
#### 心跳
![6ea6acf2cd4b956e.png](http://www.wailian.work/images/2019/10/26/6ea6acf2cd4b956e.png)
#### 消息单发
![e000fda2f18e86f3.png](http://www.wailian.work/images/2019/10/26/e000fda2f18e86f3.png)
#### 小群消息群发
![749fc468746055a8ecf3fba913b66885.png](http://s1.wailian.download/2019/12/26/749fc468746055a8ecf3fba913b66885.png)
#### 大群消息群发
![e3f92bdbb3eef199d185c28292307497.png](http://s1.wailian.download/2019/12/26/e3f92bdbb3eef199d185c28292307497.png)
### github
https://github.com/alberliu/gim
