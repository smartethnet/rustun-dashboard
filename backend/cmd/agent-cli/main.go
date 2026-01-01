package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Message string        `json:"message"`
	History []ChatMessage `json:"history,omitempty"`
}

type ChatResponse struct {
	Message   string        `json:"message"`
	ToolCalls []ToolCallLog `json:"tool_calls,omitempty"`
	Error     string        `json:"error,omitempty"`
}

type ToolCallLog struct {
	Tool      string `json:"tool"`
	Arguments string `json:"arguments"`
	Result    string `json:"result"`
}

func main() {
	// Parse command line flags
	apiURL := flag.String("api", "http://localhost:8080", "API server URL")
	username := flag.String("user", "admin", "Username for Basic Auth")
	password := flag.String("pass", "admin123", "Password for Basic Auth")
	message := flag.String("msg", "", "Direct message to send (non-interactive mode)")
	flag.Parse()

	baseURL := *apiURL
	auth := basicAuth(*username, *password)

	// Check server health
	if err := checkHealth(baseURL); err != nil {
		fmt.Printf("⚠️  Warning: Cannot connect to server at %s\n", baseURL)
		fmt.Printf("   Error: %v\n", err)
		fmt.Printf("   Make sure the backend is running.\n\n")
		return
	}

	fmt.Println("🤖 Rustun AI Agent CLI")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Connected to: %s\n", baseURL)
	fmt.Println()

	// Non-interactive mode
	if *message != "" {
		sendMessage(baseURL, auth, *message)
		return
	}

	// Interactive mode
	fmt.Println("💡 Tips:")
	fmt.Println("  - Type your request in natural language")
	fmt.Println("  - The agent remembers the conversation context")
	fmt.Println("  - Type 'clear' to reset conversation history")
	fmt.Println("  - Type 'exit' or 'quit' to leave")
	fmt.Println("  - Type 'help' for examples")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	var conversationHistory []ChatMessage

	for {
		fmt.Print("You: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		// Handle commands
		switch strings.ToLower(input) {
		case "exit", "quit":
			fmt.Println("👋 Goodbye!")
			return
		case "help":
			showHelp()
			continue
		case "clear":
			conversationHistory = nil
			fmt.Println("🔄 Conversation history cleared")
			continue
		}

		// Send message with history
		response := sendMessageWithHistory(baseURL, auth, input, conversationHistory)
		if response != nil {
			// Update conversation history
			conversationHistory = append(conversationHistory,
				ChatMessage{Role: "user", Content: input},
				ChatMessage{Role: "assistant", Content: response.Message},
			)
		}
		fmt.Println()
	}
}

func checkHealth(baseURL string) error {
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	return nil
}

func sendMessageWithHistory(baseURL, auth, message string, history []ChatMessage) *ChatResponse {
	req := ChatRequest{
		Message: message,
		History: history,
	}

	return sendRequest(baseURL, auth, req)
}

func sendMessage(baseURL, auth, message string) {
	req := ChatRequest{
		Message: message,
	}
	sendRequest(baseURL, auth, req)
}

func sendRequest(baseURL, auth string, req ChatRequest) *ChatResponse {

	jsonData, err := json.Marshal(req)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return nil
	}

	httpReq, err := http.NewRequest("POST", baseURL+"/api/agent/chat", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return nil
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", auth)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ Error reading response: %v\n", err)
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ Server error (status %d): %s\n", resp.StatusCode, string(body))
		return nil
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		fmt.Printf("❌ Error parsing response: %v\n", err)
		return nil
	}

	if chatResp.Error != "" {
		fmt.Printf("❌ Error: %s\n", chatResp.Error)
		return nil
	}

	// Display tool calls if any
	if len(chatResp.ToolCalls) > 0 {
		fmt.Println("🔧 Actions taken:")
		for _, tc := range chatResp.ToolCalls {
			fmt.Printf("  └─ %s\n", tc.Tool)
		}
		fmt.Println()
	}

	// Display agent response
	fmt.Printf("🤖 Agent: %s\n", chatResp.Message)
	
	return &chatResp
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return "Basic " + encodeBase64([]byte(auth))
}

func encodeBase64(data []byte) string {
	const base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	var result strings.Builder

	for i := 0; i < len(data); i += 3 {
		b := (data[i] & 0xFC) >> 2
		result.WriteByte(base64Table[b])

		b = (data[i] & 0x03) << 4
		if i+1 < len(data) {
			b |= (data[i+1] & 0xF0) >> 4
			result.WriteByte(base64Table[b])

			b = (data[i+1] & 0x0F) << 2
			if i+2 < len(data) {
				b |= (data[i+2] & 0xC0) >> 6
				result.WriteByte(base64Table[b])
				b = data[i+2] & 0x3F
				result.WriteByte(base64Table[b])
			} else {
				result.WriteByte(base64Table[b])
				result.WriteByte('=')
			}
		} else {
			result.WriteByte(base64Table[b])
			result.WriteString("==")
		}
	}

	return result.String()
}

func showHelp() {
	fmt.Println()
	fmt.Println("📚 Example commands:")
	fmt.Println()
	fmt.Println("  Basic operations:")
	fmt.Println("  • 帮我创建一个NAS客户端")
	fmt.Println("  • 查询production集群有哪些客户端")
	fmt.Println("  • 删除IP为10.12.0.10的客户端")
	fmt.Println("  • 把办公室客户端的名字改成总部")
	fmt.Println("  • 列出所有集群")
	fmt.Println()
	fmt.Println("  Multi-turn conversation (Agent will ask for details):")
	fmt.Println("  You: 帮我创建一个客户端")
	fmt.Println("  Agent: 好的，请问客户端名称和集群是什么？")
	fmt.Println("  You: 名字叫办公室，放在production集群")
	fmt.Println("  Agent: ✅ 已成功创建！")
	fmt.Println()
	fmt.Println("  Commands:")
	fmt.Println("  • clear - Clear conversation history")
	fmt.Println("  • help - Show this help")
	fmt.Println("  • exit/quit - Exit the program")
	fmt.Println()
}

