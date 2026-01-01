package agent

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/smartethnet/rustun-dashboard/internal/service"
)

//go:embed knowledge.md
var productKnowledge string

// Agent is the main AI agent that orchestrates conversations
type Agent struct {
	llmClient    *LLMClient
	toolExecutor *ToolExecutor
	systemPrompt string
}

// NewAgent creates a new agent instance
func NewAgent(apiKey, model, baseURL, provider string, routeService *service.RouteService) *Agent {
	return &Agent{
		llmClient:    NewLLMClient(apiKey, model, baseURL, provider),
		toolExecutor: NewToolExecutor(routeService),
		systemPrompt: getSystemPrompt(),
	}
}

// AgentChatRequest represents a chat request from the user
type AgentChatRequest struct {
	Message             string        `json:"message" binding:"required"`
	ConversationID      string        `json:"conversation_id,omitempty"`
	ConversationHistory []ChatMessage `json:"history,omitempty"`
}

// AgentChatResponse represents the agent's response
type AgentChatResponse struct {
	Message        string        `json:"message"`
	ConversationID string        `json:"conversation_id,omitempty"`
	ToolCalls      []ToolCallLog `json:"tool_calls,omitempty"`
	Error          string        `json:"error,omitempty"`
}

// ToolCallLog logs what tools were called
type ToolCallLog struct {
	Tool      string `json:"tool"`
	Arguments string `json:"arguments"`
	Result    string `json:"result"`
}

// StreamEvent represents a single event in the SSE stream
type StreamEvent struct {
	Type        string        `json:"type"` // "content", "tool_call", "done", "error"
	Content     string        `json:"content,omitempty"`
	ToolCall    *ToolCallLog  `json:"tool_call,omitempty"`
	FullMessage string        `json:"fullMessage,omitempty"` // Only for "done" event
	ToolCalls   []ToolCallLog `json:"toolCalls,omitempty"`   // Only for "done" event
	Error       string        `json:"error,omitempty"`
}

// Chat processes a user message and returns a response
func (a *Agent) Chat(req AgentChatRequest) (*AgentChatResponse, error) {
	// Initialize conversation with system prompt and history
	messages := []ChatMessage{
		{
			Role:    "system",
			Content: a.systemPrompt,
		},
	}

	// Add conversation history if provided
	if len(req.ConversationHistory) > 0 {
		messages = append(messages, req.ConversationHistory...)
	}

	// Add user message
	messages = append(messages, ChatMessage{
		Role:    "user",
		Content: req.Message,
	})

	// Get available tools
	tools := a.toolExecutor.GetTools()
	log.Printf("[Agent] Chat: Registered %d tools", len(tools))

	// Track tool calls for response
	var toolCallLogs []ToolCallLog

	// Maximum iterations to prevent infinite loops
	maxIterations := 10
	for i := 0; i < maxIterations; i++ {
		log.Printf("[Agent] Chat: Iteration %d/%d", i+1, maxIterations)

		// Call LLM
		resp, err := a.llmClient.Chat(messages, tools)
		if err != nil {
			return nil, fmt.Errorf("LLM调用失败: %w", err)
		}

		if len(resp.Choices) == 0 {
			return nil, fmt.Errorf("LLM未返回任何响应")
		}

		choice := resp.Choices[0]
		assistantMsg := choice.Message

		log.Printf("[Agent] Chat: Got response, tool_calls count: %d, content length: %d",
			len(assistantMsg.ToolCalls), len(assistantMsg.Content))

		// Add assistant message to conversation
		messages = append(messages, assistantMsg)

		// If no tool calls, we're done
		if len(assistantMsg.ToolCalls) == 0 {
			log.Printf("[Agent] Chat: No tool calls, returning final response")
			return &AgentChatResponse{
				Message:   assistantMsg.Content,
				ToolCalls: toolCallLogs,
			}, nil
		}

		// Execute tool calls
		log.Printf("[Agent] Chat: Executing %d tool call(s)", len(assistantMsg.ToolCalls))
		for _, toolCall := range assistantMsg.ToolCalls {
			if toolCall.Type != "function" {
				log.Printf("[Agent] Chat: Skipping non-function tool call: %s", toolCall.Type)
				continue
			}

			functionName := toolCall.Function.Name
			log.Printf("[Agent] Chat: Calling tool: %s", functionName)
			arguments := toolCall.Function.Arguments

			// Execute the tool
			result, err := a.toolExecutor.ExecuteTool(functionName, arguments)
			if err != nil {
				result = fmt.Sprintf(`{"error": "%s"}`, err.Error())
			}

			// Log tool call
			toolCallLogs = append(toolCallLogs, ToolCallLog{
				Tool:      functionName,
				Arguments: arguments,
				Result:    result,
			})

			// Add tool result to conversation
			messages = append(messages, ChatMessage{
				Role:       "tool",
				Content:    result,
				ToolCallID: toolCall.ID,
				Name:       functionName,
			})
		}
	}

	return nil, fmt.Errorf("达到最大迭代次数，对话可能过于复杂")
}

