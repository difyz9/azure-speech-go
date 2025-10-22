#!/bin/bash

# Azure TTS API 测试脚本

# 配置 API 基础 URL（可通过环境变量覆盖）
API_BASE="${API_BASE:-http://43.160.253.168:8080/api}"

echo "🧪 Azure TTS API 测试脚本"
echo "========================="
echo "测试目标: $API_BASE"

# 检查依赖
check_dependencies() {
    local missing_deps=()
    
    if ! command -v curl &> /dev/null; then
        missing_deps+=("curl")
    fi
    
    if ! command -v jq &> /dev/null; then
        echo -e "${YELLOW}⚠️  建议安装 jq 来美化 JSON 输出${NC}"
    fi
    
    if [ ${#missing_deps[@]} -ne 0 ]; then
        echo -e "${RED}❌ 缺少依赖: ${missing_deps[*]}${NC}"
        echo "请安装缺少的依赖后重试"
        exit 1
    fi
}

# 等待服务启动
wait_for_service() {
    local max_attempts=30
    local attempt=1
    
    echo -e "\n${YELLOW}等待 API 服务启动...${NC}"
    
    while [ $attempt -le $max_attempts ]; do
        if curl -s --connect-timeout 2 "$API_BASE/health" > /dev/null 2>&1; then
            echo -e "${GREEN}✅ 服务已就绪${NC}"
            return 0
        fi
        
        echo -n "."
        sleep 2
        attempt=$((attempt + 1))
    done
    
    echo -e "\n${RED}❌ 服务启动超时${NC}"
    return 1
}

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 检查依赖
check_dependencies

# 测试函数
test_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4
    
    echo -e "\n${BLUE}📝 测试: $description${NC}"
    echo "请求: $method $API_BASE$endpoint"
    
    if [[ -n "$data" ]]; then
        echo "数据: $data" | jq . 2>/dev/null || echo "数据: $data"
    fi
    
    local start_time=$(date +%s.%3N)
    
    if [[ $method == "GET" ]]; then
        response=$(curl -s -w "\n%{http_code}" --connect-timeout 10 --max-time 30 "$API_BASE$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" --connect-timeout 10 --max-time 30 \
            -X "$method" -H "Content-Type: application/json" -d "$data" "$API_BASE$endpoint")
    fi
    
    local end_time=$(date +%s.%3N)
    local duration=$(echo "$end_time - $start_time" | bc -l 2>/dev/null || echo "N/A")
    
    # 分离响应体和状态码
    http_code=$(echo "$response" | tail -n1)
    response_body=$(echo "$response" | head -n -1)
    
    echo "响应时间: ${duration}s"
    
    if [[ $http_code -ge 200 && $http_code -lt 300 ]]; then
        echo -e "${GREEN}✅ 成功 (HTTP $http_code)${NC}"
        if command -v jq &> /dev/null; then
            echo "$response_body" | jq . 2>/dev/null || echo "$response_body"
        else
            echo "$response_body"
        fi
    else
        echo -e "${RED}❌ 失败 (HTTP $http_code)${NC}"
        echo "响应: $response_body"
    fi
    
    return $http_code
}

# 检查 API 服务是否运行
if ! wait_for_service; then
    echo -e "${RED}❌ API 服务未运行或无法访问${NC}"
    echo "请确保服务已启动:"
    echo "  docker-compose up azure-tts-api"
    echo "或者设置正确的 API_BASE 环境变量:"
    echo "  export API_BASE=http://your-server:8080/api"
    exit 1
fi

# 测试计数器
total_tests=0
passed_tests=0

run_test() {
    total_tests=$((total_tests + 1))
    test_endpoint "$@"
    if [ $? -ge 200 ] && [ $? -lt 300 ]; then
        passed_tests=$((passed_tests + 1))
    fi
}

# 1. 健康检查
run_test "GET" "/health" "" "健康检查"

# 2. API 文档
run_test "GET" "/" "" "API 文档"

# 3. 单句文本转语音
run_test "POST" "/tts" '{
    "text": "你好，这是Azure TTS API自动化测试",
    "language": "zh-CN",
    "voice": "zh-CN-XiaoxiaoNeural",
    "format": "wav"
}' "单句文本转语音 (wav格式)"

# 4. 英文语音测试
run_test "POST" "/tts" '{
    "text": "Hello, this is Azure TTS API automated testing",
    "language": "en-US",
    "voice": "en-US-JennyNeural",
    "format": "mp3"
}' "英文语音测试 (mp3格式)"

# 5. 批量文本转语音
run_test "POST" "/batch-tts" '{
    "texts": ["第一句自动化测试", "第二句自动化测试", "第三句自动化测试"],
    "language": "zh-CN",
    "voice": "zh-CN-XiaoxiaoNeural",
    "format": "wav"
}' "批量文本转语音"

