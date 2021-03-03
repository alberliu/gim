### 简要介绍
im是一个即时通讯服务器，代码全部使用golang完成。主要功能  
1.支持tcp，websocket接入  
2.离线消息同步    
3.单用户多设备同时在线    
4.单聊，群聊，以及超大群聊天场景  
5.支持服务水平扩展  
gim和im有什么区别？im可以作为一个im中台提供给业务方使用，而gim可以作为以业务服务器的一个组件，为业务服务器提供im的能力，业务服务器的user服务只需要实现user.int.proto协议中定义的GRPC接口，为im服务提供基本的用户功能即可，其实以我目前的认知，我更推荐这种方式，这种模式相比于im,我认为最大好处在于一下两点：  
1.gim不需要考虑多个app的场景，相比im,业务复杂度降低了一个维度  
2.各个业务服务可以互不影响，可以做到风险隔离
### 使用技术：
数据库：MySQL+Redis  
通讯框架：GRPC  
长连接通讯协议：Protocol Buffers  
日志框架：Zap  
ORM框架：GORM
### 安装部署
1.首先安装MySQL，Redis  
2.创建数据库gim，执行sql/create_table.sql，完成初始化表的创建（数据库包含提供测试的一些初始数据）  
3.修改config下配置文件，使之和你本地配置一致  
4.分别切换到cmd的tcp_conn,ws_conn,logic,user目录下，执行go run main.go,启动TCP连接层服务器,WebSocket连接层服务器,逻辑层服务器,用户服务器  
（注意：tcp_conn只能在linux下启动，如果想在其他平台下启动，请安装docker，执行run.sh）  
### 项目目录简介
项目结构遵循 https://github.com/golang-standards/project-layout
```
cmd:          服务启动入口
config:       服务配置
internal:     每个服务私有代码
pkg:          服务共有代码
sql:          项目sql文件
test:         长连接测试脚本
```
### 服务简介
1.tcp_conn  
维持与客户端的TCP长连接，心跳，以及TCP拆包粘包，消息编解码  
2.ws_conn  
维持与客户端的WebSocket长连接，心跳，消息编解码  
3.logic  
设备信息，好友信息，群组信息管理，消息转发逻辑  
4.business  
一个简单的业务服务器服务，可以根据自己的业务需求，进行扩展,但是前提是，你的业务服务器实现了user.int.proto接口
### 客户端接入流程
1.调用LogicExt.RegisterDevice接口，完成设备注册，获取设备ID（device_id）,注意，一个设备只需完成一次注册即可，后续如果本地有device_id,就不需要注册了，举个例子，如果是APP第一次安装，就需要调用这个接口，后面即便是换账号登录，也不需要重新注册。  
2.调用UserExt.SignIn接口，完成账户登录，获取账户登录的token。  
3.建立长连接，使用步骤2拿到的token，完成长连接登录。  
如果是web端,需要调用建立WebSocket时,将user_id,device_id,token，以URL参数的形式传递到服务器，完成长连接登录，例如：ws://localhost:8081/ws?user_id={user_id}&device_id={device_id}&token={token}  
如果是APP端，就需要建立TCP长连接，在完成建立TCP长连接时，第一个包应该是长连接登录包（SignInInput），如果信息无误，客户端就会成功建立长连接。  
4.使用长连接发送消息同步包（SyncInput），完成离线消息同步，注意：seq字段是客户端接收到消息的最大同步序列号，如果用户是换设备登录或者第一次登录，seq应该传0。  
接下来，用户可以使用LogicExt.SendMessage接口来发送消息，消息接收方可以使用长连接接收到对应的消息。  
### 网络模型
TCP的网络层使用linux的epoll实现，相比golang原生，能减少goroutine使用，从而节省系统资源占用
### 单用户多设备支持，离线消息同步
每个用户都会维护一个自增的序列号，当用户A给用户B发送消息是，首先会获取A的最大序列号，设置为这条消息的seq，持久化到用户A的消息列表，
再通过长连接下发到用户A账号登录的所有设备，再获取用户B的最大序列号，设置为这条消息的seq，持久化到用户B的消息列表，再通过长连接下发
到用户B账号登录的所有设备。  
假如用户的某个设备不在线，在设备长连接登录时，用本地收到消息的最大序列号，到服务器做消息同步，这样就可以保证离线消息不丢失。
### 读扩散和写扩散
首先解释一下，什么是读扩散，什么是写扩散  
#### 读扩散
**简介**：群组成员发送消息时，先建立一个会话，都将这个消息写入这个会话中，同步离线消息时，需要同步这个会话的未同步消息  
**优点**：每个消息只需要写入数据库一次就行，减少数据库访问次数，节省数据库空间  
**缺点**：一个用户有n个群组，客户端每次同步消息时，要上传n个序列号，服务器要对这n个群组分别做消息同步  
#### 写扩散
**简介**：在群组中，每个用户维持一个自己的消息列表，当群组中有人发送消息时，给群组的每个用户的消息列表插入一条消息即可  
**优点**：每个用户只需要维护一个序列号和消息列表  
**缺点**：一个群组有多少人，就要插入多少条消息，当群组成员很多时，DB的压力会增大
### 消息转发逻辑选型以及特点
#### 普通群组：
采用写扩散，群组成员信息持久化到数据库保存。支持消息离线同步。  
#### 超大群组：  
采用读扩散，群组成员信息保存到redis,不支持离线消息同步。
### 核心流程时序图
#### 长连接登录
![eaf3a08af9c64bbd.png](http://www.wailian.work/images/2019/10/26/eaf3a08af9c64bbd.png)
#### 离线消息同步
![ef9c9452e65be3ced63573164fec7ed5.png](http://s1.wailian.download/2019/12/25/ef9c9452e65be3ced63573164fec7ed5.png)
#### 心跳
![6ea6acf2cd4b956e.png](http://www.wailian.work/images/2019/10/26/6ea6acf2cd4b956e.png)
#### 消息单发
c1.d1和c1.d2分别表示c1用户的两个设备d1和d2,c2.d3和c2.d4同理
![e000fda2f18e86f3.png](http://www.wailian.work/images/2019/10/26/e000fda2f18e86f3.png)
#### 小群消息群发
c1,c2.c3表示一个群组中的三个用户
![749fc468746055a8ecf3fba913b66885.png](http://s1.wailian.download/2019/12/26/749fc468746055a8ecf3fba913b66885.png)
#### 大群消息群发
![e3f92bdbb3eef199d185c28292307497.png](http://s1.wailian.download/2019/12/26/e3f92bdbb3eef199d185c28292307497.png)
#### APP
基于Flutter写了一个简单的客户端  
GitHub地址：https://github.com/alberliu/fim  
APP下载：https://github.com/alberliu/fim/releases/download/v1.0.0/FIM.apk  
APP截图：  
[![8a3937b2c4a53b568fea418ff6f0ee4a.md.png](https://p.130014.xyz/2021/02/07/8a3937b2c4a53b568fea418ff6f0ee4a.md.png)](http://www.wailian.work/image/QauQqB)
[![dd47b5027e4983a068ca91b2ef1b8b1c.md.png](https://p.130014.xyz/2021/02/07/dd47b5027e4983a068ca91b2ef1b8b1c.md.png)](http://www.wailian.work/image/Qau58S)
[![8be32a9ab20de799a2e2d53d5c23eddb.md.png](https://p.130014.xyz/2021/02/07/8be32a9ab20de799a2e2d53d5c23eddb.md.png)](http://www.wailian.work/image/QauBFI)
[![2fa211eb51022831c21da470a2d87330.md.png](https://p.130014.xyz/2021/02/07/2fa211eb51022831c21da470a2d87330.md.png)](http://www.wailian.work/image/QauMRb)
[![8d3b9dfa3b022c727c640e8938b34105.md.png](https://p.130014.xyz/2021/02/07/8d3b9dfa3b022c727c640e8938b34105.md.png)](http://www.wailian.work/image/Qan758)
[![63c839904f676c69892ef54946094684.md.png](https://p.130014.xyz/2021/02/07/63c839904f676c69892ef54946094684.md.png)](http://www.wailian.work/image/QanLaO)
### 联系方式
![my.png](https://p.130014.xyz/2021/03/03/my.png)
### 赞赏支持
如果觉得项目对你有帮助，请支持一下  
[![WechatIMG9.md.jpg](https://p.130014.xyz/2021/03/01/WechatIMG9.md.jpg)](https://www.wailian.work/image/Q90fg0)
### github
https://github.com/alberliu/gim
