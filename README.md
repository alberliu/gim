### 简要介绍
gim是一个即时通讯服务器，代码全部使用golang完成。主要特性  
1.支持tcp，websocket接入  
2.离线消息同步    
3.单用户多设备同时在线    
4.单聊，群聊，以及房间聊天场景  
5.支持服务水平扩展  
6.使用领域驱动设计  
7.支持裸机部署和k8s部署  
gim可以作为以业务服务器的一个组件，为现有业务服务器提供im的能力，业务服务器
只需要实现user.int.proto协议中定义的GRPC接口，为gim服务提供基本的用户功能即可
### 使用技术：
数据库：MySQL+Redis  
通讯框架：GRPC  
长连接通讯协议：Protocol Buffers  
ORM框架：GORM
### 安装部署
#### docker compose部署
直接执行脚本deploy_compose.sh，即可部署一个单机集群
#### k8s部署
直接执行脚本deploy_k8s.sh   
### 项目目录简介
项目结构遵循 https://github.com/golang-standards/project-layout
```
cmd:          服务启动入口
config:       服务配置
deploy        部署配置文件
internal:     服务私有代码
pkg:          服务共有代码
sql:          项目sql文件
```
### 服务简介
1.connect  
维持与客户端的TCP和WebSocket长连接，心跳，以及TCP拆包粘包，消息编解码   
2.logic  
设备信息，好友信息，群组信息管理，消息转发逻辑  
3.user  
一个简单的用户服务，可以根据自己的业务需求，进行扩展,但是前提是，你的业务服务器实现了user.int.proto接口
### 客户端接入流程
1.调用RegisterDevice接口，完成设备注册，获取设备ID（device_id）,注意，一个设备只需完成一次注册即可，后续如果本地有device_id,就不需要注册了，举个例子，如果是APP第一次安装，就需要调用这个接口，后面即便是换账号登录，也不需要重新注册。  
2.调用SignIn接口，完成账户登录，获取账户登录的token。  
3.建立长连接，使用步骤2拿到的token，完成长连接登录。  
如果是web端,需要调用建立WebSocket时,如果是APP端，就需要建立TCP长连接。  
在完成建立TCP长连接时，第一个包应该是长连接登录包（SignInInput），如果信息无误，客户端就会成功建立长连接。  
4.使用长连接发送消息同步包（SyncInput），完成离线消息同步，注意：seq字段是客户端接收到消息的最大同步序列号，如果用户是换设备登录或者第一次登录，seq应该传0。  
接下来，用户可以使用LogicExt.SendMessage接口来发送消息，消息接收方可以使用长连接接收到对应的消息。
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
#### 群组：
采用写扩散，群组成员信息持久化到数据库保存。支持消息离线同步。  
#### 房间：  
采用读扩散，会将消息短暂的保存到Redis，长连接登录消息同步不会同步离线消息。
### 核心流程时序图
#### 长连接登录
![登录.png](https://upload-images.jianshu.io/upload_images/5760439-2e54d3c5dd0a44c1.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)
#### 离线消息同步
![离线消息同步.png](https://upload-images.jianshu.io/upload_images/5760439-aa513ea0de851e12.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)
#### 心跳
![心跳.png](https://upload-images.jianshu.io/upload_images/5760439-26d491374da3843b.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)
#### 消息单发
c1.d1和c1.d2分别表示c1用户的两个设备d1和d2,c2.d3和c2.d4同理
![消息单发.png](https://upload-images.jianshu.io/upload_images/5760439-35f1a91c8d7fffa6.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)
#### 群组消息群发
c1,c2.c3表示一个群组中的三个用户
![消息群发.png](https://upload-images.jianshu.io/upload_images/5760439-47a87c45b899b3f9.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)
### github
https://github.com/alberliu/gim
