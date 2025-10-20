import { ref, onUnmounted, computed } from 'vue'

interface SubscriptionOptions {
  scope: 'project' | 'table' | 'workflow' | 'global'
  resourceId?: string
  events?: string[]
  onMessage: (data: any) => void
  batchUpdates?: boolean
  throttle?: number
}

interface Subscription {
  id: string
  options: SubscriptionOptions
  unsubscribe: () => void
}

class SmartWebSocketManager {
  private ws: WebSocket | null = null
  private subscriptions = new Map<string, Subscription>()
  private reconnectAttempts = 0
  private maxReconnectAttempts = 5
  private reconnectDelay = 1000
  private messageQueue: any[] = []
  private isConnecting = false
  private heartbeatInterval: number | null = null
  private url: string

  constructor(url: string) {
    this.url = url
  }

  connect(token?: string) {
    if (this.ws?.readyState === WebSocket.OPEN || this.isConnecting) {
      return
    }

    this.isConnecting = true
    const wsUrl = token ? `${this.url}?token=${token}` : this.url

    try {
      this.ws = new WebSocket(wsUrl)

      this.ws.onopen = () => {
        console.log('[WebSocket] Connected')
        this.isConnecting = false
        this.reconnectAttempts = 0
        this.reconnectDelay = 1000
        this.flushMessageQueue()
        this.startHeartbeat()

        // Resubscribe to all active subscriptions
        this.subscriptions.forEach((sub) => {
          this.sendSubscribeMessage(sub.options)
        })
      }

      this.ws.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data)
          this.handleMessage(message)
        } catch (error) {
          console.error('[WebSocket] Failed to parse message:', error)
        }
      }

      this.ws.onerror = (error) => {
        console.error('[WebSocket] Error:', error)
      }

      this.ws.onclose = () => {
        console.log('[WebSocket] Disconnected')
        this.isConnecting = false
        this.stopHeartbeat()
        this.scheduleReconnect()
      }
    } catch (error) {
      console.error('[WebSocket] Connection failed:', error)
      this.isConnecting = false
      this.scheduleReconnect()
    }
  }

  private startHeartbeat() {
    this.stopHeartbeat()
    this.heartbeatInterval = window.setInterval(() => {
      if (this.ws?.readyState === WebSocket.OPEN) {
        this.send({ type: 'ping' })
      }
    }, 30000) // 30 seconds
  }

  private stopHeartbeat() {
    if (this.heartbeatInterval) {
      clearInterval(this.heartbeatInterval)
      this.heartbeatInterval = null
    }
  }

  private scheduleReconnect() {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.warn('[WebSocket] Max reconnection attempts reached. Call connect() manually to retry.')
      return
    }

    this.reconnectAttempts++
    const delay = Math.min(this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1), 30000) // Cap at 30s

    console.log(`[WebSocket] Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts})`)

    setTimeout(() => {
      const token = localStorage.getItem('auth_token')
      this.connect(token || undefined)
    }, delay)
  }

  resetReconnectAttempts() {
    this.reconnectAttempts = 0
    this.reconnectDelay = 1000
  }

  private handleMessage(message: any) {
    // Route message to appropriate subscriptions
    const { type, scope, resourceId, data } = message

    this.subscriptions.forEach((sub) => {
      const { options } = sub

      // Check if subscription matches message
      if (options.scope !== scope && options.scope !== 'global') {
        return
      }

      if (options.resourceId && options.resourceId !== resourceId) {
        return
      }

      if (options.events && !options.events.includes(type)) {
        return
      }

      // Handle throttling
      if (options.throttle) {
        this.throttledCallback(sub.id, () => options.onMessage(data), options.throttle)
      } else {
        options.onMessage(data)
      }
    })
  }

  private throttleTimers = new Map<string, number>()

  private throttledCallback(id: string, callback: () => void, delay: number) {
    const existingTimer = this.throttleTimers.get(id)
    if (existingTimer) {
      clearTimeout(existingTimer)
    }

    const timer = window.setTimeout(() => {
      callback()
      this.throttleTimers.delete(id)
    }, delay)

    this.throttleTimers.set(id, timer)
  }

  subscribe(options: SubscriptionOptions): () => void {
    const id = this.generateSubscriptionId(options)

    // Don't duplicate subscriptions
    if (this.subscriptions.has(id)) {
      return this.subscriptions.get(id)!.unsubscribe
    }

    const unsubscribe = () => {
      this.subscriptions.delete(id)
      this.sendUnsubscribeMessage(options)
    }

    const subscription: Subscription = {
      id,
      options,
      unsubscribe,
    }

    this.subscriptions.set(id, subscription)

    // Send subscribe message if connected
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.sendSubscribeMessage(options)
    }

    return unsubscribe
  }

  private generateSubscriptionId(options: SubscriptionOptions): string {
    return `${options.scope}-${options.resourceId || 'all'}-${options.events?.join(',') || 'all'}`
  }

  private sendSubscribeMessage(options: SubscriptionOptions) {
    this.send({
      type: 'subscribe',
      scope: options.scope,
      resourceId: options.resourceId,
      events: options.events,
    })
  }

  private sendUnsubscribeMessage(options: SubscriptionOptions) {
    this.send({
      type: 'unsubscribe',
      scope: options.scope,
      resourceId: options.resourceId,
      events: options.events,
    })
  }

  private send(data: any) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(data))
    } else {
      this.messageQueue.push(data)
    }
  }

  private flushMessageQueue() {
    while (this.messageQueue.length > 0) {
      const message = this.messageQueue.shift()
      this.send(message)
    }
  }

  disconnect() {
    this.stopHeartbeat()
    this.subscriptions.clear()
    this.throttleTimers.forEach((timer) => clearTimeout(timer))
    this.throttleTimers.clear()

    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
  }

  isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN
  }
}

// Singleton instance
let wsManager: SmartWebSocketManager | null = null

export function useWebSocket(autoConnect = false) {
  // Initialize manager if needed
  if (!wsManager) {
    const wsUrl = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws'
    wsManager = new SmartWebSocketManager(wsUrl)
  }

  const isConnected = ref(false)
  const subscriptions = ref<(() => void)[]>([])

  // Connect on mount if authenticated and autoConnect is enabled
  if (autoConnect) {
    const token = localStorage.getItem('auth_token')
    if (token && !wsManager.isConnected()) {
      wsManager.connect(token)
    }
  }

  // Check connection status
  const checkConnection = () => {
    isConnected.value = wsManager?.isConnected() || false
  }
  const connectionInterval = setInterval(checkConnection, 1000)

  const subscribe = (options: SubscriptionOptions) => {
    if (!wsManager) return () => {}

    const unsubscribe = wsManager.subscribe(options)
    subscriptions.value.push(unsubscribe)
    return unsubscribe
  }

  const connect = () => {
    const token = localStorage.getItem('auth_token')
    wsManager?.connect(token || undefined)
  }

  const disconnect = () => {
    wsManager?.disconnect()
  }

  // Cleanup on unmount
  onUnmounted(() => {
    clearInterval(connectionInterval)
    subscriptions.value.forEach((unsub) => unsub())
    subscriptions.value = []
  })

  return {
    subscribe,
    connect,
    disconnect,
    isConnected: computed(() => isConnected.value),
  }
}
