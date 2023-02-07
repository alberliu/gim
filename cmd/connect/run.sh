CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
echo "打包完成"
docker run -v $(pwd)/:/app -p 8000:8000 -p 8002:8002 -p 8003:8003 alpine .//app/main
