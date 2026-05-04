<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { NBadge, NButton, NDivider, NPopover, NTooltip, useMessage } from 'naive-ui'
import {
  MessageSquare,
  Users,
  Settings,
  LogOut,
  Sun,
  Moon,
  Phone,
  Video,
} from 'lucide-vue-next'
import Avatar from '@/components/Avatar.vue'
import ConversationList from '@/views/conversation/ConversationList.vue'
import ChatPane from '@/views/conversation/ChatPane.vue'
import FriendList from '@/views/friend/FriendList.vue'
import FriendDetail from '@/views/friend/FriendDetail.vue'
import NewFriends from '@/views/friend/NewFriends.vue'
import AddFriend from '@/views/friend/AddFriend.vue'
import CreateGroup from '@/views/friend/CreateGroup.vue'
import SettingsView from '@/views/SettingsView.vue'
import { state, setActive } from '@/api/messageSync'
import { auth } from '@/store/auth'
import { theme as themeStore, PRESET_COLORS } from '@/store/theme'

const emit = defineEmits<{ (e: 'logout'): void }>()
const message = useMessage()

type NavKey = 'chat' | 'friend' | 'settings'
const nav = ref<NavKey>('chat')

// What's selected within each tab.
const selectedFriendKey = ref<string | null>(null)
const friendSubview = ref<'detail' | 'new' | 'add' | 'create-group' | null>(null)

function selectConv(key: string) {
  setActive(key)
}

function selectFriend(userId: string) {
  selectedFriendKey.value = userId
  friendSubview.value = 'detail'
}

function openFriendSubview(v: 'new' | 'add' | 'create-group') {
  selectedFriendKey.value = null
  friendSubview.value = v
}

function toggleMode() {
  themeStore.mode = themeStore.mode === 'dark' ? 'light' : 'dark'
}

function pickPrimary(c: string) {
  themeStore.primary = c
}

function doLogout() {
  emit('logout')
  message.success('已退出登录')
}

onMounted(() => {
  // nothing yet
})

const navItems = computed(() => [
  { key: 'chat' as NavKey, icon: MessageSquare, badge: state.totalUnread, label: '消息' },
  { key: 'friend' as NavKey, icon: Users, badge: state.pendingFriendRequests, label: '通讯录' },
  { key: 'settings' as NavKey, icon: Settings, badge: 0, label: '设置' },
])
</script>

