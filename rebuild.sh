helm delete gim

./build_docker.sh connect
./build_docker.sh logic
./build_docker.sh business

helm install gim ./chart