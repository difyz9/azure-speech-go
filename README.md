# Azure 语音服务 Go 开发环境

🎤 基于 Docker 的 Azure 认知服务语音 SDK Go 语言开发环境，支持语音识别、文本转语音等功能。

## 📋 项目概述

这是一个完整的 Azure 语音服务 Go 语言开发环境，提供了：

- 🗣️ **语音识别** (Speech-to-Text)
- 🔊 **文本转语音** (Text-to-Speech) 
- 🌐 **Web 界面演示**
- 📦 **批量语音合成**
- 🐳 **Docker 容器化部署**

## 🚀 快速开始

### 1. 环境准备

确保您已安装：
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

### 2. Azure 配置

1. 在 [Azure Portal](https://portal.azure.com/) 创建语音服务资源
2. 获取 API 密钥和区域信息
3. 创建 `.env` 文件：

```bash
# 复制环境变量模板
cp .env.example .env
```

编辑 `.env` 文件，填入您的 Azure 配置：

```env
# Azure 语音服务配置
SPEECH_KEY=your_speech_service_key_here
SPEECH_REGION=your_region_here
```

### 3. 启动容器

```bash
# 构建并启动容器
docker-compose up -d

# 进入容器
docker-compose exec azure-speech-go bash
```

## 🎯 功能演示

### 批量文本转语音 (`main.go`)

在容器中运行批量语音合成：

```bash
# 运行批量语音合成
go run main.go
```

功能特点：
- 支持中文语音合成
- 使用 `zh-CN-XiaoxiaoNeural` 声音
- 自动生成带时间戳的 WAV 文件
- 输出文件保存到 `/workspace/output` 目录

### 交互式语音服务 (`azure_speech_demo.go`)

运行完整的语音服务演示：

```bash
# 启动交互式演示
go run azure_speech_demo.go
```

支持的功能：
- `r` / `recognize`: 语音识别（需要麦克风）
- `t` / `tts`: 文本转语音
- `w` / `web`: 启动 Web 服务器 (http://localhost:8080)
- `q` / `quit`: 退出程序

### Web 界面访问

启动 Web 服务后，访问：
- **本地访问**: http://localhost:8080
- **容器内访问**: http://container_ip:8080

## 📁 项目结构

```
├── main.go                 # 批量语音合成主程序
├── azure_speech_demo.go    # 交互式语音服务演示
├── go.mod                  # Go 模块依赖
├── Dockerfile              # Docker 镜像构建文件
├── docker-compose.yml      # Docker Compose 配置
├── .env.example           # 环境变量模板
├── output/                # 生成的音频文件目录
└── README.md              # 项目文档
```

## 🔧 技术栈

- **语言**: Go 1.22
- **SDK**: Microsoft Cognitive Services Speech SDK
- **容器**: Docker + Ubuntu 24.04
- **音频格式**: WAV, MP3
- **语音模型**: zh-CN-XiaoxiaoNeural

## 🛠️ 开发指南

### 本地开发

如果您想在本地开发而不使用 Docker：

1. 安装 Go 1.22+
2. 下载 [Azure Speech SDK for Linux](https://aka.ms/csspeech/linuxbinary)
3. 设置环境变量：

```bash
export SPEECHSDK_ROOT=/path/to/speechsdk
export CGO_CFLAGS="-I$SPEECHSDK_ROOT/include/c_api"
export CGO_LDFLAGS="-L$SPEECHSDK_ROOT/lib/x64 -lMicrosoft.CognitiveServices.Speech.core"
export LD_LIBRARY_PATH="$SPEECHSDK_ROOT/lib/x64:$LD_LIBRARY_PATH"
```

### 添加新功能

1. 修改 `azure_speech_demo.go` 添加新的交互选项
2. 在 `main.go` 中扩展批量处理逻辑
3. 更新 Docker 配置以支持新的依赖

## 🎵 音频配置

### 支持的音频格式
- **输入**: 16kHz, 16-bit, 单声道 WAV
- **输出**: WAV, MP3 (可配置)

### 语音选项
- **语言**: 中文 (zh-CN)
- **声音**: XiaoxiaoNeural (可更换为其他中文声音)
- **质量**: 16kHz, 32kbps (MP3) / 16kHz, 16-bit (WAV)

## 🚨 故障排除

### 常见问题

**1. 环境变量未设置**
```
❌ 请设置 SPEECH_KEY 和 SPEECH_REGION 环境变量
```
解决方案: 确保 `.env` 文件配置正确

**2. 音频设备权限问题**
```
❌ 创建音频配置失败
```
解决方案: 检查 Docker 音频设备映射配置

**3. 网络连接问题**
```
❌ 语音合成失败: network error
```
解决方案: 检查网络连接和 Azure 服务状态

### 调试命令

```bash
# 查看容器日志
docker-compose logs azure-speech-go

# 检查 Go 模块
go mod verify

# 测试网络连接
curl -I https://cognitiveservices.azure.com/
```

## 📚 参考资料

- [Azure 语音服务文档](https://docs.microsoft.com/zh-cn/azure/cognitive-services/speech-service/)
- [Speech SDK for Go](https://github.com/Microsoft/cognitive-services-speech-sdk-go)
- [支持的语音列表](https://docs.microsoft.com/zh-cn/azure/cognitive-services/speech-service/language-support)

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 👥 联系方式

如果您有任何问题或建议，请通过以下方式联系：

- 📧 邮件: your-email@example.com
- 🐛 问题反馈: [GitHub Issues](https://github.com/your-username/docker-azure-golang-env/issues)

---

⭐ 如果这个项目对您有帮助，请给个 Star！