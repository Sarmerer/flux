// src/api/wsClient.ts
import mitt from 'mitt'

export type ServerMessage =
  | { type: 'welcome'; userId: string }
  | { type: 'user_joined'; roomId: string; userId: string }
  | { type: 'user_left'; roomId: string; userId: string }
  | { type: 'project_created'; project: { id: string; name: string } }
  | { type: 'project_updated'; project: { id: string; name: string } }
  | { type: 'project_deleted'; projectId: string }
  | { type: 'error'; message: string }

export type ClientMessage =
  | { type: 'join_room'; roomId: string }
  | { type: 'leave_room'; roomId: string }
  | { type: 'create_project'; name: string; description?: string }

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
    const token = localStorage.getItem('token')
    const url = token ? `${WS_URL}?token=${token}` : WS_URL

    this.socket = new WebSocket(url)
    this.connectionPromise = new Promise((resolve, reject) => {
      this.socket!.onopen = () => resolve()
      this.socket!.onerror = (err) => reject(err)
      this.socket!.onclose = () => setTimeout(() => this.connect(), 3000)
      this.socket!.onmessage = (e) => {
        try {
          const data: ServerMessage = JSON.parse(e.data)
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
