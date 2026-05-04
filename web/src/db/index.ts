// IndexedDB wrapper. Schema is per-user — we open a different DB name per logged-in account
// so signing out and back in as a different user does not leak history.

import { auth } from '@/store/auth'

export interface ConvKey {
  type: 'user' | 'group'
  id: string // peer user_id or group_id
}

export interface ConvRecord {
  key: string // `${type}:${id}`
  type: 'user' | 'group'
  peerId: string
  name: string
  avatarUrl: string
  lastSeq: number
  lastMessageId: string
  lastPreview: string
  lastTimestamp: number
  unread: number
}

export interface MessageRecord {
  id: string // `${convKey}:${seq}` or `${convKey}:local:${nonce}` for outbound-not-yet-acked
  convKey: string
  seq: number
  fromUserId: string
  fromNickname: string
  fromAvatarUrl: string
  // For 1-1 chats this is the peer user id; for groups, the group id.
  toId: string
  isGroup: boolean
  content: string // markdown text
  createdAt: number
  self: boolean
  pending?: boolean
  failed?: boolean
}

let db: IDBDatabase | null = null
let dbName = ''

function open(name: string): Promise<IDBDatabase> {
  return new Promise((resolve, reject) => {
    const req = indexedDB.open(name, 1)
    req.onupgradeneeded = () => {
      const d = req.result
      if (!d.objectStoreNames.contains('conversations')) {
        d.createObjectStore('conversations', { keyPath: 'key' })
      }
      if (!d.objectStoreNames.contains('messages')) {
        const s = d.createObjectStore('messages', { keyPath: 'id' })
        s.createIndex('byConv', 'convKey')
        s.createIndex('byConvSeq', ['convKey', 'seq'])
      }
      if (!d.objectStoreNames.contains('meta')) {
        d.createObjectStore('meta', { keyPath: 'k' })
      }
    }
    req.onsuccess = () => resolve(req.result)
    req.onerror = () => reject(req.error)
  })
}

export async function initDb(userId: string): Promise<void> {
  const name = `gim_${userId}`
  if (db && dbName === name) return
  if (db) {
    db.close()
    db = null
  }
  db = await open(name)
  dbName = name
}

// Lazily reopen the DB if the module-level handle was lost (e.g. Vite HMR
// re-evaluated this module, or a tx ran before initDb finished).
async function ensureDb(): Promise<IDBDatabase> {
  if (db) return db
  if (!auth.userId) throw new Error('db not initialized: no user logged in')
  await initDb(auth.userId)
  if (!db) throw new Error('db not initialized')
  return db
}

function unwrap<T>(v: T): T {
  if (v == null || typeof v !== 'object') return v
  return JSON.parse(JSON.stringify(v))
}

async function tx<T>(stores: string | string[], mode: IDBTransactionMode, fn: (t: IDBTransaction) => Promise<T> | T): Promise<T> {
  const handle = await ensureDb()
  return new Promise((resolve, reject) => {
    const t = handle.transaction(stores, mode)
    t.onerror = () => reject(t.error)
    t.onabort = () => reject(t.error)
    Promise.resolve(fn(t))
      .then((v) => {
        t.oncomplete = () => resolve(v)
      })
      .catch(reject)
  })
}

function reqProm<T>(req: IDBRequest<T>): Promise<T> {
  return new Promise((res, rej) => {
    req.onsuccess = () => res(req.result)
    req.onerror = () => rej(req.error)
  })
}

export function convKey(type: 'user' | 'group', id: string | bigint): string {
  return `${type}:${id.toString()}`
}

export async function getMaxSeq(): Promise<number> {
  return tx('meta', 'readonly', async (t) => {
    const r = await reqProm(t.objectStore('meta').get('maxSeq'))
    return r ? Number((r as any).v) : 0
  })
}

export async function setMaxSeq(seq: number): Promise<void> {
  await tx('meta', 'readwrite', (t) => {
    t.objectStore('meta').put({ k: 'maxSeq', v: seq })
  })
}

