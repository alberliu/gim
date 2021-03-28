CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
echo "打包完成"
docker run -v $(pwd)/:/app -p 8080:8080 alpine .//app/main
