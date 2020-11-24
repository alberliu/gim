# 拉取远程代码
git pull

cd cmd/user
rm -f user
go build -o user main.go
echo "打包user成功"
pkill user
echo "停止user服务"
nohup ./user &
echo "启动user服务"

cd ../cmd/logic
rm -f logic
go build -o logic main.go
echo "打包logic成功"
pkill logic
echo "停止logic服务"
nohup ./logic &
echo "启动logic服务"

cd ../cmd/tcp_conn
rm -f tcp_conn
go build -o tcp_conn main.go
echo "打包tcp_conn成功"
pkill tcp_conn
echo "停止tcp_conn服务"
nohup ./tcp_conn &
echo "启动tcp_conn服务"

