import { fromBinary } from '@bufbuild/protobuf'
import { reactive, ref } from 'vue'
import { messageLogicClient, friendClient, userClient, groupClient } from './transport'
import { ws, type IncomingMessage } from './ws'
import { MessageCommand } from '@/gen/proto/connect/connect.ext_pb'
import { UserMessagePushSchema, GroupMessagePushSchema } from '@/gen/proto/logic/message.ext_pb'
import { AddFriendPushSchema, AgreeAddFriendPushSchema, UpdateGroupPushSchema } from '@/gen/proto/business/friend.ext_pb'
import {
  convKey,
  type ConvRecord,
  type FriendRequestRecord,
  type MessageRecord,
  findMessageBySeq,
  getConversation,
  getFriendRequests,
  getMaxSeq,
  listConversations,
  listMessages,
  putMessage,
  setFriendRequests,
  setMaxSeq,
  upsertConversation,
  clearUnread as dbClearUnread,
  deleteConversation as dbDeleteConversation,
} from '@/db'
import { auth } from '@/store/auth'
import { summarizeMarkdown } from '@/utils/markdown'

interface State {
  conversations: ConvRecord[]
  activeKey: string | null
  messages: Record<string, MessageRecord[]>
  friendRequests: FriendRequestRecord[]
  totalUnread: number
  pendingFriendRequests: number
  // Profiles cache: user_id -> {nickname, avatar}.
  profiles: Record<string, { nickname: string; avatarUrl: string }>
  // Group meta cache: group_id -> name, avatar
  groups: Record<string, { name: string; avatarUrl: string }>
}

export const state = reactive<State>({
  conversations: [],
  activeKey: null,
  messages: {},
  friendRequests: [],
  totalUnread: 0,
  pendingFriendRequests: 0,
  profiles: {},
  groups: {},
})

export const wsConnected = ref(false)

function decoder() {
  return new TextDecoder('utf-8')
}

function normalizeTs(v: number): number {
  // Server may emit seconds; client expects ms.
  if (!v) return Date.now()
  return v < 1e12 ? v * 1000 : v
}

function recomputeBadges() {
  state.totalUnread = state.conversations.reduce((acc, c) => acc + (c.unread || 0), 0)
  state.pendingFriendRequests = state.friendRequests.filter((r) => r.status === 'pending').length
}

export async function refreshConversations() {
  state.conversations = await listConversations()
  recomputeBadges()
}

export async function refreshMessages(key: string) {
  state.messages[key] = await listMessages(key)
}

export async function setActive(key: string | null) {
  state.activeKey = key
  if (key) {
    await dbClearUnread(key)
    await refreshConversations()
    await refreshMessages(key)
  }
}

async function ensureUser(userId: string): Promise<{ nickname: string; avatarUrl: string }> {
  if (state.profiles[userId]) return state.profiles[userId]
  try {
    const r = await userClient.getUser({ userId: BigInt(userId) })
    if (r.user) {
      state.profiles[userId] = { nickname: r.user.nickname, avatarUrl: r.user.avatarUrl }
      return state.profiles[userId]
    }
  } catch (e) {
    console.warn('getUser failed', userId, e)
  }
  return { nickname: `用户${userId}`, avatarUrl: '' }
}

async function ensureGroup(groupId: string): Promise<{ name: string; avatarUrl: string }> {
  if (state.groups[groupId]) return state.groups[groupId]
  try {
    const r = await groupClient.get({ groupId: BigInt(groupId) })
    if (r.group) {
      state.groups[groupId] = { name: r.group.name, avatarUrl: r.group.avatarUrl }
      return state.groups[groupId]
    }
  } catch (e) {
    console.warn('group.get failed', groupId, e)
  }
  return { name: `群组${groupId}`, avatarUrl: '' }
}

