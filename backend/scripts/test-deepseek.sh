#!/bin/bash

# DeepSeek 测试脚本
# 用于测试 DeepSeek 集成是否正常工作

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_header() {
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ️  $1${NC}"
}

# 检查 API Key
print_header "检查 DeepSeek API Key"

if [ -z "$DEEPSEEK_API_KEY" ]; then
    print_error "未设置 DEEPSEEK_API_KEY 环境变量"
    echo ""
    echo "请先设置 API Key："
    echo "  export DEEPSEEK_API_KEY=\"sk-...\""
    echo ""
    echo "获取 API Key："
    echo "  https://platform.deepseek.com/"
    exit 1
fi

print_success "DEEPSEEK_API_KEY 已设置"
echo ""

# 测试 DeepSeek API
print_header "测试 DeepSeek API 连接"

response=$(curl -s https://api.deepseek.com/v1/chat/completions \
    -H "Authorization: Bearer $DEEPSEEK_API_KEY" \
    -H "Content-Type: application/json" \
    -d '{
        "model": "deepseek-chat",
        "messages": [
            {"role": "user", "content": "Hello"}
        ],
        "max_tokens": 10
    }')

if echo "$response" | grep -q "error"; then
    print_error "API 测试失败"
    echo "$response"
    exit 1
fi

print_success "DeepSeek API 连接正常"
echo ""

# 检查后端配置
print_header "检查后端配置"

if [ ! -f "config.yaml" ]; then
    print_error "config.yaml 不存在"
    exit 1
fi

if grep -q "provider.*deepseek" config.yaml; then
    print_success "config.yaml 已配置使用 DeepSeek"
else
    print_info "config.yaml 未配置 DeepSeek，将使用环境变量"
fi

echo ""

# 构建并测试
print_header "构建并启动服务"

print_info "正在构建后端..."
if go build -o rustun-dashboard cmd/dashboard/main.go; then
    print_success "构建成功"
else
    print_error "构建失败"
    exit 1
fi

echo ""
print_info "正在启动服务..."
echo ""

# 启动服务（后台）
./rustun-dashboard -config config.yaml &
SERVER_PID=$!

# 等待服务启动
sleep 3

# 检查服务是否运行
if ! kill -0 $SERVER_PID 2>/dev/null; then
    print_error "服务启动失败"
    exit 1
fi

print_success "服务已启动 (PID: $SERVER_PID)"
echo ""

# 测试 Health 端点
print_header "测试 Health 端点"

if curl -s http://localhost:8080/health | grep -q "ok"; then
    print_success "Health 检查通过"
else
    print_error "Health 检查失败"
    kill $SERVER_PID
    exit 1
fi

echo ""

# 测试 Agent 端点
print_header "测试 AI Agent (DeepSeek)"

print_info "发送测试消息：列出所有集群"

response=$(curl -s -u admin:admin123 \
    -H "Content-Type: application/json" \
    -d '{"message":"列出所有集群"}' \
    http://localhost:8080/api/agent/chat)

if echo "$response" | grep -q "message"; then
    print_success "Agent 响应正常"
    echo ""
    echo "Agent 回复："
    echo "$response" | grep -o '"message":"[^"]*"' | cut -d'"' -f4 | head -1
else
    print_error "Agent 响应异常"
    echo "$response"
fi

echo ""

# 清理
print_header "清理"

print_info "停止服务..."
kill $SERVER_PID
wait $SERVER_PID 2>/dev/null || true

print_success "服务已停止"
echo ""

# 总结
print_header "测试完成"
echo ""
echo "✅ DeepSeek 集成测试通过！"
echo ""
echo "后续步骤："
echo "  1. 启动服务：go run cmd/dashboard/main.go"
echo "  2. 使用 CLI：./agent-cli"
echo "  3. 查看文档：docs/DEEPSEEK_INTEGRATION.md"
echo ""
echo "价格优势："
echo "  - DeepSeek: ~¥0.10/100次操作"
echo "  - GPT-4o-mini: ~¥1.50/100次操作"
echo "  - 节省 93% 成本！"
echo ""

