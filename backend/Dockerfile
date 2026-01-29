# =========================
# 构建阶段 / Build Stage
# =========================
FROM golang:1.23-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

# 安装构建时需要的工具
RUN apk add --no-cache git ca-certificates tzdata

# 先复制 go.mod / go.sum 方便利用缓存
COPY go.mod go.sum ./
RUN go mod download

# 再复制源码
COPY . .

# 可选：加入编译信息（版本号、构建时间等），这里先简单写
# 你项目如果有 main 包路径，可以在这里设 LDFLAGS
# ENV LDFLAGS="-s -w"
# RUN go build -ldflags="$LDFLAGS" -o sol-green .

# 构建二进制
RUN go build -o sol-green .

# =========================
# 运行阶段 / Runtime Stage
# =========================
FROM alpine:3.19

# 基本运行依赖
RUN apk --no-cache add ca-certificates tzdata

ENV TZ=Asia/Shanghai

WORKDIR /app

# 创建非 root 用户
RUN addgroup -S app && adduser -S app -G app
USER app

# 从构建阶段复制二进制
COPY --from=builder /app/sol-green .

EXPOSE 8080

CMD ["./sol-green"]