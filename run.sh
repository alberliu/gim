# 拉取远程代码
git pull

cd cmd/user
rm user
go build -o user main.go
pkill user
nohup ./user &

cd ../cmd/logic
rm logic
go build -o logic main.go
pkill logic
nohup ./logic &

cd ../cmd/tcp_conn
rm tcp_conn
go build -o tcp_conn main.go
pkill tcp_conn
nohup ./tcp_conn &

