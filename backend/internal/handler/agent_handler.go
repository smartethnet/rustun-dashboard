package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/smartethnet/rustun-dashboard/internal/agent"
	"github.com/smartethnet/rustun-dashboard/internal/model"

	"github.com/gin-gonic/gin"
)

// AgentHandler handles AI agent requests
type AgentHandler struct {
	agent *agent.Agent
}

// NewAgentHandler creates a new agent handler
func NewAgentHandler(agent *agent.Agent) *AgentHandler {
	return &AgentHandler{
		agent: agent,
	}
}

// Chat handles chat requests
// @Summary Chat with AI agent
// @Description Send a message to the AI agent and get a response
// @Tags agent
// @Accept json
// @Produce json
// @Param request body agent.AgentChatRequest true "Chat request"
// @Success 200 {object} agent.AgentChatResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/agent/chat [post]
func (h *AgentHandler) Chat(c *gin.Context) {
	var req agent.AgentChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponseWithCode(
			http.StatusBadRequest,
			"Invalid request",
			err.Error(),
		))
		return
	}

	resp, err := h.agent.Chat(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponseWithCode(
			http.StatusInternalServerError,
			"Agent error",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, resp)
}

// StreamChat handles streaming chat requests
// @Summary Stream chat with AI agent
// @Description Send a message and receive streaming response
// @Tags agent
// @Accept json
// @Produce text/event-stream
// @Param request body agent.AgentChatRequest true "Chat request"
// @Success 200 {string} string "SSE stream"
// @Failure 400 {object} model.ErrorResponse
// @Router /api/agent/chat/stream [post]
func (h *AgentHandler) StreamChat(c *gin.Context) {
	var req agent.AgentChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponseWithCode(
			http.StatusBadRequest,
			"Invalid request",
			err.Error(),
		))
		return
	}

	// Set headers for SSE
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no") // Disable nginx buffering

	// Create context for streaming
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	// Create a channel for streaming responses
	streamChan := make(chan agent.StreamEvent, 10)

	// Start agent chat in goroutine
	go func() {
		defer close(streamChan)

		h.agent.ChatStream(ctx, req, streamChan)
	}()

	// Send SSE events
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, model.ErrorResponseWithCode(
			http.StatusInternalServerError,
			"Streaming not supported",
			"",
		))
		return
	}

	for event := range streamChan {
		data, err := json.Marshal(event)
		if err != nil {
			continue
		}

		fmt.Fprintf(c.Writer, "data: %s\n\n", data)
		flusher.Flush()

		// Check if client disconnected
		if c.Writer.Status() == http.StatusRequestTimeout {
			break
		}
	}
}
