<script setup lang="ts">
import { computed, nextTick, onMounted, onUpdated, ref, watch } from 'vue'
import { NButton, NDrawer, NDrawerContent, NInput, NModal, NPopover, useMessage } from 'naive-ui'
import { ChevronRight, MoreHorizontal, Phone, Settings, Video } from 'lucide-vue-next'
import Avatar from '@/components/Avatar.vue'
import MessageEditor from './MessageEditor.vue'
import UserCard from '@/components/UserCard.vue'
import { state, recordOutgoingMessage, setFriendRemarks } from '@/api/messageSync'
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

const chatBodyEl = ref<HTMLElement | null>(null)
const showDrawer = ref(false)
const showSettings = ref(false)
const remarksInput = ref('')
const savingRemarks = ref(false)

watch(showSettings, (v) => {
  if (v) remarksInput.value = conv.value?.remarks ?? ''
})
watch(() => props.convKey, () => {
  showDrawer.value = false
  showSettings.value = false
})

function openSettings() {
  showDrawer.value = false
  showSettings.value = true
}

async function onSaveAndClose() {
  if (!conv.value || conv.value.type !== 'user') return
  if (savingRemarks.value) return
  savingRemarks.value = true
  try {
    await setFriendRemarks(conv.value.peerId, remarksInput.value.trim())
    message.success('已保存')
    showSettings.value = false
  } catch (e: any) {
    message.error('保存失败：' + (e?.message ?? e))
  } finally {
    savingRemarks.value = false
  }
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
      <div class="title">{{ conv.remarks || conv.name }}</div>
      <div class="actions">
        <button class="hbtn" disabled><Phone :size="18" /></button>
        <button class="hbtn" disabled><Video :size="18" /></button>
        <button
          class="hbtn"
          :class="{ active: showDrawer }"
          title="更多"
          @click="showDrawer = !showDrawer"
        ><MoreHorizontal :size="18" /></button>
      </div>
    </div>
    <div class="chat-body" ref="chatBodyEl">
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
                <UserCard :user-id="it.msg!.fromUserId" @open-chat="() => $emit('open-user', it.msg!.fromUserId)" />
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

      <n-drawer
        v-model:show="showDrawer"
        :width="240"
        placement="right"
        :to="chatBodyEl ?? undefined"
        :show-mask="false"
      >
        <n-drawer-content >
          <div class="drawer-stack" v-if="conv">
            <div class="info-card">
              <Avatar
                :src="conv.avatarUrl"
                :name="conv.remarks || conv.name || conv.peerId"
                :size="56"
                :rounded="10"
              />
              <div class="info-text">
                <div class="info-name">{{ conv.remarks || conv.name || (conv.type === 'group' ? '群组' : '用户') + conv.peerId }}</div>
                <div v-if="conv.type === 'user'" class="info-sub">
                  昵称：{{ conv.name || ('用户' + conv.peerId) }}
                </div>
                <div v-if="conv.type === 'user'" class="info-sub">
                  备注：<span :class="{ 'info-muted': !conv.remarks }">{{ conv.remarks || '未设置' }}</span>
                </div>
                <div class="info-sub">
                  {{ conv.type === 'group' ? '群ID' : 'UID' }}：{{ conv.peerId }}
                </div>
              </div>
            </div>

            <div class="menu-list">
              <div
                v-if="conv.type === 'user'"
                class="menu-item"
                role="button"
                tabindex="0"
                @click="openSettings"
                @keyup.enter="openSettings"
              >
                <Settings :size="16" />
                <span class="menu-label">设置</span>
                <ChevronRight :size="16" class="menu-arrow" />
              </div>
              <div v-else class="menu-empty">群聊设置暂未开放</div>
            </div>
          </div>
        </n-drawer-content>
      </n-drawer>
    </div>

    <n-modal
      v-model:show="showSettings"
      preset="card"
      title="聊天设置"
      :style="{ width: '420px' }"
      :bordered="false"
      :mask-closable="!savingRemarks"
      :close-on-esc="!savingRemarks"
    >
      <div class="dialog-content" v-if="conv">
        <div class="dialog-peer">
          <Avatar
            :src="conv.avatarUrl"
            :name="conv.remarks || conv.name || conv.peerId"
            :size="44"
            :rounded="8"
          />
          <div class="dialog-peer-text">
            <div class="dialog-peer-name">{{ conv.remarks || conv.name || ('用户' + conv.peerId) }}</div>
            <div class="dialog-peer-uid">UID: {{ conv.peerId }}</div>
          </div>
        </div>
        <div class="dialog-form">
          <label class="dialog-label" for="remarks-input">好友备注</label>
          <n-input
            id="remarks-input"
            v-model:value="remarksInput"
            placeholder="为该好友设置备注"
            :maxlength="32"
            show-count
            clearable
            @keyup.enter="onSaveAndClose"
          />
          <div class="dialog-hint">备注仅自己可见，不会通知对方</div>
        </div>
      </div>
      <template #footer>
        <div class="dialog-footer">
          <n-button
            quaternary
            :disabled="savingRemarks"
            @click="showSettings = false"
          >取消</n-button>
          <n-button
            type="primary"
            :loading="savingRemarks"
            :disabled="savingRemarks"
            @click="onSaveAndClose"
          >保存</n-button>
        </div>
      </template>
    </n-modal>
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
.hbtn.active { background: rgba(127,127,127,0.18); opacity: 1; }
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
.chat-body {
  flex: 1;
  min-height: 0;
  position: relative;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}
.drawer-stack {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin: -8px -16px 0;
}
.info-card {
  display: flex;
  gap: 12px;
  align-items: center;
  padding: 16px;
  background: rgba(127, 127, 127, 0.08);
  border-bottom: 1px solid rgba(127, 127, 127, 0.12);
}
.info-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
  flex: 1;
}
.info-name {
  font-size: 15px;
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.info-sub {
  font-size: 12px;
  opacity: 0.7;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.info-muted {
  opacity: 0.55;
}
.menu-list {
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: 0 8px;
}
.menu-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  cursor: pointer;
  border-radius: 6px;
  font-size: 13px;
}
.menu-item:hover {
  background: rgba(127, 127, 127, 0.10);
}
.menu-label {
  flex: 1;
}
.menu-arrow {
  opacity: 0.45;
}
.menu-empty {
  text-align: center;
  font-size: 12px;
  opacity: 0.55;
  padding: 24px 8px;
}
.dialog-content {
  display: flex;
  flex-direction: column;
  gap: 18px;
}
.dialog-peer {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 14px;
  background: rgba(127, 127, 127, 0.08);
  border-radius: 8px;
}
.dialog-peer-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
  flex: 1;
}
.dialog-peer-name {
  font-size: 14px;
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.dialog-peer-uid {
  font-size: 11px;
  opacity: 0.6;
}
.dialog-form {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.dialog-label {
  font-size: 13px;
  font-weight: 500;
}
.dialog-hint {
  font-size: 11px;
  opacity: 0.55;
  margin-top: 2px;
}
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}
</style>
