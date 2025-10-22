#!/bin/bash

# Azure TTS API 重建脚本

echo "🔄 重建 Azure TTS API 容器..."

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# 停止并删除现有容器
echo -e "${YELLOW}停止现有容器...${NC}"
sudo docker compose down

# 重新构建镜像
echo -e "${YELLOW}重新构建镜像...${NC}"
sudo docker compose build --no-cache azure-tts-api

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ 镜像构建成功${NC}"
    
    # 启动容器
    echo -e "${YELLOW}启动容器...${NC}"
    sudo docker compose up -d azure-tts-api
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ 容器启动成功${NC}"
        
        # 等待服务启动
        echo -e "${YELLOW}等待服务启动...${NC}"
        sleep 5
        
        # 检查服务状态
        echo -e "${YELLOW}检查服务状态...${NC}"
        sudo docker compose ps
        
        # 检查日志
        echo -e "\n${YELLOW}最近的日志:${NC}"
        sudo docker compose logs --tail=20 azure-tts-api
        
        # 测试健康检查
        echo -e "\n${YELLOW}测试健康检查...${NC}"
        sleep 3
        if curl -s http://localhost:8080/api/health > /dev/null; then
            echo -e "${GREEN}✅ 服务健康检查通过${NC}"
            echo -e "\n${GREEN}🎉 重建完成！服务已就绪${NC}"
            echo -e "访问: http://localhost:8080/"
            echo -e "健康检查: http://localhost:8080/api/health"
        else
            echo -e "${RED}❌ 健康检查失败${NC}"
            echo -e "\n${YELLOW}查看详细日志:${NC}"
            sudo docker compose logs azure-tts-api
        fi
    else
        echo -e "${RED}❌ 容器启动失败${NC}"
    fi
else
    echo -e "${RED}❌ 镜像构建失败${NC}"
fi

echo -e "\n${YELLOW}💡 有用的命令:${NC}"
echo "查看日志: sudo docker compose logs -f azure-tts-api"
echo "进入容器: sudo docker compose exec azure-tts-api /bin/bash"
echo "停止服务: sudo docker compose down"
echo "运行测试: ./test-api.sh"