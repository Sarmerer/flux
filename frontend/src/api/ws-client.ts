// src/api/wsClient.ts
import mitt from 'mitt'

export type ServerMessage =
  | { type: 'progress'; id: string; data: { operation_id: string; step: string; progress: number; message: string } }
  | { type: 'error'; id?: string; data: { operation_id?: string; message: string; error?: string } }
  | { type: 'success'; id?: string; data: any }
  | { type: 'notification'; data: { message: string; data?: any } }
  | { type: 'ping' }

export type ClientMessage =
  | { type: 'pong' }

export type ServerEventMap = {
  [K in ServerMessage['type']]: Extract<ServerMessage, { type: K }>
}

export type ClientEventMap = {
  [K in ClientMessage['type']]: Extract<ClientMessage, { type: K }>
}

class WebSocketClient {
  private socket: WebSocket | null = null
  private emitter = mitt<ServerEventMap>()
  private connectionPromise: Promise<void> | null = null

  async connect() {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) return

    const WS_URL = import.meta.env.VITE_WS_URL
    this.socket = new WebSocket(WS_URL)
    this.connectionPromise = new Promise((resolve, reject) => {
      this.socket!.onopen = () => {
        // Send auth message immediately after connection opens
        const token = localStorage.getItem('token')
        if (token) {
          const authMsg = { type: 'auth', token: `Bearer ${token}` }
          this.socket!.send(JSON.stringify(authMsg))
        }
        resolve()
      }
      this.socket!.onerror = (err) => reject(err)
      this.socket!.onclose = () => setTimeout(() => this.connect(), 3000)
      this.socket!.onmessage = (e) => {
        try {
          const data: ServerMessage = JSON.parse(e.data)
          // Respond to ping
          if ((data as any).type === 'ping') {
            this.socket!.send(JSON.stringify({ type: 'pong' }))
            return
          }
          this.emitter.emit(data.type, data)
        } catch {
          console.error('Invalid WS message:', e.data)
        }
      }
    })

    return this.connectionPromise
  }

  async send<T extends ClientMessage>(msg: T) {
    await this.connect()
    this.socket!.send(JSON.stringify(msg))
  }

  on<K extends keyof ServerEventMap>(type: K, handler: (event: ServerEventMap[K]) => void) {
    this.emitter.on(type, handler)
    return () => this.emitter.off(type, handler)
  }
}

export const ws = new WebSocketClient()
