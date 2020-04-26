CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
docker run -v $(pwd)/:/app -p 8085:8080 alpine .//app/main