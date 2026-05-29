#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$SCRIPT_DIR"

if [ ! -f .env ]; then
    cp .env.example .env
    echo "已创建 .env 文件，请修改配置后重新运行"
    exit 1
fi

source .env

MODE="${1:-install}"

case "$MODE" in
    install)
        echo "=== 安装依赖 ==="
        cd web
        npm install
        cd ..
        echo "依赖安装完成"
        ;;

    build)
        echo "=== 构建前端 ==="
        cd web
        npm run build
        cd ..

        echo "=== 构建Go服务端 ==="
        cd server
        go build -o ../bin/consumable-tracker .
        cd ..
        echo "构建完成"
        ;;

    docker)
        echo "=== Docker部署（单容器集成数据库）==="
        docker compose build
        docker compose up -d
        echo "Docker部署完成"
        ;;

    init-db)
        echo "=== 初始化数据库（首次启动自动执行，一般无需手动）==="
        docker compose exec app psql -U admin -d consumable_db < database/init.sql 2>/dev/null || true
        echo "完成"
        ;;

    backup)
        echo "=== 手动备份 ==="
        docker compose exec app pg_dump -U admin consumable_db | gzip > "backups/backup_$(date +%Y%m%d_%H%M%S).sql.gz"
        echo "备份完成"
        ;;

    restore)
        BACKUP_FILE="$2"
        if [ -z "$BACKUP_FILE" ]; then
            echo "用法: $0 restore <备份文件名>"
            exit 1
        fi
        gunzip -c "$BACKUP_FILE" | docker compose exec -T app psql -U admin consumable_db
        echo "恢复完成"
        ;;

    logs)
        docker compose logs -f app
        ;;

    cert)
        if [ "$HTTPS_MODE" = "selfsigned" ]; then
            openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
                -keyout certs/privkey.pem -out certs/fullchain.pem \
                -subj "/CN=${DOMAIN:-localhost}"
            echo "自签名证书已生成"
        fi
        ;;

    *)
        echo "用法: $0 {install|build|docker|init-db|backup|restore|logs|cert}"
        echo ""
        echo "  install   安装前端依赖"
        echo "  build     构建前端和后端"
        echo "  docker    Docker编排部署"
        echo "  init-db   初始化数据库表结构"
        echo "  backup    手动备份数据库"
        echo "  restore   从备份文件恢复数据库"
        echo "  logs      查看API日志"
        echo "  cert      生成自签名证书"
        ;;
esac