async function applyMessage(seq: number, command: MessageCommand, content: Uint8Array, createdAt: number) {
  if (command === MessageCommand.MC_USER_MESSAGE) {
    const push = fromBinary(UserMessagePushSchema, content)
    const fromId = push.fromUser?.userId.toString() ?? '0'
    const toId = push.toUserId.toString()
    const isSelf = fromId === auth.userId
    const peerId = isSelf ? toId : fromId
    const key = convKey('user', peerId)
    const text = decoder().decode(push.content)
    if (push.fromUser) {
      state.profiles[fromId] = { nickname: push.fromUser.nickname, avatarUrl: push.fromUser.avatarUrl }
    }
    let peer = state.profiles[peerId]
    if (!peer) peer = await ensureUser(peerId)
    const rec: MessageRecord = {
      id: `${key}:${seq}`,
      convKey: key,
      seq,
      fromUserId: fromId,
      fromNickname: push.fromUser?.nickname || state.profiles[fromId]?.nickname || `用户${fromId}`,
      fromAvatarUrl: push.fromUser?.avatarUrl || state.profiles[fromId]?.avatarUrl || '',
      toId: peerId,
      isGroup: false,
      content: text,
      createdAt,
      self: isSelf,
    }
    await putMessage(rec)
    const conv: ConvRecord = (await getConversation(key)) || {
      key,
      type: 'user',
      peerId,
      name: peer.nickname,
      avatarUrl: peer.avatarUrl,
      lastSeq: 0,
      lastMessageId: '',
      lastPreview: '',
      lastTimestamp: 0,
      unread: 0,
    }
    conv.name = peer.nickname || conv.name
    conv.avatarUrl = peer.avatarUrl || conv.avatarUrl
    if (seq >= conv.lastSeq) {
      conv.lastSeq = seq
      conv.lastMessageId = rec.id
      conv.lastPreview = summarizeMarkdown(text)
      conv.lastTimestamp = createdAt
    }
    if (!isSelf && state.activeKey !== key) {
      conv.unread = (conv.unread || 0) + 1
    }
    await upsertConversation(conv)
  } else if (command === MessageCommand.MC_GROUP_MESSAGE) {
    const push = fromBinary(GroupMessagePushSchema, content)
    const fromId = push.fromUser?.userId.toString() ?? '0'
    const groupId = push.groupId.toString()
    const isSelf = fromId === auth.userId
    const key = convKey('group', groupId)
    const text = decoder().decode(push.content)
    if (push.fromUser) {
      state.profiles[fromId] = { nickname: push.fromUser.nickname, avatarUrl: push.fromUser.avatarUrl }
    }
    const g = await ensureGroup(groupId)
    const rec: MessageRecord = {
      id: `${key}:${seq}`,
      convKey: key,
      seq,
      fromUserId: fromId,
      fromNickname: push.fromUser?.nickname || state.profiles[fromId]?.nickname || `用户${fromId}`,
      fromAvatarUrl: push.fromUser?.avatarUrl || state.profiles[fromId]?.avatarUrl || '',
      toId: groupId,
      isGroup: true,
      content: text,
      createdAt,
      self: isSelf,
    }
    await putMessage(rec)
    const conv: ConvRecord = (await getConversation(key)) || {
      key,
      type: 'group',
      peerId: groupId,
      name: g.name,
      avatarUrl: g.avatarUrl,
      lastSeq: 0,
      lastMessageId: '',
      lastPreview: '',
      lastTimestamp: 0,
      unread: 0,
    }
    conv.name = g.name || conv.name
    conv.avatarUrl = g.avatarUrl || conv.avatarUrl
    if (seq >= conv.lastSeq) {
      conv.lastSeq = seq
      conv.lastMessageId = rec.id
      conv.lastPreview = summarizeMarkdown(text)
      conv.lastTimestamp = createdAt
    }
    if (!isSelf && state.activeKey !== key) {
      conv.unread = (conv.unread || 0) + 1
    }
    await upsertConversation(conv)
  }
}

async function applyControlPush(command: MessageCommand, content: Uint8Array) {
  console.log('[ctrl]', { command, contentLen: content.length })
  // Currently not used since pushes are sent as MC_USER_MESSAGE / MC_GROUP_MESSAGE only.
  // The command numbers 110/111/120 (ADD_FRIEND etc.) are package-level Command codes but
  // they arrive over the WebSocket inside connect.Message wrappers in some setups; we read
  // the wrapper command (MessageCommand). We additionally handle the explicit codes here in
  // case the server uses them as MessageCommand values.
  const code = command as unknown as number
  if (code === 110) {
    const push = fromBinary(AddFriendPushSchema, content)
    await pushFriendRequest({
      friendId: push.friendId.toString(),
      nickname: push.nickname,
      avatarUrl: push.avatarUrl,
      description: push.description,
      receivedAt: Date.now(),
      status: 'pending',
    })
  } else if (code === 111) {
    const push = fromBinary(AgreeAddFriendPushSchema, content)
    state.profiles[push.friendId.toString()] = { nickname: push.nickname, avatarUrl: push.avatarUrl }
  } else if (code === 120) {
    const push = fromBinary(UpdateGroupPushSchema, content)
    void push
    // A future iteration could refresh group meta here.
  }
}

