# 构建：使用 golang:1.15 版本
FROM golang:1.15 as build
# 设置工作区
WORKDIR /go/release
# 把需要编译的文件添加到 /go/release 目录
ADD main.go .
# 编译：把cmd/main.go编译成可执行的二进制文件，命名为 httpserver
#RUN GO111MODULE=off GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -installsuffix cgo -o httpserver main.go
RUN GO111MODULE=on GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -installsuffix cgo -o httpserver main.go

# 运行：使用 ubuntu/alpine 作为基础镜像
# FROM ubuntu
FROM alpine:3.14
# 设置时区
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo 'Asia/Shanghai' > /etc/timezone
# 复制build阶段构建的可执行 go 二进制文件 
COPY --from=build /go/release/httpserver /
EXPOSE 80
ENTRYPOINT ["/httpserver"]
