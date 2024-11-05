# 构建后端项目
FROM golang:1.23 AS back
WORKDIR /app
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o sre-copilot .

# 运行环境
FROM ubuntu:latest
WORKDIR /app
COPY --from=back /app/sre-copilot .
COPY conf/server.yaml config.yaml
# 安装CA证书
RUN apt-get update && apt-get install -y ca-certificates

EXPOSE 8080
CMD ["./sre-copilot", "server", "-c", "config.yaml"]
