export tag=v1.0
root:
	export ROOT=https://gitee.com/devops87/geektime.git

build:
	echo "building httpserver binary"
	mkdir -p bin/amd64
	GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -installsuffix cgo -o bin/amd64 .
	# CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64 .

# 通过 dickerfile 中的分阶段进行编译，未使用 build
release:
	echo "building httpserver container"
	sudo docker build -t devops87/httpserver:${tag} .
	# 删除在build镜像的过程中，可能会产生一些临时的不具有名称也没有作用的镜像
	sudo docker image rm $(sudo docker image ls -f dangling=true -q)

push: release
	echo "pushing devops87/httpserver"
	docker push devops87/httpserver:v1.0
