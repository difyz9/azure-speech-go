# Azure Speech & Translation Toolkit

基于 Azure 认知服务的语音工具包，包含文本转语音 (TTS) API 和文本翻译功能，支持多种语言和音频格式。

## 🚀 快速开始

### 环境要求

- Docker & Docker Compose
- Azure 认知服务订阅 (Speech Service)
- Azure 翻译服务订阅 (Translator Service) - 用于文本翻译功能

### 1. 克隆项目

```bash
git clone https://github.com/difyz9/azure-speech-go
cd azure-speech-go
```

### 2. 设置环境变量

创建 `.env` 文件或设置环境变量：

```bash
export SPEECH_KEY="your-azure-speech-key"
export SPEECH_REGION="your-azure-region"  # 例如: eastus, westus2
```

**注意**: `SPEECH_KEY` 和 `SPEECH_REGION` 同时用于语音服务和翻译服务。

### 3. 构建基础镜像

首先构建包含 Azure Speech SDK 的基础镜像：

```bash
sudo docker build -f Dockerfile.base -t azure-speech-sdk:latest .
```

### 4. 启动服务

```bash
# 构建并启动服务
docker-compose up -d azure-tts-api

# 或使用重建脚本（推荐）
./rebuild.sh
```

### 5. 验证服务

```bash
# 运行测试脚本
./test-api.sh

# 或手动检查健康状态
curl http://localhost:8080/api/health
```

## 📚 功能特性

### 🎤 文本转语音 (TTS)
- 支持单句和批量文本转语音
- 多种语言和语音选择
- WAV/MP3 音频格式输出
- 音频文件管理和下载

### 🔄 文本翻译
- 基于 Azure Translator API
- 支持多语言互译
- 简单易用的翻译示例

## 📚 API 文档

### 基础信息

- **服务地址**: `http://localhost:8080`
- **API 基础路径**: `/api`
- **文档**: `http://localhost:8080/`

### 主要接口

#### 1. 健康检查

```http
GET /api/health
```

#### 2. 单句文本转语音

```http
POST /api/tts
Content-Type: application/json

{
  "text": "你好，这是测试文本",
  "language": "zh-CN",
  "voice": "zh-CN-XiaoxiaoNeural",
  "format": "wav"
}
```

**响应**:
```json
{
  "success": true,
  "message": "语音合成成功",
  "filename": "tts_20241022_150405.wav",
  "duration": "2.5s"
}
```

#### 3. 批量文本转语音

```http
POST /api/batch-tts
Content-Type: application/json

{
  "texts": ["第一句话", "第二句话", "第三句话"],
  "language": "zh-CN",
  "voice": "zh-CN-XiaoxiaoNeural",
  "format": "wav"
}
```

#### 4. 获取文件列表

```http
GET /api/files?limit=10
```

#### 5. 下载音频文件

```http
GET /api/download/{filename}
```

## 🎯 支持的语言和语音

### 中文
- `zh-CN-XiaoxiaoNeural` (女声)
- `zh-CN-YunxiNeural` (男声)
- `zh-CN-YunyangNeural` (男声)

### 英文
- `en-US-JennyNeural` (女声)
- `en-US-GuyNeural` (男声)
- `en-US-AriaNeural` (女声)

### 音频格式
- `wav` - WAV 格式 (默认)
- `mp3` - MP3 格式

## 🛠️ 开发工具

### 重建服务

```bash
./rebuild.sh
```

### 测试 API

```bash
./test-api.sh
```

### 文本翻译示例

项目包含独立的翻译功能示例 (`api/translate.go`)，可以直接运行：

```bash
# 进入 api 目录
cd api

# 设置环境变量
export SPEECH_KEY="your-azure-speech-key"
export SPEECH_REGION="your-azure-region"

# 运行翻译示例
go run translate.go
```

**翻译示例说明**:
- 默认从英文翻译到中文 (en → zh)
- 示例文本: "I would really like to drive your car around the block a few times."
- 返回 JSON 格式的翻译结果

### 查看日志

```bash
docker-compose logs -f azure-tts-api
```

### 进入容器

```bash
docker-compose exec azure-tts-api /bin/bash
```

## 📁 项目结构

```
azure-speech-go-toolkit/
├── api/                    # Go 源码目录
│   ├── main.go            # TTS API 主程序
│   ├── translate.go       # 文本翻译示例
│   ├── go.mod             # Go 依赖管理
│   └── go.sum
├── output/                 # 音频文件输出目录
├── Dockerfile.base         # 基础镜像 Dockerfile
├── Dockerfile.api          # API 应用 Dockerfile
├── docker-compose.yml      # Docker Compose 配置
├── rebuild.sh             # 重建脚本
├── test-api.sh            # API 测试脚本
└── README.md              # 项目文档
```

