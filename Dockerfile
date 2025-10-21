# 使用 Ubuntu 24.04 作为基础镜像
FROM ubuntu:24.04

# 设置环境变量
ENV DEBIAN_FRONTEND=noninteractive
ENV GO_VERSION=1.22.0
ENV SPEECHSDK_ROOT=/opt/speechsdk

# 安装系统依赖
RUN apt-get update && apt-get install -y \
    build-essential \
    ca-certificates \
    libasound2-dev \
    libssl-dev \
    wget \
    curl \
    git \
    pkg-config \
    && rm -rf /var/lib/apt/lists/*

# 安装 Go
RUN wget https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz \
    && rm go${GO_VERSION}.linux-amd64.tar.gz

# 设置 Go 环境变量
ENV PATH=/usr/local/go/bin:$PATH
ENV GOPATH=/go
ENV GOROOT=/usr/local/go
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
ENV CGO_ENABLED=1

# 创建 Speech SDK 目录
RUN mkdir -p $SPEECHSDK_ROOT

# 下载并安装 Azure Speech SDK
RUN wget -O SpeechSDK-Linux.tar.gz https://aka.ms/csspeech/linuxbinary \
    && tar --strip 1 -xzf SpeechSDK-Linux.tar.gz -C $SPEECHSDK_ROOT \
    && rm SpeechSDK-Linux.tar.gz

# 设置 Speech SDK 环境变量
ENV CGO_CFLAGS="-I$SPEECHSDK_ROOT/include/c_api"
ENV CGO_LDFLAGS="-L$SPEECHSDK_ROOT/lib/x64 -lMicrosoft.CognitiveServices.Speech.core"
ENV LD_LIBRARY_PATH="$SPEECHSDK_ROOT/lib/x64:$LD_LIBRARY_PATH"

# 复制库文件到系统路径（可选，便于链接）
RUN cp $SPEECHSDK_ROOT/lib/x64/libMicrosoft.CognitiveServices.Speech.core.so /usr/local/lib/ \
    && ldconfig

# 配置 Git (Go 模块需要)
RUN git config --global user.email "docker@example.com" \
    && git config --global user.name "Docker Build"

# 设置工作目录
WORKDIR /workspace

# 复制项目文件
COPY . .

# 初始化 Go 模块并添加依赖
RUN go mod init azure-speech-service || true \
    && go get github.com/Microsoft/cognitive-services-speech-sdk-go@latest \
    && go mod tidy

# 构建应用
RUN mkdir -p bin \
    && go build -o bin/text-to-speech text_to_speech.go 

# 暴露端口
EXPOSE 8080

# 启动命令
CMD ["bash"]