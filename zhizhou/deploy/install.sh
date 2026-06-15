#!/bin/bash
# 知舟 - 一键部署脚本
set -e

echo "========================================="
echo "  知舟 - 一键部署"
echo "  个人知识管理工具"
echo "========================================="
echo ""

# 检查 Docker
if ! command -v docker &> /dev/null; then
    echo "错误: 请先安装 Docker"
    echo "curl -fsSL https://get.docker.com | bash"
    exit 1
fi

if ! command -v docker compose &> /dev/null; then
    echo "错误: 请先安装 Docker Compose"
    exit 1
fi

# 收集配置
echo "请输入以下配置信息："
echo ""

read -p "数据库密码: " DB_PASSWORD
read -p "JWT Secret (随机字符串): " JWT_SECRET
read -p "加密密钥 (必须是32字节): " ENCRYPTION_KEY
read -p "部署模式 (cloud=官方云端 / local=本地部署，默认 local): " DEPLOY_MODE

DEPLOY_MODE=${DEPLOY_MODE:-local}
EMBEDDING_DIM=${EMBEDDING_DIM:-1536}

# 导出环境变量
export DB_PASSWORD JWT_SECRET ENCRYPTION_KEY DEPLOY_MODE EMBEDDING_DIM

# 创建 .env 文件
cat > .env << EOF
DATABASE_URL=postgres://zhizhou:${DB_PASSWORD}@localhost:5432/zhizhou?sslmode=disable
JWT_SECRET=${JWT_SECRET}
ENCRYPTION_KEY=${ENCRYPTION_KEY}
DEPLOY_MODE=${DEPLOY_MODE}
EMBEDDING_DIM=${EMBEDDING_DIM}
SERVER_PORT=8080
EOF

echo ""
echo "配置已保存到 .env 文件"
echo ""

# 启动服务
echo "正在启动服务..."
docker compose -f deploy/docker-compose.yml up -d

echo ""
echo "========================================="
echo "  部署完成！"
echo "========================================="
echo ""
echo "  访问地址: http://localhost:3000"
echo "  API 地址: http://localhost:8080"
echo ""
echo "  常用命令:"
echo "    docker compose -f deploy/docker-compose.yml logs -f  查看日志"
echo "    docker compose -f deploy/docker-compose.yml down     停止服务"
echo "    docker compose -f deploy/docker-compose.yml restart  重启服务"
echo ""