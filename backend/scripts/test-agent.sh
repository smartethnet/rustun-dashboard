#!/bin/bash

# AI Agent 测试脚本
# 用于快速验证Agent功能是否正常工作

set -e

# 配置
API_URL="${API_URL:-http://localhost:8080}"
USERNAME="${USERNAME:-admin}"
PASSWORD="${PASSWORD:-admin123}"
AGENT_CLI="./agent-cli"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 检查是否设置了OpenAI API Key
if [ -z "$OPENAI_API_KEY" ]; then
    echo -e "${RED}错误: 未设置OPENAI_API_KEY环境变量${NC}"
    echo "请运行: export OPENAI_API_KEY=\"sk-...\""
    exit 1
fi

# 打印标题
print_header() {
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

# 打印测试
print_test() {
    echo -e "\n${YELLOW}📝 测试: $1${NC}"
}

# 打印成功
print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

# 打印失败
print_fail() {
    echo -e "${RED}❌ $1${NC}"
}

# 测试API端点
test_api() {
    local message="$1"
    local description="$2"
    
    print_test "$description"
    echo "发送消息: $message"
    
    response=$(curl -s -u "$USERNAME:$PASSWORD" \
        -H "Content-Type: application/json" \
        -d "{\"message\":\"$message\"}" \
        "$API_URL/api/agent/chat")
    
    # 检查是否有错误
    error=$(echo "$response" | grep -o '"error":"[^"]*"' | cut -d'"' -f4)
    if [ -n "$error" ]; then
        print_fail "错误: $error"
        return 1
    fi
    
    # 提取消息
    agent_message=$(echo "$response" | grep -o '"message":"[^"]*"' | cut -d'"' -f4 | sed 's/\\n/\n/g')
    echo -e "\n🤖 Agent回复:\n$agent_message\n"
    
    print_success "测试通过"
    return 0
}

# 检查服务器健康
check_server() {
    print_header "检查服务器状态"
    
    if curl -s "$API_URL/health" > /dev/null; then
        print_success "服务器运行正常: $API_URL"
    else
        print_fail "无法连接到服务器: $API_URL"
        echo "请确保后端正在运行"
        exit 1
    fi
}

# 构建CLI工具
build_cli() {
    print_header "构建CLI工具"
    
    if [ ! -f "$AGENT_CLI" ]; then
        echo "正在构建agent-cli..."
        go build -o "$AGENT_CLI" cmd/agent-cli/main.go
        print_success "构建完成"
    else
        print_success "CLI工具已存在"
    fi
}

# 运行测试
run_tests() {
    print_header "运行AI Agent功能测试"
    
    # 测试1: 列出集群
    test_api "列出所有集群" "查询集群列表"
    sleep 2
    
    # 测试2: 创建客户端
    test_api "帮我创建一个名为测试NAS的客户端，放在test-cluster集群" "创建新客户端"
    sleep 2
    
    # 测试3: 查询客户端
    test_api "查询test-cluster集群的所有客户端" "查询特定集群的客户端"
    sleep 2
    
    # 测试4: 更复杂的查询
    test_api "一共有多少个客户端？" "统计查询"
    sleep 2
    
    echo ""
    print_success "所有测试完成！"
}

# 交互模式
interactive_mode() {
    print_header "进入交互模式"
    echo "您现在可以直接与AI Agent对话"
    echo "输入 'exit' 退出"
    echo ""
    
    "$AGENT_CLI" -api "$API_URL" -user "$USERNAME" -pass "$PASSWORD"
}

# 主函数
main() {
    echo -e "${GREEN}"
    cat << "EOF"
    ____            __               ___    ____
   / __ \__  _______/ /___  ______   /   |  /  _/
  / /_/ / / / / ___/ __/ / / / __ \ / /| |  / /  
 / _, _/ /_/ (__  ) /_/ /_/ / / / // ___ |_/ /   
/_/ |_|\__,_/____/\__/\__,_/_/ /_//_/  |_/___/   
                                                  
        AI Agent 测试工具
EOF
    echo -e "${NC}"
    
    check_server
    build_cli
    
    # 根据参数决定运行模式
    if [ "$1" == "--interactive" ] || [ "$1" == "-i" ]; then
        interactive_mode
    else
        run_tests
        
        echo ""
        echo -e "${YELLOW}提示: 运行 $0 --interactive 进入交互模式${NC}"
    fi
}

# 运行
main "$@"

