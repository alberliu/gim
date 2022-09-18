if [[ $? -ne 0 ]]; then
    exit 1
fi

server=$1
cd cmd/$server
# 打包可执行文件
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main
pwd
mv main ../../docker/
cd ../../docker/
pwd
# 构建镜像
docker build -t $1 .

kind load docker-image $server --name kind
