#!/bin/bash
set -e

echo "=== 耗材使用管理系统 - 一键部署 ==="

# Check Docker
if ! command -v docker &> /dev/null; then
    echo "错误: 未安装 Docker，请先安装 Docker"
    echo "   brew install --cask docker"
    exit 1
fi

# Create .env if not exists
if [ ! -f .env ]; then
    echo "创建 .env 配置文件..."
    cp .env.example .env
    JWT_SECRET=$(openssl rand -hex 16 2>/dev/null || echo "change_this_to_a_random_32_char_string")
    if [[ "$OSTYPE" == "darwin"* ]]; then
        sed -i '' "s/change_this_to_a_random_32_char_string/$JWT_SECRET/" .env
    else
        sed -i "s/change_this_to_a_random_32_char_string/$JWT_SECRET/" .env
    fi
    echo "  已生成随机 JWT_SECRET"
fi

# Create directories
mkdir -p data backups certs

# Build and run
echo "构建 Docker 镜像..."
docker build -t consumable-tracker .

echo "启动容器..."
docker stop consumable-tracker 2>/dev/null || true
docker rm consumable-tracker 2>/dev/null || true

docker run -d \
    --name consumable-tracker \
    --restart unless-stopped \
    -p 8443:8443 \
    -v "$(pwd)/data:/app/data" \
    -v "$(pwd)/backups:/app/backups" \
    -v "$(pwd)/certs:/app/certs" \
    --env-file .env \
    consumable-tracker

echo ""
echo "=== 部署完成 ==="
echo "访问地址: http://localhost:8443"
echo "默认账号: admin / admin123"
echo ""
echo "查看日志: docker logs -f consumable-tracker"
echo "停止服务: docker stop consumable-tracker"