<template>
  <div class="layout" :class="themeStore.mode">
    <!-- Left navigation rail -->
    <div class="nav-rail">
      <div class="nav-top">
        <Avatar :src="auth.avatarUrl" :name="auth.nickname" :size="40" :rounded="6" />
      </div>
      <div class="nav-center">
        <n-tooltip v-for="item in navItems" :key="item.key" placement="right">
          <template #trigger>
            <button
              class="nav-btn"
              :class="{ active: nav === item.key }"
              @click="nav = item.key"
            >
              <n-badge :value="item.badge" :max="99" :show="item.badge > 0">
                <component :is="item.icon" :size="22" />
              </n-badge>
            </button>
          </template>
          {{ item.label }}
        </n-tooltip>
      </div>
      <div class="nav-bottom">
        <n-tooltip placement="right">
          <template #trigger>
            <button class="nav-btn small" @click="toggleMode">
              <component :is="themeStore.mode === 'dark' ? Sun : Moon" :size="18" />
            </button>
          </template>
          {{ themeStore.mode === 'dark' ? '切换浅色' : '切换深色' }}
        </n-tooltip>

        <n-popover trigger="click" placement="right" :show-arrow="false">
          <template #trigger>
            <button class="nav-btn small" :title="'主题色'">
              <span class="dot" :style="{ background: themeStore.primary }" />
            </button>
          </template>
          <div class="color-grid">
            <button
              v-for="c in PRESET_COLORS"
              :key="c.value"
              class="color-cell"
              :style="{ background: c.value }"
              :class="{ chosen: themeStore.primary === c.value }"
              :title="c.label"
              @click="pickPrimary(c.value)"
            />
          </div>
        </n-popover>

        <n-tooltip placement="right">
          <template #trigger>
            <button class="nav-btn small" @click="doLogout">
              <LogOut :size="18" />
            </button>
          </template>
          退出登录
        </n-tooltip>
      </div>
    </div>

    <!-- Middle list -->
    <div class="middle-pane">
      <template v-if="nav === 'chat'">
        <ConversationList :active-key="state.activeKey" @select="selectConv" />
      </template>
      <template v-else-if="nav === 'friend'">
        <FriendList
          :active="selectedFriendKey"
          @select="selectFriend"
          @open-new="openFriendSubview('new')"
          @open-add="openFriendSubview('add')"
          @open-create-group="openFriendSubview('create-group')"
        />
      </template>
      <template v-else>
        <div class="middle-stub">
          <div class="middle-stub-title">设置</div>
        </div>
      </template>
    </div>

    <!-- Right content -->
    <div class="content-pane">
      <template v-if="nav === 'chat'">
        <ChatPane v-if="state.activeKey" :conv-key="state.activeKey" @open-user="selectFriend" />
        <div v-else class="empty">
          <div class="empty-inner">选择一个会话开始聊天</div>
        </div>
      </template>
      <template v-else-if="nav === 'friend'">
        <FriendDetail
          v-if="friendSubview === 'detail' && selectedFriendKey"
          :user-id="selectedFriendKey"
          @open-chat="(k) => { nav = 'chat'; setActive(k) }"
        />
        <NewFriends v-else-if="friendSubview === 'new'" />
        <AddFriend v-else-if="friendSubview === 'add'" />
        <CreateGroup
          v-else-if="friendSubview === 'create-group'"
          @created="(key) => { nav = 'chat'; setActive(key); friendSubview = null }"
        />
        <div v-else class="empty"><div class="empty-inner">选择联系人或操作</div></div>
      </template>
      <template v-else>
        <SettingsView />
      </template>
    </div>
  </div>
</template>

<style scoped>
.layout {
  display: grid;
  grid-template-columns: 64px 280px 1fr;
  height: 100%;
  width: 100%;
}
.layout.dark { background: #1a1a1a; color: #ddd; }
.layout.light { background: #f3f3f3; color: #222; }
.nav-rail {
  border-right: 1px solid rgba(127, 127, 127, 0.16);
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 14px 0;
}
.layout.dark .nav-rail { background: #2a2a2a; }
.layout.light .nav-rail { background: #ededed; }
.nav-top {
  margin-bottom: 16px;
}
.nav-center {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 14px;
}
.nav-bottom {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  padding-bottom: 8px;
}
.nav-btn {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  border: none;
  background: transparent;
  color: inherit;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0.65;
  transition: background 0.15s, opacity 0.15s;
}
.nav-btn:hover {
  background: rgba(127, 127, 127, 0.18);
  opacity: 1;
}
.nav-btn.active {
  opacity: 1;
  color: var(--n-primary-color, #07c160);
}
.nav-btn.small {
  width: 36px;
  height: 36px;
  border-radius: 8px;
}
.dot {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  display: inline-block;
}
.color-grid {
  display: grid;
  grid-template-columns: repeat(3, 28px);
  gap: 8px;
}
.color-cell {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  border: 2px solid transparent;
  cursor: pointer;
}
.color-cell.chosen {
  outline: 2px solid var(--n-text-color, #ccc);
  outline-offset: 2px;
}
.middle-pane {
  border-right: 1px solid rgba(127, 127, 127, 0.16);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.layout.dark .middle-pane { background: #1f1f1f; }
.layout.light .middle-pane { background: #f7f7f7; }
.content-pane {
  display: flex;
  flex-direction: column;
  min-width: 0;
  min-height: 0;
}
.layout.dark .content-pane { background: #2b2b2b; }
.layout.light .content-pane { background: #ffffff; }
.empty {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}
.empty-inner {
  color: rgba(160, 160, 160, 0.7);
  font-size: 14px;
}
.middle-stub {
  padding: 24px;
}
.middle-stub-title {
  font-size: 17px;
  font-weight: 600;
}
</style>
