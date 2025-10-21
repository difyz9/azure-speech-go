#!/bin/bash

echo "🧪 Azure 语音服务环境测试脚本"
echo "================================="

# 检查环境变量
echo "📋 检查环境变量..."
if [ -z "$SPEECH_KEY" ]; then
    echo "❌ SPEECH_KEY 环境变量未设置"
    exit 1
else
    echo "✅ SPEECH_KEY 已设置"
fi

if [ -z "$SPEECH_REGION" ]; then
    echo "❌ SPEECH_REGION 环境变量未设置"
    exit 1
else
    echo "✅ SPEECH_REGION 已设置为: $SPEECH_REGION"
fi

# 检查 Speech SDK 安装
echo ""
echo "📦 检查 Speech SDK 安装..."
if [ -d "$SPEECHSDK_ROOT" ]; then
    echo "✅ Speech SDK 目录存在: $SPEECHSDK_ROOT"
    
    if [ -f "$SPEECHSDK_ROOT/lib/x64/libMicrosoft.CognitiveServices.Speech.core.so" ]; then
        echo "✅ Speech SDK 核心库存在"
    else
        echo "❌ Speech SDK 核心库不存在"
        exit 1
    fi
    
    if [ -d "$SPEECHSDK_ROOT/include/c_api" ]; then
        echo "✅ Speech SDK 头文件存在"
    else
        echo "❌ Speech SDK 头文件不存在"
        exit 1
    fi
else
    echo "❌ Speech SDK 目录不存在: $SPEECHSDK_ROOT"
    exit 1
fi

# 检查 Go 环境
echo ""
echo "🐹 检查 Go 环境..."
if command -v go &> /dev/null; then
    GO_VERSION=$(go version)
    echo "✅ Go 已安装: $GO_VERSION"
else
    echo "❌ Go 未安装"
    exit 1
fi

# 检查网络连接
echo ""
echo "🌐 检查网络连接..."
if curl -s --connect-timeout 5 https://cognitive.microsofttranslator.com > /dev/null; then
    echo "✅ 网络连接正常"
else
    echo "⚠️  网络连接可能有问题"
fi

# 检查 Go 模块
echo ""
echo "📚 检查 Go 模块..."
if [ -f "go.mod" ]; then
    echo "✅ go.mod 文件存在"
    if go mod verify; then
        echo "✅ Go 模块验证通过"
    else
        echo "❌ Go 模块验证失败"
        exit 1
    fi
else
    echo "❌ go.mod 文件不存在"
    exit 1
fi

# 尝试编译测试程序
echo ""
echo "🔨 编译测试..."
if go build -o test_build azure_speech_demo.go; then
    echo "✅ 程序编译成功"
    rm -f test_build
else
    echo "❌ 程序编译失败"
    exit 1
fi

# 检查音频设备（可选）
echo ""
echo "🎵 检查音频设备..."
if [ -d "/dev/snd" ]; then
    AUDIO_DEVICES=$(ls /dev/snd/ 2>/dev/null | wc -l)
    if [ $AUDIO_DEVICES -gt 0 ]; then
        echo "✅ 发现 $AUDIO_DEVICES 个音频设备"
    else
        echo "⚠️  未发现音频设备"
    fi
else
    echo "⚠️  音频设备目录不存在"
fi

echo ""
echo "🎉 环境测试完成！"
echo "💡 现在可以运行以下命令："
echo "   make demo  # 运行综合演示"
echo "   make run   # 语音识别"
echo "   make tts   # 文本转语音"