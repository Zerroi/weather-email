# 构建阶段
FROM golang:1.22.12 AS builder
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# 关键修改：禁用 CGO + 静态编译
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o gin-app .

# 运行阶段
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/gin-app .
# 添加执行权限（可选）
RUN chmod +x gin-app
CMD ["./gin-app"]
