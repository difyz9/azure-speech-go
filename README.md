# Azure 语音服务 Go SDK - Docker Linux 24.04 环境

🎤 使用 Docker 在 Ubuntu 24.04 环境中运行 Azure 语音识别和文本转语音服务的完整解决方案。

## 📋 前置要求

- Docker 和 Docker Compose
- Azure 语音服务订阅密钥
- 支持音频输入/输出的设备

## 🚀 快速开始

### 1. 获取 Azure 语音服务密钥

1. 访问 [Azure Portal](https://portal.azure.com)
2. 创建"语音服务"资源
3. 在"密钥和终结点"页面获取：
   - 密钥 (Key)
   - 区域 (Region)

### 2. 配置环境变量

复制环境变量示例文件：

```bash
cp .env.example .env
```

编辑 `.env` 文件，填入您的 Azure 配置：

```bash
SPEECH_KEY=您的语音服务密钥
SPEECH_REGION=您的服务区域
```

常用区域示例：
- `eastus` - 美国东部
- `westus2` - 美国西部 2  
- `eastasia` - 东亚
- `southeastasia` - 东南亚
- `chinaeast2` - 中国东部 2（中国区域）

### 3. 构建和启动 Docker 环境

```bash
# 构建 Docker 镜像（基于 Ubuntu 24.04）
make docker-build

# 启动容器服务
```

### 4. 进入容器并运行程序

```bash
# 进入运行中的容器
docker-compose exec azure-speech-go bash

# 在容器中运行不同的示例程序
```

## 🎯 功能演示

### 方式 1：命令行模式

```bash
# 进入容器后运行综合 Demo
make demo

# 或者单独运行各个功能
make run     # 语音识别
make tts     # 文本转语音
```

### 方式 2：Web 界面模式

```bash
# 在容器中启动 Web 服务
make demo
# 然后选择 'w' 启动 Web 服务器

# 在浏览器中访问：http://localhost:8080
```

## 📁 项目结构

```
.
├── Dockerfile              # Ubuntu 24.04 基础镜像配置
├── docker-compose.yml      # Docker Compose 服务配置
├── go.mod                  # Go 模块定义
├── go.sum                  # Go 模块依赖校验
├── .env.example           # 环境变量配置示例
├── .env                   # 环境变量（需自行创建）
├── start.sh              # 启动脚本
├── main.go               # 语音识别示例
├── text_to_speech.go     # 文本转语音示例
├── azure_speech_demo.go  # 综合 Demo 程序
├── Makefile             # 构建和管理脚本
└── README.md            # 项目说明文档
```

## 🎤 功能说明

### 1. 语音识别 (Speech-to-Text)
- 实时语音识别
- 支持中文语音输入
- 连续识别模式
- 麦克风音频采集

### 2. 文本转语音 (Text-to-Speech)
- 中文文本转语音
- 多种语音角色支持
- 自然语音合成
- 扬声器音频输出

### 3. Web 界面
- 友好的 Web 操作界面
- 实时文本转语音功能
- 服务状态监控

## 🔧 技术特性

- **基础环境**：Ubuntu 24.04 LTS
- **Go 版本**：1.22.0
- **Azure Speech SDK**：最新版本
- **音频支持**：ALSA 音频系统
- **SSL 支持**：OpenSSL 3.x
- **容器化**：Docker + Docker Compose

## ⚠️ 常见问题

### 1. 音频设备权限

如需在容器中使用音频设备，请在 `docker-compose.yml` 中添加：

```yaml
devices:
  - /dev/snd:/dev/snd
volumes:
  - /dev/shm:/dev/shm
```

### 2. 中国区域配置

如使用 Azure 中国区域，请确保：
- 使用正确的区域代码（如：`chinaeast2`）
- 访问中国区域的 Azure Portal

### 3. 网络连接问题

如遇到网络问题，可以尝试：
```bash
# 检查网络连接
docker-compose exec azure-speech-go ping cognitive.azure.com

# 检查 DNS 解析
docker-compose exec azure-speech-go nslookup cognitive.azure.com
```

## 🛠️ 高级配置

### 自定义语音设置

在代码中可以修改以下语音参数：

```go
// 语音识别语言
config.SetSpeechRecognitionLanguage("zh-CN")

// 文本转语音语音角色
config.SetSpeechSynthesisVoiceName("zh-CN-XiaoxiaoNeural")

// 其他可用的中文语音角色：
// zh-CN-XiaoyiNeural (女声)
// zh-CN-YunjianNeural (男声)  
// zh-CN-YunxiNeural (男声)
```

### 环境清理

```bash
# 停止所有服务
make docker-stop

# 清理构建文件
make clean

# 完全清理（包括 Docker 镜像）
docker-compose down --rmi all --volumes --remove-orphans
```

## 📚 参考资料

- [Azure 语音服务官方文档](https://learn.microsoft.com/zh-cn/azure/ai-services/speech-service/)
- [Go 语言 Speech SDK](https://github.com/Microsoft/cognitive-services-speech-sdk-go)
- [Azure 语音服务定价](https://azure.microsoft.com/pricing/details/cognitive-services/speech-services/)
- [Ubuntu 24.04 LTS 发布说明](https://releases.ubuntu.com/24.04/)

## 📝 许可证

本项目遵循 MIT 许可证。Azure Speech SDK 遵循微软的许可条款。