# 等待一下让文件生成完成
sleep 2

# 6. 获取文件列表
run_test "GET" "/files?limit=10" "" "获取音频文件列表"

# 7. 测试错误情况
echo -e "\n${BLUE}🔍 错误处理测试${NC}"

run_test "POST" "/tts" '{
    "text": ""
}' "空文本测试（应该失败）"

run_test "POST" "/tts" '{
    "invalid": "data"
}' "无效请求数据（应该失败）"

run_test "GET" "/download/nonexistent.wav" "" "下载不存在的文件（应该失败）"

# 8. 性能测试
echo -e "\n${BLUE}⚡ 性能测试${NC}"

echo -e "${YELLOW}测试长文本处理...${NC}"
long_text="这是一个很长的文本用于测试Azure TTS API的性能表现。我们需要确保即使是较长的文本也能正确处理并生成高质量的语音文件。Azure认知服务提供了强大的文本转语音功能，支持多种语言和语音风格。"

run_test "POST" "/tts" "{
    \"text\": \"$long_text\",
    \"language\": \"zh-CN\",
    \"voice\": \"zh-CN-XiaoxiaoNeural\",
    \"format\": \"wav\"
}" "长文本性能测试"

# 9. 下载测试（如果有文件）
echo -e "\n${BLUE}📁 文件下载测试${NC}"
files_response=$(curl -s "$API_BASE/files?limit=1")
filename=$(echo "$files_response" | jq -r '.files[0].filename' 2>/dev/null)

if [[ "$filename" != "null" && "$filename" != "" && "$filename" != "undefined" ]]; then
    echo "尝试下载文件: $filename"
    download_status=$(curl -s -o "/tmp/test_download_$$" -w "%{http_code}" "$API_BASE/download/$filename")
    if [[ $download_status == "200" ]]; then
        file_size=$(stat -f%z "/tmp/test_download_$$" 2>/dev/null || stat -c%s "/tmp/test_download_$$" 2>/dev/null || echo "unknown")
        echo -e "${GREEN}✅ 文件下载成功 (大小: ${file_size} bytes)${NC}"
        rm -f "/tmp/test_download_$$"
        passed_tests=$((passed_tests + 1))
    else
        echo -e "${RED}❌ 文件下载失败 (HTTP $download_status)${NC}"
    fi
    total_tests=$((total_tests + 1))
else
    echo -e "${YELLOW}⚠️  没有可下载的文件，跳过下载测试${NC}"
fi

# 10. 并发测试
echo -e "\n${BLUE}🚀 并发测试${NC}"
echo "启动 3 个并发请求..."

pids=()
for i in {1..3}; do
    (
        curl -s -X POST -H "Content-Type: application/json" \
        -d "{\"text\":\"并发测试第${i}句\",\"language\":\"zh-CN\",\"voice\":\"zh-CN-XiaoxiaoNeural\"}" \
        "$API_BASE/tts" > "/tmp/concurrent_test_${i}_$$" &
        echo $! 
    ) &
    pids+=($!)
done

# 等待所有并发请求完成
for pid in "${pids[@]}"; do
    wait $pid
done

echo -e "${GREEN}✅ 并发测试完成${NC}"
rm -f /tmp/concurrent_test_*_$$

# 测试结果统计
echo -e "\n${GREEN}🎉 测试完成！${NC}"
echo "========================="
echo -e "总测试数: ${BLUE}$total_tests${NC}"
echo -e "通过测试: ${GREEN}$passed_tests${NC}"
echo -e "失败测试: ${RED}$((total_tests - passed_tests))${NC}"

if [ $passed_tests -eq $total_tests ]; then
    echo -e "\n${GREEN}🎊 所有测试通过！API 运行正常${NC}"
else
    echo -e "\n${YELLOW}⚠️  部分测试失败，请检查 API 配置${NC}"
fi
echo -e "\n${YELLOW}💡 使用提示:${NC}"
echo "- 设置自定义 API 地址: export API_BASE=http://your-server:8080/api"
echo "- API 文档: ${API_BASE%/api}/"
echo "- 健康检查: $API_BASE/health"
echo "- 文件列表: $API_BASE/files"
echo ""
echo -e "${BLUE}🔧 常用命令:${NC}"
echo "- 启动服务: docker-compose up azure-tts-api"
echo "- 查看日志: docker-compose logs -f azure-tts-api"
echo "- 停止服务: docker-compose down"
echo ""
echo -e "${GREEN}📝 测试示例:${NC}"
echo "curl -X POST -H 'Content-Type: application/json' \\"
echo "  -d '{\"text\":\"测试文本\",\"language\":\"zh-CN\"}' \\"
echo "  $API_BASE/tts"