## 🔧 配置说明

### 环境变量

| 变量名 | 说明 | 默认值 | 必需 |
|--------|------|--------|------|
| `SPEECH_KEY` | Azure 语音/翻译服务密钥 | - | ✅ |
| `SPEECH_REGION` | Azure 服务区域 | - | ✅ |
| `OUTPUT_DIR` | 音频文件输出目录 | `/app/output` | ❌ |
| `PORT` | API 服务端口 | `8080` | ❌ |
| `GIN_MODE` | Gin 框架模式 | `release` | ❌ |

### Docker 配置

- **基础镜像**: `ubuntu:24.04`
- **Go 版本**: `1.22.0`
- **Azure Speech SDK**: 最新版本
- **端口映射**: `8080:8080`
- **Volume 挂载**: `./output:/app/output`

## 🚨 故障排除

### 1. 权限问题

如果遇到输出目录权限问题：

```bash
# 设置输出目录权限
mkdir -p ./output
chmod 755 ./output
```

### 2. 基础镜像不存在

```bash
# 重新构建基础镜像
sudo docker build -f Dockerfile.base -t azure-speech-sdk:latest .
```

### 3. Azure 服务配置错误

检查环境变量：
```bash
echo $SPEECH_KEY
echo $SPEECH_REGION
```

### 4. 服务启动失败

查看详细日志：
```bash
docker-compose logs azure-tts-api
```

## 📝 使用示例

### cURL 示例

```bash
# 健康检查
curl http://localhost:8080/api/health

# 生成中文语音
curl -X POST http://localhost:8080/api/tts \
  -H "Content-Type: application/json" \
  -d '{
    "text": "你好，欢迎使用Azure语音服务",
    "language": "zh-CN",
    "voice": "zh-CN-XiaoxiaoNeural",
    "format": "wav"
  }'

# 生成英文语音
curl -X POST http://localhost:8080/api/tts \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Hello, welcome to Azure Speech Service",
    "language": "en-US",
    "voice": "en-US-JennyNeural",
    "format": "mp3"
  }'

# 获取文件列表
curl http://localhost:8080/api/files

# 下载文件
curl -O http://localhost:8080/api/download/tts_20241022_150405.wav
```

### JavaScript 示例

```javascript
// 文本转语音
const response = await fetch('http://localhost:8080/api/tts', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    text: '这是一个测试文本',
    language: 'zh-CN',
    voice: 'zh-CN-XiaoxiaoNeural',
    format: 'wav'
  })
});

const result = await response.json();
console.log(result);

// 下载生成的文件
if (result.success) {
  const audioUrl = `http://localhost:8080/api/download/${result.filename}`;
  window.open(audioUrl);
}
```

### Go 翻译示例

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "os"
    "net/http"
    "net/url"
)

func translateText() {
    key := os.Getenv("SPEECH_KEY")
    endpoint := "https://api.cognitive.microsofttranslator.com/"
    uri := endpoint + "/translate?api-version=3.0"
    location := os.Getenv("SPEECH_REGION")

    // 构建请求URL
    u, _ := url.Parse(uri)
    q := u.Query()
    q.Add("from", "en")  // 源语言
    q.Add("to", "zh")    // 目标语言
    u.RawQuery = q.Encode()

    // 请求体
    body := []struct {
        Text string
    }{
        {Text: "Hello, how are you?"},
    }
    b, _ := json.Marshal(body)

    // 创建HTTP请求
    req, _ := http.NewRequest("POST", u.String(), bytes.NewBuffer(b))
    req.Header.Add("Ocp-Apim-Subscription-Key", key)
    req.Header.Add("Ocp-Apim-Subscription-Region", location)
    req.Header.Add("Content-Type", "application/json")

    // 发送请求
    res, _ := http.DefaultClient.Do(req)
    
    // 解析响应
    var result interface{}
    json.NewDecoder(res.Body).Decode(&result)
    prettyJSON, _ := json.MarshalIndent(result, "", "  ")
    fmt.Printf("%s\n", prettyJSON)
}
```

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License

## 🔗 相关链接

- [Azure 认知服务文档](https://docs.microsoft.com/azure/cognitive-services/speech-service/)
- [Azure Speech SDK for Go](https://github.com/Microsoft/cognitive-services-speech-sdk-go)
- [Azure Translator 文档](https://docs.microsoft.com/azure/cognitive-services/translator/)
- [Azure Translator API 参考](https://docs.microsoft.com/azure/cognitive-services/translator/reference/v3-0-translate)
- [Docker 官方文档](https://docs.docker.com/)
- [Go 官方文档](https://golang.org/doc/)