// ChatStream processes a user message and streams the response
func (a *Agent) ChatStream(ctx context.Context, req AgentChatRequest, streamChan chan<- StreamEvent) {
	messages := []ChatMessage{
		{
			Role:    "system",
			Content: a.systemPrompt,
		},
	}

	if len(req.ConversationHistory) > 0 {
		messages = append(messages, req.ConversationHistory...)
	}

	messages = append(messages, ChatMessage{
		Role:    "user",
		Content: req.Message,
	})

	tools := a.toolExecutor.GetTools()
	log.Printf("[Agent] ChatStream: Registered %d tools", len(tools))

	var (
		fullContent  strings.Builder
		toolCallLogs []ToolCallLog
		wg           sync.WaitGroup
	)

	maxIterations := 10
	for i := 0; i < maxIterations; i++ {
		select {
		case <-ctx.Done():
			streamChan <- StreamEvent{Type: "error", Error: "Stream cancelled by client"}
			return
		default:
			// Continue
		}

		log.Printf("[Agent] ChatStream: Iteration %d/%d", i+1, maxIterations)

		llmStream, err := a.llmClient.ChatStream(ctx, messages, tools)
		if err != nil {
			streamChan <- StreamEvent{Type: "error", Error: fmt.Sprintf("LLM stream failed: %v", err)}
			return
		}

		var (
			assistantMsg ChatMessage
			isToolCall   bool
		)

		for streamResp := range llmStream {
			if streamResp.Error != nil {
				streamChan <- StreamEvent{Type: "error", Error: streamResp.Error.Message}
				return
			}

			choice := streamResp.Choices[0]
			delta := choice.Delta

			if delta.Role != "" {
				assistantMsg.Role = delta.Role
			}
			if delta.Content != "" {
				fullContent.WriteString(delta.Content)
				streamChan <- StreamEvent{Type: "content", Content: delta.Content}
			}

			if len(delta.ToolCalls) > 0 {
				isToolCall = true
				// Accumulate tool calls
				if len(assistantMsg.ToolCalls) == 0 {
					assistantMsg.ToolCalls = make([]ToolCall, len(delta.ToolCalls))
				}
				for i, newToolCall := range delta.ToolCalls {
					if newToolCall.ID != "" {
						assistantMsg.ToolCalls[i].ID = newToolCall.ID
					}
					if newToolCall.Type != "" {
						assistantMsg.ToolCalls[i].Type = newToolCall.Type
					}
					if newToolCall.Function.Name != "" {
						assistantMsg.ToolCalls[i].Function.Name = newToolCall.Function.Name
					}
					assistantMsg.ToolCalls[i].Function.Arguments += newToolCall.Function.Arguments
				}
			}
		}

		log.Printf("[Agent] ChatStream: Got response, tool_calls: %d, isToolCall: %v",
			len(assistantMsg.ToolCalls), isToolCall)

		messages = append(messages, assistantMsg)

		if !isToolCall {
			// If no tool calls, we're done with the LLM part
			streamChan <- StreamEvent{
				Type:        "done",
				FullMessage: fullContent.String(),
				ToolCalls:   toolCallLogs,
			}
			return
		}

		// Execute tool calls in parallel
		log.Printf("[Agent] ChatStream: Executing %d tool calls", len(assistantMsg.ToolCalls))
		for _, toolCall := range assistantMsg.ToolCalls {
			wg.Add(1)
			go func(tc ToolCall) {
				defer wg.Done()
				functionName := tc.Function.Name
				arguments := tc.Function.Arguments

				log.Printf("[Agent] ChatStream: Calling tool: %s with args: %s", functionName, arguments)

				result, err := a.toolExecutor.ExecuteTool(functionName, arguments)
				if err != nil {
					result = fmt.Sprintf(`{"error": "%s"}`, err.Error())
				}

				logEntry := ToolCallLog{
					Tool:      functionName,
					Arguments: arguments,
					Result:    result,
				}
				toolCallLogs = append(toolCallLogs, logEntry)

				// Send tool_call event immediately after execution
				streamChan <- StreamEvent{Type: "tool_call", ToolCall: &logEntry}

				messages = append(messages, ChatMessage{
					Role:       "tool",
					Content:    result,
					ToolCallID: tc.ID,
					Name:       functionName,
				})
			}(toolCall)
		}
		wg.Wait() // Wait for all tool calls to finish before next LLM call

		// Reset fullContent for the next LLM response
		fullContent.Reset()
	}

	streamChan <- StreamEvent{Type: "error", Error: "Maximum iterations reached, conversation might be too complex"}
}

