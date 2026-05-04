<script setup lang="ts">
import { computed, nextTick, onMounted, onUpdated, ref, watch } from 'vue'
import { NPopover, NSpin, useMessage } from 'naive-ui'
import { Phone, Video, MoreHorizontal } from 'lucide-vue-next'
import Avatar from '@/components/Avatar.vue'
import MessageEditor from './MessageEditor.vue'
import UserCard from '@/components/UserCard.vue'
import { state, recordOutgoingMessage } from '@/api/messageSync'
import { messageClient } from '@/api/transport'
import { renderMarkdown } from '@/utils/markdown'
import { formatTimeLabel } from '@/utils/time'
import { auth } from '@/store/auth'

const props = defineProps<{ convKey: string }>()
const emit = defineEmits<{ (e: 'open-user', userId: string): void }>()
const message = useMessage()
const scrollerRef = ref<HTMLDivElement>()

const conv = computed(() => state.conversations.find((c) => c.key === props.convKey))
const messages = computed(() => state.messages[props.convKey] || [])

interface MessageVM {
  type: 'message' | 'time'
  id: string
  // for time:
  label?: string
  // for message:
  msg?: typeof messages.value[number]
  showAvatar?: boolean
  showName?: boolean
}

const items = computed<MessageVM[]>(() => {
  const out: MessageVM[] = []
  let lastTs = 0
  let lastSender = ''
  const list = messages.value
  for (const m of list) {
    if (!lastTs || m.createdAt - lastTs > 5 * 60 * 1000) {
      out.push({ type: 'time', id: `time-${m.id}`, label: formatTimeLabel(m.createdAt) })
      lastSender = ''
    }
    const showName = !m.self && m.isGroup && lastSender !== m.fromUserId
    out.push({ type: 'message', id: m.id, msg: m, showAvatar: true, showName })
    lastTs = m.createdAt
    lastSender = m.fromUserId
  }
  return out
})

function scrollToBottom() {
  nextTick(() => {
    const el = scrollerRef.value
    if (el) el.scrollTop = el.scrollHeight
  })
}

watch(() => props.convKey, () => scrollToBottom())
watch(() => messages.value.length, () => scrollToBottom())
onMounted(() => scrollToBottom())
onUpdated(() => {
  // keep scrolled to bottom unless user has scrolled up
})

const enc = new TextEncoder()
const sending = ref(false)

async function onSend(text: string) {
  if (!conv.value) return
  if (sending.value) return
  sending.value = true
  try {
    const bytes = enc.encode(text)
    if (conv.value.type === 'user') {
      const r = await messageClient.sendFriendMessage({ userId: BigInt(conv.value.peerId), content: bytes })
      const seq = Number(r.fromUserSeq)
      await recordOutgoingMessage({ isGroup: false, peerId: conv.value.peerId, seq, text })
    } else {
      const r = await messageClient.sendGroupMessage({ groupId: BigInt(conv.value.peerId), content: bytes })
      const seq = Number(r.fromUserSeq)
      await recordOutgoingMessage({ isGroup: true, peerId: conv.value.peerId, seq, text })
    }
    scrollToBottom()
  } catch (e: any) {
    console.error(e)
    message.error('发送失败：' + (e?.message ?? e))
  } finally {
    sending.value = false
  }
}

function openUserCard(userId: string) {
  emit('open-user', userId)
}

function nineGridUrls(_groupKey: string): string[] {
  // Lookup last few unique sender avatars in this group as a stand-in 9-grid.
  const seen = new Map<string, string>()
  for (const m of messages.value) {
    if (m.fromAvatarUrl && !seen.has(m.fromUserId)) {
      seen.set(m.fromUserId, m.fromAvatarUrl)
      if (seen.size >= 9) break
    }
  }
  return Array.from(seen.values())
}
</script>

<template>
  <div class="chat-root">
    <div class="header" v-if="conv">
      <div class="title">{{ conv.name }}</div>
      <div class="actions">
        <button class="hbtn" disabled><Phone :size="18" /></button>
        <button class="hbtn" disabled><Video :size="18" /></button>
        <button class="hbtn" disabled><MoreHorizontal :size="18" /></button>
      </div>
    </div>
    <div class="scroller" ref="scrollerRef">
      <div v-if="!items.length" class="hint-empty">还没有消息，开始聊天吧</div>
      <template v-for="it in items" :key="it.id">
        <div v-if="it.type === 'time'" class="time-sep">{{ it.label }}</div>
        <div v-else class="msg-row" :class="{ self: it.msg!.self }">
          <div class="msg-avatar" v-if="!it.msg!.self">
            <n-popover trigger="click" placement="top-start">
              <template #trigger>
                <span><Avatar :src="it.msg!.fromAvatarUrl" :name="it.msg!.fromNickname" :size="36" :rounded="6" /></span>
              </template>
              <UserCard :user-id="it.msg!.fromUserId" @open-chat="(k) => $emit('open-user', it.msg!.fromUserId)" />
            </n-popover>
          </div>
          <div class="msg-content">
            <div class="msg-name" v-if="it.showName">{{ it.msg!.fromNickname }}</div>
            <div class="msg-bubble" :class="{ self: it.msg!.self }">
              <div class="md-body user-text" v-html="renderMarkdown(it.msg!.content)"></div>
            </div>
          </div>
          <div class="msg-avatar" v-if="it.msg!.self">
            <Avatar :src="auth.avatarUrl" :name="auth.nickname" :size="36" :rounded="6" />
          </div>
        </div>
      </template>
    </div>
    <MessageEditor @send="onSend" />
  </div>
</template>

<style scoped>
.chat-root {
  display: flex;
  flex-direction: column;
  height: 100%;
}
.header {
  height: 52px;
  border-bottom: 1px solid rgba(127, 127, 127, 0.16);
  padding: 0 18px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.title {
  font-size: 15px;
  font-weight: 600;
}
.actions {
  display: flex;
  gap: 4px;
}
.hbtn {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  border: none;
  background: transparent;
  color: inherit;
  cursor: pointer;
  opacity: 0.7;
  display: flex;
  align-items: center;
  justify-content: center;
}
.hbtn:hover { background: rgba(127,127,127,0.18); opacity: 1; }
.hbtn:disabled { opacity: 0.3; cursor: not-allowed; }
.scroller {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding: 16px 18px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.hint-empty {
  margin: auto;
  color: rgba(127,127,127,0.6);
  font-size: 13px;
}
.time-sep {
  align-self: center;
  color: rgba(127,127,127,0.55);
  font-size: 11px;
  background: rgba(127,127,127,0.12);
  padding: 1px 8px;
  border-radius: 8px;
  margin: 8px 0 4px;
}
.msg-row {
  display: flex;
  gap: 8px;
  align-items: flex-start;
  max-width: 100%;
}
.msg-row.self {
  justify-content: flex-end;
}
.msg-content {
  max-width: min(560px, 70%);
  display: flex;
  flex-direction: column;
}
.msg-row.self .msg-content {
  align-items: flex-end;
}
.msg-name {
  font-size: 11px;
  color: rgba(127,127,127,0.85);
  margin: 2px 4px;
}
.msg-bubble {
  background: var(--n-action-color, rgba(127,127,127,0.18));
  color: inherit;
  padding: 8px 12px;
  border-radius: 8px;
  word-break: break-word;
}
.msg-bubble.self {
  background: var(--n-primary-color, #07c160);
  color: #fff;
}
.msg-bubble.self :deep(a) { color: #fff; text-decoration: underline; }
.msg-bubble.self :deep(code) { background: rgba(255,255,255,0.18); }
</style>
