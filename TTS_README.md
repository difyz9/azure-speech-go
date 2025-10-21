# Azure 语音服务 - 文本转语音示例

## 快速开始

### 1. 设置环境变量

```bash
export SPEECH_KEY='your-azure-speech-key'
export SPEECH_REGION='your-azure-region'  # 例如: eastasia
```

### 2. 构建并运行 Docker 容器

```bash
# 构建 Docker 镜像
docker-compose build

# 运行文本转语音程序
./run_tts.sh
```

或者手动运行：

```bash
# 创建输出目录
mkdir -p output

# 运行容器
docker-compose run --rm azure-speech-go /workspace/bin/text-to-speech
```

### 3. 查看生成的音频文件

音频文件会保存在 `./output/` 目录下，格式为 `speech_N_timestamp.wav`

```bash
ls -lh output/
```

## 测试文本

程序会将以下文本转换为语音：

1. 你好，这是 Azure 语音服务的测试。
2. 今天天气很好，适合出去散步。
3. 人工智能正在改变我们的生活方式。
4. 学习新技能需要耐心和持续的努力。
5. 科技让世界变得更加美好和便捷。

## 自定义文本

如果需要自定义文本，请修改 `text_to_speech.go` 文件中的 `texts` 数组。

## 音频格式

- 格式：WAV
- 采样率：16kHz
- 比特率：32kbps
- 声道：单声道
- 语音：zh-CN-XiaoxiaoNeural (中文女声)

## 故障排除

1. 如果遇到权限问题，确保 output 目录有写入权限
2. 如果提示缺少环境变量，请重新设置 SPEECH_KEY 和 SPEECH_REGION
3. 确保 Azure 语音服务密钥有效且配额充足
