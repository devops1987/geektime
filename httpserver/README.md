# 模块三作业
```shell
cd /opt/git-repo/geektime/httpserver
## 构建镜像
docker build -t httpserver:v2.0 .

## 推送镜像到dockerhub
docker push registry/httpserver:v2.0

## 启动容器
docker run -d -p 80:80 httpserver:v2.0

## 获取启动容器的 pid
docker inspect ef7d30cf7f4c|grep -i pid

## 进入network namespace
nsenter -n -t 9332
```


# 模块八作业

## 通过 make 进行构建管理
```shell
cd /opt/git-repo/geektime/httpserver
bash build.sh
```

## 部署 nginx-ingress-controller
```shell
kubectl apply -f ingress/nginx-ingress-controller.yaml
```

## 部署 httpserver 相关 namespace、deployment、service、ingress
```shell
kubectl apply -f deployment/
```

## 查看部署后相关信息
```shell
# k -n ingress-nginx get all
NAME                                 READY   STATUS    RESTARTS   AGE
pod/nginx-ingress-controller-ksbq5   1/1     Running   0          47m

NAME                                      DESIRED   CURRENT   READY   UP-TO-DATE   AVAILABLE   NODE SELECTOR       AGE
daemonset.apps/nginx-ingress-controller   1         1         1       1            1           ingress-node=lb01   3h7m

# k -n httpserver get all
NAME                             READY   STATUS    RESTARTS   AGE
pod/httpserver-dcf44d4d7-vzh85   1/1     Running   0          31m

NAME                 TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)   AGE
service/httpserver   ClusterIP   172.21.30.246   <none>        80/TCP    31m

NAME                         READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/httpserver   1/1     1            1           31m

NAME                                   DESIRED   CURRENT   READY   AGE
replicaset.apps/httpserver-dcf44d4d7   1         1         1       31m

# k -n httpserver get ingress
NAME                 HOSTS                    ADDRESS         PORTS   AGE
httpserver-ingress   ingress.httpserver.org   10.253.129.75   80      18m
```

## 验证请求
```shell
# curl ingress.httpserver.org
=================== Details of the http response header ============
X-Scheme=http
User-Agent=curl/7.58.0
X-Request-Id=9139576c0fc9651caea7cb5bba5eaea0
X-Forwarded-Host=ingress.httpserver.org
Version=null
X-Forwarded-Proto=http
X-Original-Uri=/
X-Real-Ip=10.253.129.54
X-Forwarded-Port=80
X-Forwarded-For=10.253.129.54
Accept=*/*

# curl ingress.httpserver.org/healthz
200
```

## 配置 https
```shell
## 自签证书
# openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/C=CN/ST=Chongqing/L=Chongqing/O=devops/CN=ingress.httpserver.org"
## 创建secret
# k -n httpserver create secret tls httpserver --cert=tls.crt --key=tls.key
## 编辑 ingress，支持 tls
# k apply -f ingress-httpserver.yaml

## 确认监听 443
# k -n httpserver get ingress        
NAME                 HOSTS                    ADDRESS         PORTS     AGE
httpserver-ingress   ingress.httpserver.org   10.253.129.75   80, 443   48m

## 验证
# curl -H "Host: ingress.httpserver.org" https://10.253.129.75 -v -k
* Rebuilt URL to: https://10.253.129.75/
*   Trying 10.253.129.75...
* TCP_NODELAY set
* Connected to 10.253.129.75 (10.253.129.75) port 443 (#0)
* ALPN, offering h2
* ALPN, offering http/1.1
* successfully set certificate verify locations:
*   CAfile: /etc/ssl/certs/ca-certificates.crt
  CApath: /etc/ssl/certs
* (304) (OUT), TLS handshake, Client hello (1):
* (304) (IN), TLS handshake, Server hello (2):
* TLSv1.2 (IN), TLS handshake, Certificate (11):
* TLSv1.2 (IN), TLS handshake, Server key exchange (12):
* TLSv1.2 (IN), TLS handshake, Server finished (14):
* TLSv1.2 (OUT), TLS handshake, Client key exchange (16):
* TLSv1.2 (OUT), TLS change cipher, Client hello (1):
* TLSv1.2 (OUT), TLS handshake, Finished (20):
* TLSv1.2 (IN), TLS handshake, Finished (20):
* SSL connection using TLSv1.2 / ECDHE-RSA-AES256-GCM-SHA384
* ALPN, server accepted to use h2
* Server certificate:
*  subject: O=Acme Co; CN=Kubernetes Ingress Controller Fake Certificate
*  start date: Nov 29 05:48:55 2021 GMT
*  expire date: Nov 29 05:48:55 2022 GMT
*  issuer: O=Acme Co; CN=Kubernetes Ingress Controller Fake Certificate
*  SSL certificate verify result: unable to get local issuer certificate (20), continuing anyway.
* Using HTTP2, server supports multi-use
* Connection state changed (HTTP/2 confirmed)
* Copying HTTP/2 data in stream buffer to connection buffer after upgrade: len=0
* Using Stream ID: 1 (easy handle 0x55ea157f2580)
> GET / HTTP/2
> Host: ingress.httpserver.org
> User-Agent: curl/7.58.0
> Accept: */*
> 
* Connection state changed (MAX_CONCURRENT_STREAMS updated)!
< HTTP/2 200 
< server: nginx/1.15.10
< date: Mon, 29 Nov 2021 07:08:32 GMT
< content-type: text/plain; charset=utf-8
< content-length: 325
< vary: Accept-Encoding
< accept: */*
< user-agent: curl/7.58.0
< version: null
< x-forwarded-for: 127.0.0.1
< x-forwarded-host: ingress.httpserver.org
< x-forwarded-port: 443
< x-forwarded-proto: https
< x-original-uri: /
< x-real-ip: 127.0.0.1
< x-request-id: e115872eb9bcaa6d284b263cf1ed264e
< x-scheme: https
< 
=================== Details of the http response header ============
Version=null
X-Forwarded-Proto=https
User-Agent=curl/7.58.0
Accept=*/*
X-Forwarded-For=127.0.0.1
X-Forwarded-Host=ingress.httpserver.org
X-Scheme=https
X-Original-Uri=/
X-Request-Id=e115872eb9bcaa6d284b263cf1ed264e
X-Real-Ip=127.0.0.1
X-Forwarded-Port=443
* Connection #0 to host 10.253.129.75 left intact
```