export async function pushFriendRequest(req: FriendRequestRecord) {
  const dup = state.friendRequests.find((r) => r.friendId === req.friendId && r.status === 'pending')
  if (dup) return
  state.friendRequests.unshift(req)
  await setFriendRequests([...state.friendRequests])
  recomputeBadges()
}

export async function loadFriendRequests() {
  state.friendRequests = await getFriendRequests()
  recomputeBadges()
}

export async function updateFriendRequest(friendId: string, status: 'agreed' | 'ignored') {
  state.friendRequests = state.friendRequests.map((r) => (r.friendId === friendId ? { ...r, status } : r))
  await setFriendRequests([...state.friendRequests])
  recomputeBadges()
}

async function syncOffline() {
  let from = await getMaxSeq()
  for (let i = 0; i < 20; i++) {
    const r = await messageLogicClient.sync({ seq: BigInt(from) })
    if (!r.messages || r.messages.length === 0) break
    let max = from
    for (const m of r.messages) {
      const seq = Number(m.seq)
      const createdAt = normalizeTs(Number(m.createdAt))
      const cmd = m.command as MessageCommand
      const code = cmd as unknown as number
      console.log('[sync]', { seq, code, contentLen: m.content.length })
      try {
        if (code === MessageCommand.MC_USER_MESSAGE || code === MessageCommand.MC_GROUP_MESSAGE) {
          await applyMessage(seq, cmd, m.content, createdAt)
        } else {
          await applyControlPush(cmd, m.content)
        }
      } catch (e) {
        console.warn('apply offline message failed', e, { code, seq })
      }
      if (seq > max) max = seq
    }
    if (max > from) {
      await setMaxSeq(max)
      from = max
      // ack server
      try {
        await messageLogicClient.messageACK({ deviceAck: BigInt(max) })
      } catch (e) {
        console.warn('ack failed', e)
      }
    }
    if (!r.hasMore) break
  }
  await refreshConversations()
  if (state.activeKey) await refreshMessages(state.activeKey)
}

let initialized = false
let listenersBound = false

export async function bootstrap() {
  if (initialized) return
  initialized = true

  await loadFriendRequests()
  await refreshConversations()
  void syncFriendRemarks()

  if (listenersBound) {
    ws.connect()
    return
  }
  listenersBound = true

  ws.on('open', async () => {
    wsConnected.value = true
    try {
      await syncOffline()
    } catch (e) {
      console.warn('sync failed', e)
    }
  })
  ws.on('close', () => {
    wsConnected.value = false
  })
  function normalizeIncoming(i: IncomingMessage): IncomingMessage {
    return { ...i, createdAt: normalizeTs(i.createdAt) }
  }
  ws.on('message', async (incoming: IncomingMessage) => {
    incoming = normalizeIncoming(incoming)
    try {
      const code = incoming.command as unknown as number
      if (code === MessageCommand.MC_USER_MESSAGE || code === MessageCommand.MC_GROUP_MESSAGE) {
        await applyMessage(incoming.seq, incoming.command, incoming.content, incoming.createdAt)
      } else {
        await applyControlPush(incoming.command, incoming.content)
      }
      const cur = await getMaxSeq()
      if (incoming.seq > cur) {
        await setMaxSeq(incoming.seq)
        try {
          await messageLogicClient.messageACK({ deviceAck: BigInt(incoming.seq) })
        } catch {}
      }
      await refreshConversations()
      if (state.activeKey) await refreshMessages(state.activeKey)
    } catch (e) {
      console.warn('on message error', e)
    }
  })

  ws.connect()
}

export function shutdown() {
  ws.disconnect()
  initialized = false
  state.conversations = []
  state.messages = {}
  state.friendRequests = []
  state.activeKey = null
  state.profiles = {}
  state.groups = {}
}