export async function listConversations(): Promise<ConvRecord[]> {
  return tx('conversations', 'readonly', async (t) => {
    const r = await reqProm(t.objectStore('conversations').getAll())
    const arr = (r ?? []) as ConvRecord[]
    arr.sort((a, b) => b.lastTimestamp - a.lastTimestamp)
    return arr
  })
}

export async function getConversation(key: string): Promise<ConvRecord | undefined> {
  return tx('conversations', 'readonly', async (t) => {
    const r = await reqProm(t.objectStore('conversations').get(key))
    return r as ConvRecord | undefined
  })
}

export async function upsertConversation(c: ConvRecord): Promise<void> {
  await tx('conversations', 'readwrite', (t) => {
    t.objectStore('conversations').put(unwrap(c))
  })
}

export async function clearUnread(key: string): Promise<void> {
  await tx('conversations', 'readwrite', async (t) => {
    const s = t.objectStore('conversations')
    const cur = (await reqProm(s.get(key))) as ConvRecord | undefined
    if (cur) {
      cur.unread = 0
      s.put(cur)
    }
  })
}

export async function deleteConversation(key: string): Promise<void> {
  await tx(['conversations', 'messages'], 'readwrite', async (t) => {
    t.objectStore('conversations').delete(key)
    const idx = t.objectStore('messages').index('byConv')
    const range = IDBKeyRange.only(key)
    const cursor = idx.openCursor(range)
    cursor.onsuccess = () => {
      const c = cursor.result
      if (c) {
        c.delete()
        c.continue()
      }
    }
  })
}

export async function listMessages(key: string): Promise<MessageRecord[]> {
  return tx('messages', 'readonly', async (t) => {
    const idx = t.objectStore('messages').index('byConv')
    const r = await reqProm(idx.getAll(IDBKeyRange.only(key)))
    const arr = (r ?? []) as MessageRecord[]
    arr.sort((a, b) => a.createdAt - b.createdAt || a.seq - b.seq)
    return arr
  })
}

export async function putMessage(m: MessageRecord): Promise<void> {
  await tx('messages', 'readwrite', (t) => {
    t.objectStore('messages').put(unwrap(m))
  })
}

export async function deleteMessage(id: string): Promise<void> {
  await tx('messages', 'readwrite', (t) => {
    t.objectStore('messages').delete(id)
  })
}

export async function findMessageBySeq(convKey: string, seq: number): Promise<MessageRecord | undefined> {
  return tx('messages', 'readonly', async (t) => {
    const idx = t.objectStore('messages').index('byConvSeq')
    const r = await reqProm(idx.get([convKey, seq]))
    return r as MessageRecord | undefined
  })
}

// Hidden friend requests (ignored ones) — we just persist a set of friend ids.
export async function getIgnoredFriendRequests(): Promise<string[]> {
  return tx('meta', 'readonly', async (t) => {
    const r = await reqProm(t.objectStore('meta').get('ignoredFR'))
    return r ? ((r as any).v as string[]) : []
  })
}

export async function setIgnoredFriendRequests(ids: string[]): Promise<void> {
  await tx('meta', 'readwrite', (t) => {
    t.objectStore('meta').put({ k: 'ignoredFR', v: unwrap(ids) })
  })
}

export async function getFriendRequests(): Promise<FriendRequestRecord[]> {
  return tx('meta', 'readonly', async (t) => {
    const r = await reqProm(t.objectStore('meta').get('friendRequests'))
    return r ? ((r as any).v as FriendRequestRecord[]) : []
  })
}

export async function setFriendRequests(list: FriendRequestRecord[]): Promise<void> {
  await tx('meta', 'readwrite', (t) => {
    t.objectStore('meta').put({ k: 'friendRequests', v: unwrap(list) })
  })
}

export interface FriendRequestRecord {
  friendId: string
  nickname: string
  avatarUrl: string
  description: string
  receivedAt: number
  status: 'pending' | 'agreed' | 'ignored'
}
