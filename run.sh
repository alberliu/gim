cd cmd/business
rm -f business
go build -o business main.go
echo "打包business成功"
pkill business
echo "停止business服务"
nohup ./business &
echo "启动business服务"

cd ../logic
rm -f logic
go build -o logic main.go
echo "打包logic成功"
pkill logic
echo "停止logic服务"
nohup ./logic &
echo "启动logic服务"

cd ../connect
rm -f connect
go build -o connect main.go
echo "打包connect成功"
pkill connect
echo "停止connect服务"
sleep 2
nohup ./connect &
echo "启动connect服务"

cd ../file
rm -f file
go build -o file main.go
echo "打包file成功"
pkill file
echo "停止file服务"
nohup ./file &
echo "启动file服务"