export async function recordOutgoingMessage(opts: {
  isGroup: boolean
  peerId: string
  seq: number
  text: string
}) {
  const key = convKey(opts.isGroup ? 'group' : 'user', opts.peerId)
  const rec: MessageRecord = {
    id: `${key}:${opts.seq}`,
    convKey: key,
    seq: opts.seq,
    fromUserId: auth.userId,
    fromNickname: auth.nickname,
    fromAvatarUrl: auth.avatarUrl,
    toId: opts.peerId,
    isGroup: opts.isGroup,
    content: opts.text,
    createdAt: Date.now(),
    self: true,
  }
  // De-dupe vs. eventual sync.
  const existing = await findMessageBySeq(key, opts.seq)
  if (!existing) await putMessage(rec)
  let conv = await getConversation(key)
  if (!conv) {
    if (opts.isGroup) {
      const g = await ensureGroup(opts.peerId)
      conv = {
        key,
        type: 'group',
        peerId: opts.peerId,
        name: g.name,
        avatarUrl: g.avatarUrl,
        lastSeq: opts.seq,
        lastMessageId: rec.id,
        lastPreview: summarizeMarkdown(opts.text),
        lastTimestamp: rec.createdAt,
        unread: 0,
      }
    } else {
      const p = await ensureUser(opts.peerId)
      conv = {
        key,
        type: 'user',
        peerId: opts.peerId,
        name: p.nickname,
        avatarUrl: p.avatarUrl,
        lastSeq: opts.seq,
        lastMessageId: rec.id,
        lastPreview: summarizeMarkdown(opts.text),
        lastTimestamp: rec.createdAt,
        unread: 0,
      }
    }
  } else {
    if (opts.seq >= conv.lastSeq) {
      conv.lastSeq = opts.seq
      conv.lastMessageId = rec.id
      conv.lastPreview = summarizeMarkdown(opts.text)
      conv.lastTimestamp = rec.createdAt
    }
  }
  await upsertConversation(conv)
  const cur = await getMaxSeq()
  if (opts.seq > cur) await setMaxSeq(opts.seq)
  await refreshConversations()
  if (state.activeKey === key) await refreshMessages(key)
}

export async function setFriendRemarks(peerId: string, remarks: string) {
  await friendClient.set({ friendId: BigInt(peerId), remarks })
  const key = convKey('user', peerId)
  let conv = await getConversation(key)
  if (!conv) {
    const p = await ensureUser(peerId)
    conv = {
      key,
      type: 'user',
      peerId,
      name: p.nickname,
      avatarUrl: p.avatarUrl,
      lastSeq: 0,
      lastMessageId: '',
      lastPreview: '',
      lastTimestamp: Date.now(),
      unread: 0,
      remarks,
    }
  } else {
    conv.remarks = remarks
  }
  await upsertConversation(conv)
  await refreshConversations()
}

async function syncFriendRemarks() {
  try {
    const r = await friendClient.getFriends({})
    for (const f of r.friends || []) {
      const peerId = f.userId.toString()
      const key = convKey('user', peerId)
      const conv = await getConversation(key)
      if (!conv) continue
      if (conv.remarks !== f.remarks) {
        conv.remarks = f.remarks
        await upsertConversation(conv)
      }
    }
    await refreshConversations()
  } catch (e) {
    console.warn('syncFriendRemarks failed', e)
  }
}

export async function ensurePeerConversation(peerId: string, peerNickname: string, peerAvatar: string) {
  const key = convKey('user', peerId)
  const cur = await getConversation(key)
  if (!cur) {
    await upsertConversation({
      key,
      type: 'user',
      peerId,
      name: peerNickname,
      avatarUrl: peerAvatar,
      lastSeq: 0,
      lastMessageId: '',
      lastPreview: '',
      lastTimestamp: Date.now(),
      unread: 0,
    })
    await refreshConversations()
  }
  state.profiles[peerId] = { nickname: peerNickname, avatarUrl: peerAvatar }
  return key
}

export async function ensureGroupConversation(groupId: string, name: string, avatarUrl: string) {
  const key = convKey('group', groupId)
  const cur = await getConversation(key)
  if (!cur) {
    await upsertConversation({
      key,
      type: 'group',
      peerId: groupId,
      name,
      avatarUrl,
      lastSeq: 0,
      lastMessageId: '',
      lastPreview: '',
      lastTimestamp: Date.now(),
      unread: 0,
    })
    await refreshConversations()
  }
  state.groups[groupId] = { name, avatarUrl }
  return key
}

export async function deleteConversation(key: string) {
  await dbDeleteConversation(key)
  if (state.activeKey === key) state.activeKey = null
  await refreshConversations()
  delete state.messages[key]
}
