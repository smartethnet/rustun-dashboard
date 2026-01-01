import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

// Create axios instance with basic auth
const createAuthAxios = () => {
  const username = localStorage.getItem('username')
  const password = localStorage.getItem('password')
  
  return axios.create({
    baseURL: API_BASE_URL,
    auth: {
      username,
      password
    }
  })
}

// Agent API
export const agentApi = {
  // Send chat message
  chat: async (message, history = []) => {
    const api = createAuthAxios()
    const response = await api.post('/api/agent/chat', {
      message,
      history
    })
    return response.data
  },

  // Send chat message with streaming
  chatStream: async (message, history = [], onChunk) => {
    const username = localStorage.getItem('username')
    const password = localStorage.getItem('password')
    const auth = 'Basic ' + btoa(username + ':' + password)

    const response = await fetch(`${API_BASE_URL}/api/agent/chat/stream`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': auth,
      },
      body: JSON.stringify({
        message,
        history
      }),
    })

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }

    const reader = response.body.getReader()
    const decoder = new TextDecoder()

    let fullMessage = ''
    let toolCalls = []

    while (true) {
      const { done, value } = await reader.read()
      
      if (done) {
        break
      }

      const chunk = decoder.decode(value)
      const lines = chunk.split('\n')

      for (const line of lines) {
        if (line.startsWith('data: ')) {
          const data = line.slice(6)
          
          try {
            const event = JSON.parse(data)
            
            if (event.type === 'content') {
              fullMessage += event.content
              onChunk({ type: 'content', content: event.content, fullMessage })
            } else if (event.type === 'tool_call') {
              toolCalls.push(event.tool_call)
              onChunk({ type: 'tool_call', toolCall: event.tool_call })
            } else if (event.type === 'done') {
              onChunk({ type: 'done', fullMessage, toolCalls })
            } else if (event.type === 'error') {
              throw new Error(event.error)
            }
          } catch (e) {
            if (e instanceof SyntaxError) {
              // Ignore JSON parse errors for incomplete chunks
              continue
            }
            throw e
          }
        }
      }
    }

    return { message: fullMessage, tool_calls: toolCalls }
  }
}

