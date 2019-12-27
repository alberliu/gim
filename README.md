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
├─ app           # 服务启动入口
│   ├── tcp_conn # TCP连接层启动入口
|   ├── ws_conn  # WebSocket连接层启动入口
│   └── logic    # 逻辑层启动入口
├─ conf # 配置
├─ tcp_conn      # TCP连接层服务代码
├─ ws_conn       # WebSocket连接层服务代码
├─ ligic         # 逻辑层服务代码
│   ├── cache    # 缓存，操作redis封装
|   ├── dao      # 操作数据库的封装
│   └── db       # Redis以及MySQL实例
│   └── model    # 模型层
│   └── service  # 业务逻辑层
│   └── rpc      # 对外提供的rpc接口
├─ public        # 连接层和逻辑层公共代码
├─ sql           # 数据库建表语句
├─ test          # 测试脚本
├─ docs          # 项目文档
```
### TCP拆包粘包
遵循LV的协议格式，一个消息包分为两部分，消息字节长度以及消息内容。
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
**特点：**  
消息存多份，群组中的每个成员保存一份  
**优缺点：**：  
1.支持离线消息同步      
2.群组成员越多，DB压力越大  
#### 超大群组：  
**特点：**  
消息只保存一份  
**优缺点：**  
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
### 错误处理,链路追踪,日志打印
   系统中的错误一般可以归类为两种，一种是业务定义的错误，一种就是未知的错误，在业务正式上线的时候，业务定义的错误的属于正常业务逻辑，不需要打印出来，
但是未知的错误，我们就需要打印出来，我们不仅要知道是什么错误，还要知道错误的调用堆栈，所以这里我对GRPC的错误进行了一些封装，使之包含调用堆栈。
```
func WrapError(err error) error {
	if err == nil {
		return nil
	}

	s := &spb.Status{
		Code:    int32(codes.Unknown),
		Message: err.Error(),
		Details: []*any.Any{
			{
				TypeUrl: TypeUrlStack,
				Value:   util.Str2bytes(stack()),
			},
		},
	}
	return status.FromProto(s).Err()
}
// Stack 获取堆栈信息
func stack() string {
	var pc = make([]uintptr, 20)
	n := runtime.Callers(3, pc)

	var build strings.Builder
	for i := 0; i < n; i++ {
		f := runtime.FuncForPC(pc[i] - 1)
		file, line := f.FileLine(pc[i] - 1)
		n := strings.Index(file, name)
		if n != -1 {
			s := fmt.Sprintf(" %s:%d \n", file[n:], line)
			build.WriteString(s)
		}
	}
	return build.String()
}
```
这样，不仅可以拿到错误的堆栈，错误的堆栈也可以跨RPC传输，但是，但是这样你只能拿到当前服务的堆栈，却不能拿到调用方的堆栈，就比如说，A服务调用
B服务，当B服务发生错误时，在A服务通过日志打印错误的时候，我们只打印了B服务的调用堆栈，怎样可以把A服务的堆栈打印出来。我们在A服务调用的地方也获取
一次堆栈。
```
func WrapRPCError(err error) error {
	if err == nil {
		return nil
	}
	e, _ := status.FromError(err)
	s := &spb.Status{
		Code:    int32(e.Code()),
		Message: e.Message(),
		Details: []*any.Any{
			{
				TypeUrl: TypeUrlStack,
				Value:   util.Str2bytes(GetErrorStack(e) + " --grpc-- \n" + stack()),
			},
		},
	}
	return status.FromProto(s).Err()
}

func interceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	err := invoker(ctx, method, req, reply, cc, opts...)
	return gerrors.WrapRPCError(err)
}

var LogicIntClient   pb.LogicIntClient

func InitLogicIntClient(addr string) {
	conn, err := grpc.DialContext(context.TODO(), addr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(interceptor))
	if err != nil {
		logger.Sugar.Error(err)
		panic(err)
	}

	LogicIntClient = pb.NewLogicIntClient(conn)
}
```
像这样，就可以获取完整一次调用堆栈。
错误打印也没有必要在函数返回错误的时候，每次都去打印。因为错误已经包含了堆栈信息
```
// 错误的方式
if err != nil {
	logger.Sugar.Error(err)
	return err
}

// 正确的方式
if err != nil {
	return err
}
```
然后，我们在上层统一打印就可以
```
func startServer {
    extListen, err := net.Listen("tcp", conf.LogicConf.ClientRPCExtListenAddr)
    if err != nil {
    	panic(err)
    }
	extServer := grpc.NewServer(grpc.UnaryInterceptor(LogicClientExtInterceptor))
	pb.RegisterLogicClientExtServer(extServer, &LogicClientExtServer{})
	err = extServer.Serve(extListen)
	if err != nil {
		panic(err)
	}
}

func LogicClientExtInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		logPanic("logic_client_ext_interceptor", ctx, req, info, &err)
	}()

	resp, err = handler(ctx, req)
	logger.Logger.Debug("logic_client_ext_interceptor", zap.Any("info", info), zap.Any("ctx", ctx), zap.Any("req", req),
		zap.Any("resp", resp), zap.Error(err))

	s, _ := status.FromError(err)
	if s.Code() != 0 && s.Code() < 1000 {
		md, _ := metadata.FromIncomingContext(ctx)
		logger.Logger.Error("logic_client_ext_interceptor", zap.String("method", info.FullMethod), zap.Any("md", md), zap.Any("req", req),
			zap.Any("resp", resp), zap.Error(err), zap.String("stack", gerrors.GetErrorStack(s)))
	}
	return
}
```
这样做的前提就是，在业务代码中透传context,golang不像其他语言，可以在线程本地保存变量，像Java的ThreadLocal,所以只能通过函数参数的形式进行传递，gim中，service层函数的第一个参数
都是context，但是dao层和cache层就不需要了，不然，显得代码臃肿。
### github
https://github.com/alberliu/gim
