#!/bin/bash

# 检查 .env 文件是否存在
if [ ! -f .env ]; then
    echo "创建 .env 文件..."
    cp .env.example .env
    echo "请编辑 .env 文件，填入您的 Azure Speech Service 密钥和区域"
    exit 1
fi

# 启动容器
echo "启动 Azure Speech SDK Go 开发环境..."
docker-compose up -d

# 等待容器启动
echo "等待容器初始化..."
sleep 5

# 显示容器状态
docker-compose ps

echo ""
echo "环境已启动！使用以下命令进入容器："
echo "docker-compose exec azure-speech-go bash"