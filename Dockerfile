# 使用 Golang 作为构建环境
FROM golang:1.20 AS builder

# 设置工作目录
WORKDIR /app

# 复制项目文件并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源码并编译
COPY . .
RUN go build -o gin-app .

# 轻量级运行环境
FROM alpine:latest
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /root/
COPY --from=builder /app/gin-app .

# 运行应用
CMD ["./gin-app"]
