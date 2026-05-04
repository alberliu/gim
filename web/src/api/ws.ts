import { create, fromBinary, toBinary } from '@bufbuild/protobuf'
import {
  PacketCommand,
  PacketSchema,
  MessageSchema,
  MessageCommand,
} from '@/gen/proto/connect/connect.ext_pb'
import { auth } from '@/store/auth'

export type WSEventName =
  | 'open'
  | 'close'
  | 'message'

type Listener = (...args: any[]) => void

const HEARTBEAT_INTERVAL = 5 * 60 * 1000 // 5 minutes
const RECONNECT_DELAY = 3000

export interface IncomingMessage {
  command: MessageCommand
  content: Uint8Array
  seq: number
  createdAt: number
  roomId: bigint
}

class WSClient {
  private socket: WebSocket | null = null
  private heartbeatTimer: number | null = null
  private reconnectTimer: number | null = null
  private listeners: Record<string, Set<Listener>> = {}
  private manuallyClosed = false

  on(ev: WSEventName, fn: Listener) {
    ;(this.listeners[ev] ||= new Set()).add(fn)
  }
  off(ev: WSEventName, fn: Listener) {
    this.listeners[ev]?.delete(fn)
  }
  private emit(ev: WSEventName, ...args: any[]) {
    this.listeners[ev]?.forEach((fn) => {
      try {
        fn(...args)
      } catch (e) {
        console.error('ws listener error', e)
      }
    })
  }

  isOpen(): boolean {
    return this.socket?.readyState === WebSocket.OPEN
  }

  connect() {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      return
    }
    if (this.socket) {
      try { this.socket.close() } catch {}
      this.socket = null
    }
    const { userId, deviceId, token } = auth
    if (!userId || !deviceId || !token) return
    const url = `ws://127.0.0.1:8003/ws?user_id=${encodeURIComponent(userId)}&device_id=${encodeURIComponent(deviceId)}&token=${encodeURIComponent(token)}`
    this.manuallyClosed = false
    const sock = new WebSocket(url)
    sock.binaryType = 'arraybuffer'
    this.socket = sock
    sock.onopen = () => {
      this.scheduleHeartbeat()
      this.emit('open')
    }
    sock.onclose = () => {
      this.clearHeartbeat()
      this.emit('close')
      if (!this.manuallyClosed) {
        this.scheduleReconnect()
      }
    }
    sock.onerror = () => {
      // browsers don't expose detail; close handler will fire next.
    }
    sock.onmessage = (ev) => this.handleMessage(ev.data as ArrayBuffer)
  }

  disconnect() {
    this.manuallyClosed = true
    this.clearHeartbeat()
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer)
      this.reconnectTimer = null
    }
    if (this.socket) {
      try {
        this.socket.close()
      } catch {}
      this.socket = null
    }
  }

  private scheduleReconnect() {
    if (this.reconnectTimer) return
    this.reconnectTimer = window.setTimeout(() => {
      this.reconnectTimer = null
      this.connect()
    }, RECONNECT_DELAY)
  }

  private scheduleHeartbeat() {
    this.clearHeartbeat()
    this.heartbeatTimer = window.setInterval(() => {
      this.sendPacket(PacketCommand.PC_HEARTBEAT, new Uint8Array(0))
    }, HEARTBEAT_INTERVAL)
  }
  private clearHeartbeat() {
    if (this.heartbeatTimer) {
      clearInterval(this.heartbeatTimer)
      this.heartbeatTimer = null
    }
  }

  private sendPacket(command: PacketCommand, content: Uint8Array) {
    if (!this.socket || this.socket.readyState !== WebSocket.OPEN) return
    const packet = create(PacketSchema, {
      command,
      content,
      requestId: '',
      seq: BigInt(0),
    })
    const buf = toBinary(PacketSchema, packet)
    this.socket.send(buf)
  }

  private handleMessage(buf: ArrayBuffer) {
    const u8 = new Uint8Array(buf)
    const packet = fromBinary(PacketSchema, u8)
    if (packet.command === PacketCommand.PC_HEARTBEAT) {
      // server echoes; nothing to do
      return
    }
    if (packet.command === PacketCommand.PC_MESSAGE) {
      const msg = fromBinary(MessageSchema, packet.content)
      const incoming: IncomingMessage = {
        command: msg.command,
        content: msg.content,
        seq: Number(msg.seq),
        createdAt: Number(msg.createdAt),
        roomId: msg.roomId,
      }
      this.emit('message', incoming)
    }
  }
}

export const ws = new WSClient()
