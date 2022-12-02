helm delete gim -n gimns

./build_docker.sh connect
./build_docker.sh logic
./build_docker.sh business

helm install gim ./chart -n gimns