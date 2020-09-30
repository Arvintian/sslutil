# SSLUtil

SSL证书工具

## Install

```
go get -u github.com/Arvintian/sslutil
```

## Usage

```
#./bin/sslutil

Usage: sslutil -action [-ca] [-ca-key] -cfg -out -prefix
Options:
  -action string
        the action ca or sign
  -ca string
        ca pem
  -ca-key string
        ca key pem
  -cfg string
        config json file
  -out string
        cert and cert-key output dir, default current dir (default "workdir/sslutil")
  -prefix string
        cert and cert-key filename prefix (default "ca")
```

## 参数说明

### action

- ca 生成根证书
- sign 生成签名证书

### ca

- 生成签名证书必填,根证书cert pem

### ca-key

- 生成签名证书必填,根证书key pem

### cfg

- 证书配置json

### out

- 输出目录

### prefix

- 证书文件前缀


## Build

```
git clone https://github.com/Arvintian/sslutil.git
cd sslutil
make build
```

## Example

### 生成根证书

```
./bin/sslutil -action ca -cfg test/ca.json -out pki -prefix ca
cat test/ca.json
{
    "country": [
        "CN"
    ],
    "organization": [
        "k8s"
    ],
    "organizationalUnit": [
        "system-root"
    ],
    "province": [
        "BeiJing"
    ],
    "locality": [
        "BeiJing"
    ],
    "hosts": [
        ""
    ],
    "years": 10
}
```

### 生成签名证书

```
./bin/sslutil -action sign -ca pki/ca-cert.pem -ca-key pki/ca-key.pem -cfg test/server.json -out pki -prefix server
cat test/server.json
{
    "country": [
        "CN"
    ],
    "organization": [
        "k8s"
    ],
    "organizationalUnit": [
        "system-server"
    ],
    "province": [
        "BeiJing"
    ],
    "locality": [
        "BeiJing"
    ],
    "hosts": [
        "127.0.0.1",
        "localhost"
    ],
    "years": 10
}
```