import { ref } from 'vue'

export interface WebSocketMessage {
  type: string
  data: any
  timestamp: string
}

class WebSocketService {
  private ws: WebSocket | null = null
  private reconnectAttempts = 0
  private maxReconnectAttempts = 5
  private reconnectInterval = 3000
  private isConnecting = false

  public isConnected = ref(false)
  public messages = ref<WebSocketMessage[]>([])

  connect(url: string) {
    if (this.isConnecting || this.isConnected.value) return

    this.isConnecting = true
    this.ws = new WebSocket(url)

    this.ws.onopen = () => {
      console.log('WebSocket connected')
      this.isConnected.value = true
      this.isConnecting = false
      this.reconnectAttempts = 0
    }

    this.ws.onmessage = (event) => {
      try {
        const message: WebSocketMessage = JSON.parse(event.data)
        this.messages.value.push(message)
        
        // Keep only last 100 messages
        if (this.messages.value.length > 100) {
          this.messages.value = this.messages.value.slice(-100)
        }
      } catch (error) {
        console.error('Failed to parse WebSocket message:', error)
      }
    }

    this.ws.onclose = () => {
      console.log('WebSocket disconnected')
      this.isConnected.value = false
      this.isConnecting = false
      this.attemptReconnect(url)
    }

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error)
      this.isConnected.value = false
      this.isConnecting = false
    }
  }

  private attemptReconnect(url: string) {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.log('Max reconnection attempts reached')
      return
    }

    this.reconnectAttempts++
    console.log(`Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts})...`)

    setTimeout(() => {
      this.connect(url)
    }, this.reconnectInterval)
  }

  disconnect() {
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
    this.isConnected.value = false
    this.reconnectAttempts = 0
  }

  send(message: any) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message))
    } else {
      console.warn('WebSocket is not connected')
    }
  }

  subscribe(eventType: string, callback: (data: any) => void) {
    // Simple subscription mechanism
    // In a real app, you might want to use a more sophisticated event system
    const unsubscribe = () => {
      // Remove callback logic here
    }
    return unsubscribe
  }
}

export const websocketService = new WebSocketService()
