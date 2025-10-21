#!/bin/bash

echo "🎤 Azure 语音服务 Docker 环境快速启动"
echo "======================================"

# 检查必要文件
if [ ! -f ".env" ]; then
    echo "📝 创建环境配置文件..."
    if [ -f ".env.example" ]; then
        cp .env.example .env
        echo "✅ 已创建 .env 文件"
        echo "⚠️  请编辑 .env 文件，填入您的 Azure Speech Service 配置："
        echo "   - SPEECH_KEY=您的语音服务密钥"
        echo "   - SPEECH_REGION=您的服务区域"
        echo ""
        echo "💡 获取密钥的步骤："
        echo "   1. 访问 https://portal.azure.com"
        echo "   2. 创建'语音服务'资源"
        echo "   3. 在'密钥和终结点'页面获取配置"
        echo ""
        read -p "是否现在打开 .env 文件进行编辑？(y/n): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            ${EDITOR:-nano} .env
        fi
    else
        echo "❌ .env.example 文件不存在"
        exit 1
    fi
fi

# 检查 Docker
if ! command -v docker &> /dev/null; then
    echo "❌ Docker 未安装，请先安装 Docker"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose 未安装，请先安装 Docker Compose"
    exit 1
fi

echo "🏗️  构建 Docker 镜像..."
if docker-compose build; then
    echo "✅ Docker 镜像构建成功"
else
    echo "❌ Docker 镜像构建失败"
    exit 1
fi

echo "🚀 启动 Docker 容器..."
if docker-compose up -d; then
    echo "✅ Docker 容器启动成功"
else
    echo "❌ Docker 容器启动失败"
    exit 1
fi

echo "⏳ 等待容器初始化..."
sleep 5

echo "🧪 运行环境测试..."
docker-compose exec azure-speech-go bash -c "cd /workspace && make test-env"

echo ""
echo "🎉 环境准备完成！"
echo ""
echo "📖 使用说明："
echo "1. 进入容器：docker-compose exec azure-speech-go bash"
echo "2. 运行 Demo： make demo"
echo "3. 语音识别：  make run"
echo "4. 文本转语音：make tts"
echo "5. Web 界面：  make demo (选择 'w'，然后访问 http://localhost:8080)"
echo ""
echo "🛠️  管理命令："
echo "- 停止容器：   docker-compose down"
echo "- 查看日志：   docker-compose logs -f"
echo "- 重新构建：   docker-compose build --no-cache"
echo ""
echo "现在可以使用 Azure 语音服务了！ 🎊"