#!/bin/bash

# Sol-Green 项目测试脚本

set -e

echo "🧪 开始运行测试..."

# 运行 Go 测试
echo "📝 运行 Go 后端测试..."
cd backend
go test -v ./tests/... -coverprofile=coverage.out
cd ..

# 显示测试覆盖率
if [ -f backend/coverage.out ]; then
    echo ""
    echo "📊 测试覆盖率:"
    cd backend
    go tool cover -func=coverage.out | tail -1
    rm coverage.out
    cd ..
fi

# 运行前端测试（如果存在）
if [ -d "frontend" ]; then
    echo ""
    echo "📝 运行前端测试..."
    cd frontend
    if [ -f "package.json" ]; then
        npm test -- --watchAll=false || true
    fi
    cd ..
fi

echo ""
echo "✅ 测试完成！"
