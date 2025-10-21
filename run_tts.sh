#!/bin/bash

# 运行文本到语音转换
echo "==================================="
echo "Azure 语音合成服务 - 文本转语音"
echo "==================================="

# 检查环境变量
if [ -z "$SPEECH_KEY" ] || [ -z "$SPEECH_REGION" ]; then
    echo "错误: 请先设置环境变量 SPEECH_KEY 和 SPEECH_REGION"
    echo "例如:"
    echo "  export SPEECH_KEY='your-key'"
    echo "  export SPEECH_REGION='your-region'"
    exit 1
fi

# 创建输出目录
mkdir -p output

# 使用 Docker Compose 运行
echo "启动 Docker 容器..."
docker-compose run --rm azure-speech-go /workspace/bin/text-to-speech

echo ""
echo "完成! 音频文件已保存到 ./output 目录"
ls -lh output/
