#!/bin/bash

# Sol-Green 项目部署脚本

set -e

echo "🚀 开始部署 Sol-Green 项目..."

# 检查环境
check_requirements() {
    echo "📋 检查部署环境..."
    
    if ! command -v docker &> /dev/null; then
        echo "❌ Docker 未安装，请先安装 Docker"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        echo "❌ Docker Compose 未安装，请先安装 Docker Compose"
        exit 1
    fi
    
    echo "✅ 环境检查通过"
}

# 加载环境变量
load_env() {
    if [ -f .env ]; then
        echo "📝 加载环境变量..."
        export $(cat .env | grep -v '^#' | xargs)
    else
        echo "⚠️  未找到 .env 文件，使用默认配置"
    fi
}

# 构建 Docker 镜像
build_images() {
    echo "🔨 构建 Docker 镜像..."
    docker-compose build
    echo "✅ 镜像构建完成"
}

# 启动服务
start_services() {
    echo "🚀 启动服务..."
    docker-compose up -d
    echo "✅ 服务启动完成"
}

# 等待服务就绪
wait_for_services() {
    echo "⏳ 等待服务就绪..."
    sleep 10
    
    # 检查后端健康状态
    max_attempts=30
    attempt=0
    while [ $attempt -lt $max_attempts ]; do
        if curl -f http://localhost:8080/api/v1/health &> /dev/null; then
            echo "✅ 后端服务已就绪"
            break
        fi
        attempt=$((attempt + 1))
        echo "等待后端服务... ($attempt/$max_attempts)"
        sleep 2
    done
    
    if [ $attempt -eq $max_attempts ]; then
        echo "❌ 后端服务启动超时"
        exit 1
    fi
}

# 运行数据库迁移
run_migrations() {
    echo "📊 运行数据库迁移..."
    # 数据库迁移已在应用启动时自动执行
    echo "✅ 数据库迁移完成"
}

# 显示服务状态
show_status() {
    echo ""
    echo "📊 服务状态:"
    docker-compose ps
    
    echo ""
    echo "🌐 服务地址:"
    echo "  后端 API: http://localhost:8080"
    echo "  前端应用: http://localhost:3000 (开发模式)"
    echo ""
    echo "✅ 部署完成！"
}

# 主函数
main() {
    check_requirements
    load_env
    build_images
    start_services
    wait_for_services
    run_migrations
    show_status
}

main