func getSystemPrompt() string {
	// Use English for the system prompt to avoid language bias
	// Knowledge base remains in Chinese as product content
	return strings.TrimSpace(`
You are Rustun VPN's AI Smart Assistant. Rustun is an open-source VPN tunnel built with Rust, featuring P2P direct connection, intelligent routing, and multi-tenant isolation.

## Language Strategy

**CRITICAL**: Automatically match the user's language:
- User asks in Chinese → Respond in Chinese
- User asks in English → Respond in English
- User asks in any language → Respond in that language
- Maintain the same language style and professionalism

## Your Roles

You are both an **Operations Assistant** and a **Technical Advisor**:
- 🔧 Operations: Manage clusters, clients, routing configs
- 📖 Technical: Answer architecture, features, deployment questions
- 💡 Solutions: Recommend best practices and configurations

## Core Capabilities

### 1. Operations Management (via Function Calling)
- Query/manage Clusters
- Query/manage Clients
- Create/update/delete client configs
- Manage routing rules (Routes/CIDRs)

### 2. Technical Consulting (based on Knowledge Base)
- Product features and capabilities
- Architecture (P2P, relay, encryption, NAT traversal)
- Installation and deployment guidance
- Configuration details (server.toml, routes.json, CLI params)
- Troubleshooting and optimization
- Multi-platform support (Linux, macOS, Windows, iOS, Android)
- Multi-tenant isolation

### 3. Scenario Solutions
- Remote work, site-to-site networking
- Intranet penetration, service exposure
- Multi-region interconnection
- Mobile device access
- Enterprise multi-tenant isolation

## Working Principles

1. **Intelligent Understanding**: Understand user needs naturally, distinguish operations vs consulting, ask for clarification when needed

2. **Knowledge Guidance**: Answer accurately based on knowledge base, explain with practical scenarios, maintain accuracy, never fabricate

3. **Multi-turn Dialogue**: Ask friendly questions when parameters missing, guide step-by-step for complex topics, remember context

4. **Tool Usage**: Use Function Calling for operations, confirm all required params before execution, report results clearly

5. **Friendly Feedback**: Use concise, professional, friendly language; explain technical terms appropriately; use Markdown formatting; highlight success, provide solutions on failure

6. **Proactive Suggestions**: Recommend best practices, remind of potential issues, suggest optimizations

## Parameter Collection Strategy

### create_client
**Required**:
- cluster (auto-created if not exists)

**Recommended**:
- name (friendly client name, strongly suggested)

**Auto-assigned** (no user input needed):
- identity (UUID)
- private_ip (virtual IP, starting from 10.12.0.10)
- mask (subnet mask, default 255.255.0.0)
- gateway (default 10.12.0.1)

**Optional**:
- routes (CIDR rules, can be added anytime)

**Collection rules**:
1. If cluster provided, create directly (name optional)
2. If cluster missing, ask first with naming suggestions
3. Suggest name but don't force

### update_client
**Required**:
- cluster
- identity (UUID)
- at least one field to update (name or routes)

### delete_client
**Required**:
- cluster
- identity (UUID)

**Safety**: Confirm before deletion, warn irreversible

## Response Style Guide

### Operations Responses

**On success**:
- Show client info clearly (name, cluster, IP, identity)
- Provide next steps (e.g., connection commands)
- Ask if routes configuration needed

**On list queries**:
- Use clear tables/lists
- Highlight key info (IP, status, routes count)

**On errors**:
- Explain the reason
- Provide solutions and suggestions
- Give correct examples

### Technical Consulting Responses

**On features**:
- Explain in layers (concept, principle, use cases)
- Provide related config options
- Give practical examples

**On deployment**:
- Provide specific commands and configs
- Step-by-step instructions
- Remind of important notes

**On troubleshooting**:
- List possible causes
- Provide systematic diagnostic steps
- Give solutions

## Special Scenarios

**Scenario 1**: User asks "How to use Rustun"
- First understand their specific scenario (remote work? intranet penetration? multi-site?)
- Provide targeted deployment plan
- Complete installation, config, startup guide
- Give client connection command examples

**Scenario 2**: User asks technical details (e.g., "How is P2P implemented")
- Explain technical principles accurately based on knowledge base
- Can reference architecture diagrams, protocol docs
- Explain technical terms in plain language
- Provide related configuration options

**Scenario 3**: User encounters problems (e.g., "Connection failed")
- Ask for specific error messages and environment
- Provide systematic diagnostic steps
- Give solutions for common issues from knowledge base
- Suggest checking logs, configs when necessary

**Scenario 4**: User needs solution advice (e.g., "How to deploy for enterprise")
- Understand enterprise scale, network environment, security requirements
- Recommend appropriate architecture (multi-tenant isolation, encryption methods)
- Provide configuration templates and best practices
- Remind of security considerations

---

## Product Knowledge Base (产品知识库)

` + productKnowledge + `

---

Remember: You're not just a tool executor, but a technical partner. Provide value with professional knowledge!
`)
}
