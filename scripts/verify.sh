#!/bin/bash

# Sol-Green 项目验证脚本

set -e

echo "🔍 开始验证 Sol-Green 项目..."

ERRORS=0

# 检查目录结构
check_structure() {
    echo ""
    echo "📁 检查项目结构..."
    
    REQUIRED_DIRS=(
        "backend"
        "backend/config"
        "backend/controller"
        "backend/middleware"
        "backend/model"
        "backend/service"
        "frontend"
        "frontend/src"
        "contracts/sol-green"
        "docs"
        "scripts"
    )
    
    for dir in "${REQUIRED_DIRS[@]}"; do
        if [ ! -d "$dir" ]; then
            echo "❌ 缺少目录: $dir"
            ERRORS=$((ERRORS + 1))
        else
            echo "✅ $dir"
        fi
    done
}

# 检查关键文件
check_files() {
    echo ""
    echo "📄 检查关键文件..."
    
    REQUIRED_FILES=(
        "backend/main.go"
        "backend/go.mod"
        "backend/Dockerfile"
        "frontend/package.json"
        "frontend/src/App.js"
        "contracts/sol-green/src/lib.rs"
        "contracts/sol-green/Cargo.toml"
        "docker-compose.yml"
        ".env.example"
    )
    
    for file in "${REQUIRED_FILES[@]}"; do
        if [ ! -f "$file" ]; then
            echo "❌ 缺少文件: $file"
            ERRORS=$((ERRORS + 1))
        else
            echo "✅ $file"
        fi
    done
}

# 验证 Go 代码语法
verify_go() {
    echo ""
    echo "🔍 验证 Go 代码..."
    
    if command -v go &> /dev/null; then
        cd backend
        echo "  检查 Go 模块..."
        if go mod verify &> /dev/null; then
            echo "✅ Go 模块验证通过"
        else
            echo "⚠️  Go 模块验证失败（可能需要运行 go mod download）"
        fi
        
        echo "  检查 Go 代码语法..."
        if go build -o /dev/null . &> /dev/null; then
            echo "✅ Go 代码编译通过"
        else
            echo "❌ Go 代码编译失败"
            go build . 2>&1 | head -10
            ERRORS=$((ERRORS + 1))
        fi
        cd ..
    else
        echo "⚠️  Go 未安装，跳过 Go 代码验证"
    fi
}

# 验证前端代码
verify_frontend() {
    echo ""
    echo "🔍 验证前端代码..."
    
    if [ -f "frontend/package.json" ]; then
        echo "✅ package.json 存在"
        
        if command -v node &> /dev/null; then
            cd frontend
            if [ -d "node_modules" ]; then
                echo "✅ node_modules 已安装"
            else
                echo "⚠️  node_modules 未安装（运行 npm install）"
            fi
            cd ..
        else
            echo "⚠️  Node.js 未安装，跳过前端验证"
        fi
    else
        echo "❌ 缺少 frontend/package.json"
        ERRORS=$((ERRORS + 1))
    fi
}

# 验证合约代码
verify_contracts() {
    echo ""
    echo "🔍 验证智能合约..."
    
    if [ -f "contracts/sol-green/Cargo.toml" ]; then
        echo "✅ Cargo.toml 存在"
        
        if [ -f "contracts/sol-green/src/lib.rs" ]; then
            echo "✅ lib.rs 存在"
        else
            echo "❌ 缺少 lib.rs"
            ERRORS=$((ERRORS + 1))
        fi
    else
        echo "❌ 缺少 Cargo.toml"
        ERRORS=$((ERRORS + 1))
    fi
}

# 验证 Docker 配置
verify_docker() {
    echo ""
    echo "🔍 验证 Docker 配置..."
    
    if [ -f "docker-compose.yml" ]; then
        echo "✅ docker-compose.yml 存在"
        
        if command -v docker-compose &> /dev/null; then
            echo "  验证 docker-compose 配置..."
            if docker-compose config &> /dev/null; then
                echo "✅ docker-compose 配置有效"
            else
                echo "❌ docker-compose 配置无效"
                docker-compose config 2>&1 | head -5
                ERRORS=$((ERRORS + 1))
            fi
        else
            echo "⚠️  docker-compose 未安装，跳过配置验证"
        fi
    else
        echo "❌ 缺少 docker-compose.yml"
        ERRORS=$((ERRORS + 1))
    fi
}

# 验证脚本
verify_scripts() {
    echo ""
    echo "🔍 验证脚本..."
    
    SCRIPTS=(
        "scripts/deploy.sh"
        "scripts/test.sh"
        "scripts/build-contract.sh"
    )
    
    for script in "${SCRIPTS[@]}"; do
        if [ -f "$script" ]; then
            if [ -x "$script" ]; then
                echo "✅ $script (可执行)"
            else
                echo "⚠️  $script (不可执行，运行 chmod +x $script)"
            fi
        else
            echo "❌ 缺少脚本: $script"
            ERRORS=$((ERRORS + 1))
        fi
    done
}

# 主函数
main() {
    check_structure
    check_files
    verify_go
    verify_frontend
    verify_contracts
    verify_docker
    verify_scripts
    
    echo ""
    echo "=========================================="
    if [ $ERRORS -eq 0 ]; then
        echo "✅ 项目验证通过！"
        echo ""
        echo "下一步："
        echo "  1. 配置环境变量: cp .env.example .env"
        echo "  2. 启动服务: docker-compose up -d"
        echo "  3. 或本地运行: cd backend && go run main.go"
        exit 0
    else
        echo "❌ 发现 $ERRORS 个问题，请修复后重试"
        exit 1
    fi
}

main
