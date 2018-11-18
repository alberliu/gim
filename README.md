### 1 简要介绍

goim是一个即时通讯服务器，代码全部使用golang完成，功能包含好友之间一对一聊天，群组聊天，支持单用户多设备同时在线，就像微信一样，当你同时使用两个设备登录账号时，两个设备可以都可以接收到消息，当你用一个设备发送消息时，另一个设备也能收到你发送的消息。目前完成了第一版，第一版不想做的太复杂庞大，但是好多细节逻辑都做了反复的推敲，其主要目的是先作出核心功能，不考略加缓存和MQ提高性能，所以还不是很完善，以后会逐渐完善。
### 2 所用技术
golang+mysql完成,web框架使用了gin(对gin进行了简单的封装)，日志框架使用了zap,当然也自己写了一些小组件，例如TCP拆包粘包，唯一消息id生成器，数据库统一事务管理等。
### 3 项目分层设计
项目主要氛围两层，connect层和logic层，public包下放置了一些connect层和logic层公用的代码和一些基础库。  
connect连接层，主要维护和客户端的tcp连接，所以connect是有状态的，connect包含了TCP拆包粘包，消息解码，客户端心跳等一些逻辑。    
logic逻辑层是无状态的，主要维护消息的转发逻辑，以及对外提供http接口，提供一些聊天系统基本的业务功能，例如，登录，注册，添加好友，删除好友，创建群组，添加人到群组，群组踢人等功能
### 4 拆包粘包以及消息协议
TCP拆包粘包是自己写的一个算法，其思想就是每次从系统缓冲读取数据流，放置到自己实现的一个buffer中，以后拆包粘包，还是消息解码都是在这个buffer完成，其目的是减少内存拷贝，提高性能。  
其中每一个TCP都遵循TLV格式（即类型，长度，值），第一部分由两个字节来标示数据类型，第二部分用两个字节来标示数据长度，第三部分则是真正要解码的数据。  
消息协议使用Google的Protocol Buffers,具体消息协议定制在/public/proto包下  

### 5 消息唯一id
唯一消息id的主要作用是用来标示一次消息发送的完整流程，消息发送->消息投递->消息投递回执，用来线上排查线上问题。  
每一个消息有唯一的消息的id，由于消息发送频率比较高s，所以性能就很重要，当时没有找到合适第三方库，所以就自己实现了一个，原理就是，每次从数据库中拿一个数据段，用完了再去数据库拿，当用完之后去从数据库拿的时候，会有一小会的阻塞，为了解决这个问题，就做成了异步的，就是保证内存中有n个可用的id，当id消耗掉小于n个时，就从数据库获取生成，当达到n个时，goroutine阻塞等待id被消耗，如此往复。

### 6 主要逻辑
client: 客户端  
connect:连接层  
logic:逻辑层  
mysql:存储层  

#### 登录
[![3496be2f9ee9d33e.jpg](http://www.wailian.work/images/2018/11/12/3496be2f9ee9d33e.jpg)](http://www.wailian.work/image/BVGV24)

#### 单发
[![00d7e21cccc9050e.jpg](http://www.wailian.work/images/2018/11/12/00d7e21cccc9050e.jpg)](http://www.wailian.work/image/BVGZkp)
#### 群发
[![7ee3ada2baf1dec0.jpg](http://www.wailian.work/images/2018/11/12/7ee3ada2baf1dec0.jpg)](http://www.wailian.work/image/BVGtLc)
### 7 日志
使用了zap的日志框架，下图展示了一次两个设备从登录，发一条消息，再到下线的一次流程的完整日志
![9f644dcd04b20287.jpg](http://www.wailian.work/images/2018/11/12/9f644dcd04b20287.jpg)
### 8 api文档
https://documenter.getpostman.com/view/4164957/RzZ4q2hJ
### 9 github
https://github.com/alberliu/goim