### 编译服务
./build_docker.sh connect
./build_docker.sh logic
./build_docker.sh business
### 部署：
helm install gim ./chart
### 升级
增量升级：
helm upgrade gim ./chart --reuse-values --set server.$1.image=$image_name
全量升级：
helm upgrade gim ./chart

### 流量转发
转发到pod 
kubectl port-forward b-deployment-5c845465f9-4w4xv 30080:80  前面是宿主机端口，后面是容器端口
转发到service 
kubectl port-forward service/b 30080:80