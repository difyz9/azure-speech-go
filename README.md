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
docker compose up -d

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


                                                    
```

ubuntu@VM-0-14-ubuntu:~/azure01$ sudo docker ps
CONTAINER ID   IMAGE                     COMMAND   CREATED          STATUS          PORTS                                         NAMES
feacc58f39c2   azure01-azure-speech-go   "bash"    11 seconds ago   Up 11 seconds   0.0.0.0:8080->8080/tcp, [::]:8080->8080/tcp   azure-speech-go-dev
ubuntu@VM-0-14-ubuntu:~/azure01$ sudo docker exec -it feacc58f39c2 /bin/bash
root@feacc58f39c2:/workspace# ls
Dockerfile  README.md      azure_speech_demo.go  docker-compose.yml  go.sum   output          run_tts.sh  test_environment.sh
Makefile    TTS_README.md  dev01.md              go.mod              main.go  quick_start.sh  start.sh    text_to_speech.go
root@feacc58f39c2:/workspace# ll
total 92
drwxr-xr-x 3 ubuntu 1001 4096 Oct 21 10:08 ./
drwxr-xr-x 1 root   root 4096 Oct 21 10:08 ../
-rw-r--r-- 1 ubuntu 1001  116 Oct 21 08:17 .env
-rw-r--r-- 1 ubuntu 1001  527 Oct 21 08:31 .env.example
-rw-r--r-- 1 ubuntu 1001 1986 Oct 21 10:05 Dockerfile
-rw-r--r-- 1 ubuntu 1001 1238 Oct 21 08:35 Makefile
-rw-r--r-- 1 ubuntu 1001 4678 Oct 21 08:33 README.md
-rw-r--r-- 1 ubuntu 1001 1438 Oct 21 10:05 TTS_README.md
-rw-r--r-- 1 ubuntu 1001 7035 Oct 21 08:32 azure_speech_demo.go
-rw-r--r-- 1 ubuntu 1001 4345 Oct 21 08:45 dev01.md
-rw-r--r-- 1 ubuntu 1001 1060 Oct 21 10:04 docker-compose.yml
-rw-r--r-- 1 ubuntu 1001  108 Oct 21 08:57 go.mod
-rw-r--r-- 1 ubuntu 1001  227 Oct 21 08:56 go.sum
-rw-r--r-- 1 ubuntu 1001 1667 Oct 21 08:14 main.go
drwxr-xr-x 2 root   root 4096 Oct 21 10:08 output/
-rw-r--r-- 1 ubuntu 1001 2367 Oct 21 08:35 quick_start.sh
-rwxr-xr-x 1 ubuntu 1001  705 Oct 21 10:05 run_tts.sh*
-rw-r--r-- 1 ubuntu 1001  530 Oct 21 08:14 start.sh
-rw-r--r-- 1 ubuntu 1001 2676 Oct 21 08:34 test_environment.sh
-rw-r--r-- 1 ubuntu 1001 3257 Oct 21 10:04 text_to_speech.go
root@feacc58f39c2:/workspace# cd output/  
root@feacc58f39c2:/workspace/output# ls
root@feacc58f39c2:/workspace/output# ll
total 8
drwxr-xr-x 2 root   root 4096 Oct 21 10:08 ./
drwxr-xr-x 3 ubuntu 1001 4096 Oct 21 10:08 ../
root@feacc58f39c2:/workspace/output# cd ..
root@feacc58f39c2:/workspace# go run text_to_speech.go 
开始批量语音合成...
===========================================
正在合成: 你好，这是 Azure 语音服务的测试。
✓ 音频已保存到: /workspace/output/speech_1_20251021_100859.wav
-------------------------------------------
正在合成: 今天天气很好，适合出去散步。
✓ 音频已保存到: /workspace/output/speech_2_20251021_100901.wav
-------------------------------------------
正在合成: 人工智能正在改变我们的生活方式。
✓ 音频已保存到: /workspace/output/speech_3_20251021_100903.wav
-------------------------------------------
正在合成: 学习新技能需要耐心和持续的努力。
✓ 音频已保存到: /workspace/output/speech_4_20251021_100904.wav
-------------------------------------------
正在合成: 科技让世界变得更加美好和便捷。
✓ 音频已保存到: /workspace/output/speech_5_20251021_100906.wav
-------------------------------------------
===========================================
合成完成！成功: 5/5
音频文件保存在: /workspace/output
root@feacc58f39c2:/workspace# 

```