.PHONY: help build run clean test demo docker-build docker-run docker-stop test-env

help:
	@echo "可用命令:"
	@echo "  make build        - 构建 Go 应用"
	@echo "  make run          - 运行语音识别示例"
	@echo "  make tts          - 运行文本转语音示例"
	@echo "  make demo         - 运行综合 Demo 程序"
	@echo "  make docker-build - 构建 Docker 镜像"
	@echo "  make docker-run   - 运行 Docker 容器"
	@echo "  make docker-stop  - 停止 Docker 容器"
	@echo "  make test         - 运行测试"
	@echo "  make test-env     - 测试环境配置"
	@echo "  make clean        - 清理构建文件"

build:
	go mod tidy
	go build -o bin/speech-recognition main.go
	go build -o bin/text-to-speech text_to_speech.go
	go build -o bin/azure-speech-demo azure_speech_demo.go

run:
	go run main.go

tts:
	go run text_to_speech.go

demo:
	go run azure_speech_demo.go

docker-build:
	docker-compose build

docker-run:
	docker-compose up -d
	@echo "Docker 容器已启动，使用以下命令进入容器："
	@echo "docker-compose exec azure-speech-go bash"

docker-stop:
	docker-compose down

test:
	go test -v ./...

test-env:
	chmod +x test_environment.sh
	./test_environment.sh

clean:
	rm -rf bin/
	go clean