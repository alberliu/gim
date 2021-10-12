cd pkg/proto
protoc --go_out=plugins=grpc:../../../ *.proto
cd ../../