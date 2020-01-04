syntax = "proto3";
import "tcp.proto";
package pb;

message SignInReq {
    int64 app_id = 1; // app_id
    int64 device_id = 2; // 设备id
    int64 user_id = 3; // 用户id
    string token = 4; // 秘钥
    string conn_addr = 5; // 服务器地址
}
message SignInResp {
}

message SyncReq {
    int64 app_id = 2; // appId
    int64 user_id = 3; // 用户id
    int64 device_id = 4; // 设备id
    int64 seq = 5; // 客户端已经同步的序列号
}
message SyncResp {
    repeated MessageItem messages = 3; // 消息列表
}

message MessageACKReq {
    int64 app_id = 1; // appId
    int64 user_id = 2; // 用户id
    int64 device_id = 3; // 设备id
    string message_id = 4; // 消息id
    int64 device_ack = 5; // 设备收到消息的确认号
    int64 receive_time = 6; // 消息接收时间戳，精确到毫秒
}
message MessageACKResp {

}

message OfflineReq {
    int64 app_id = 2; // appId
    int64 user_id = 3; // 用户id
    int64 device_id = 4; // 设备id
}
message OfflineResp {

}


service LogicInt {
    //  登录
    rpc SignIn (SignInReq) returns (SignInResp);
    //  消息同步
    rpc Sync (SyncReq) returns (SyncResp);
    //  设备收到消息回执
    rpc MessageACK (MessageACKReq) returns (MessageACKResp);
    //  设备离线
    rpc Offline (OfflineReq) returns (OfflineResp);
}