# sre-copilot

## 开发

#### 1. Git 克隆

```shell
git clone git@github.com:tang95/sre-copilot.git

cd sre-copilot
```

#### 2. 前端

注意后端请求转发配置[.umirc.ts](console/.umirc.ts)，详细配置参考 [Umi.js](https://umijs.org)

```shell
cd console
yarn
yarn start
```

访问 http://localhost:8000

#### 3. 后端

前提条件

1. 需要 Go 环境，推荐 go 1.23+
2. 生成自己的配置文件，参考 [server.yaml](conf/server.yaml)

```shell
cd sre-copilot
go mod download
# 默认读取 $HOME/.sre-copilot/server.yaml，建议放在此处。
go run ./cmd server -c server.yaml
```
