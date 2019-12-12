打包命令
protoc --go_out ../../pb/pbtcp protocol.proto
protoc --go_out=plugins=grpc:../pb/ *.proto

TODO list
1.业务推送
2.消息逻辑整理

测试
error整理