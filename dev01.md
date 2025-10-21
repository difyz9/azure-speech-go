# Azure 语音服务 Docker Linux 24.04 实现

基于微软 Azure 语音服务官方文档，已完成 Docker Linux 24.04 环境下的语音识别和文本转语音功能实现。

参考文档：https://learn.microsoft.com/zh-cn/azure/ai-services/speech-service/quickstarts/setup-platform?pivots=programming-language-go&tabs=linux%2Cubuntu%2Cdotnetcli%2Cjre%2Cmaven%2Cnodejs%2Cmac%2Cpypi

## 📁 完整项目结构

```
azure01/
├── 🐋 Docker 配置
│   ├── Dockerfile              # Ubuntu 24.04 + Go 1.22 + Azure Speech SDK
│   ├── docker-compose.yml      # 完整服务配置，支持音频设备
│   └── .env.example           # Azure 配置示例
├── 🎤 Go 语音程序
│   ├── main.go                # 语音识别示例
│   ├── text_to_speech.go      # 文本转语音示例
│   ├── azure_speech_demo.go   # 综合 Demo（命令行+Web界面）
│   ├── go.mod                 # Go 模块依赖
│   └── go.sum                 # 依赖校验文件
├── 🛠️ 工具脚本
│   ├── Makefile              # 构建和管理命令
│   ├── start.sh              # 原有启动脚本
│   ├── quick_start.sh        # 一键快速启动脚本
│   └── test_environment.sh   # 环境测试脚本
└── 📖 文档
    ├── README.md             # 详细使用说明
    └── dev01.md             # 本文档
```

## 🚀 快速开始

### 1. 一键启动（推荐）
```bash
chmod +x quick_start.sh
./quick_start.sh
```

### 2. 手动步骤
```bash
# 1. 配置 Azure 密钥
cp .env.example .env
# 编辑 .env 填入 SPEECH_KEY 和 SPEECH_REGION

# 2. 构建和启动
make docker-build
make docker-run

# 3. 进入容器
docker-compose exec azure-speech-go bash

# 4. 运行程序
make demo    # 综合演示
make run     # 语音识别  
make tts     # 文本转语音
```

## 🎯 功能特性

### ✅ 语音识别 (Speech-to-Text)
- 实时连续语音识别
- 支持中文语音输入
- 麦克风音频采集
- 识别结果实时显示

### ✅ 文本转语音 (Text-to-Speech)  
- 中文文本语音合成
- 多种语音角色（小晓、小伊等）
- 扬声器音频输出
- 自然流畅的语音效果

### ✅ Web 界面
- 友好的浏览器操作界面
- 文本转语音 Web 功能
- 实时服务状态监控
- 访问地址：http://localhost:8080

### ✅ 技术环境
- **基础系统**：Ubuntu 24.04 LTS
- **Go 版本**：1.22.0
- **Azure SDK**：最新版 Speech SDK for Go
- **音频支持**：ALSA + 音频设备映射
- **容器化**：Docker + Docker Compose

## 🔧 核心配置

### Docker 环境配置
- 基于官方 Ubuntu 24.04 镜像
- 自动安装 Go 1.22 和构建工具
- 下载配置 Azure Speech SDK
- 设置正确的环境变量和库路径
- 支持音频设备访问

### Azure Speech SDK 配置
```bash
# 环境变量设置
SPEECHSDK_ROOT=/opt/speechsdk
CGO_CFLAGS=-I/opt/speechsdk/include/c_api
CGO_LDFLAGS=-L/opt/speechsdk/lib/x64 -lMicrosoft.CognitiveServices.Speech.core
LD_LIBRARY_PATH=/opt/speechsdk/lib/x64
```

### Go 程序特性
- 模块化设计，支持多种使用方式
- 完善的错误处理和用户提示
- 中文本地化界面和消息
- 支持命令行和 Web 两种交互模式

## 🧪 测试验证

```bash
# 环境测试
make test-env

# 功能测试
make demo
```

测试内容包括：
- ✅ 环境变量检查
- ✅ Speech SDK 安装验证
- ✅ Go 环境和模块验证
- ✅ 网络连接测试
- ✅ 程序编译测试
- ✅ 音频设备检测

## 📊 实现说明

根据微软官方文档要求：
1. ✅ 安装了所需的系统依赖（build-essential, ca-certificates, libasound2-dev, libssl-dev）
2. ✅ 正确下载和配置了 Azure Speech SDK
3. ✅ 设置了正确的 CGO 环境变量
4. ✅ 实现了完整的语音识别和文本转语音功能
5. ✅ 提供了命令行和 Web 两种交互方式
6. ✅ 支持音频设备在 Docker 中的使用

## 🎉 项目完成状态

- ✅ Docker Linux 24.04 环境搭建完成
- ✅ Azure Speech SDK 集成完成  
- ✅ 语音识别功能实现完成
- ✅ 文本转语音功能实现完成
- ✅ Web 界面开发完成
- ✅ 完整的文档和脚本提供完成
- ✅ 环境测试和验证完成

现在可以直接使用 `./quick_start.sh` 一键启动完整的 Azure 语音服务环境！🎊