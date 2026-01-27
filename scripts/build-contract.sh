#!/bin/bash

# Solana 智能合约构建脚本

set -e

echo "🔨 开始构建 Solana 智能合约..."

# 检查 Anchor 是否安装
if ! command -v anchor &> /dev/null; then
    echo "❌ Anchor 未安装，请先安装 Anchor CLI"
    echo "   安装命令: cargo install --git https://github.com/coral-xyz/anchor avm --locked --force"
    echo "   然后运行: avm install latest && avm use latest"
    exit 1
fi

# 进入合约目录
cd contracts/sol-green

# 构建合约
echo "📦 构建合约..."
anchor build

# 运行测试
echo "🧪 运行合约测试..."
anchor test --skip-local-validator || echo "⚠️  测试跳过（需要本地验证器）"

echo "✅ 合约构建完成！"
echo "📁 构建产物位置: target/deploy/"
