<script setup lang="ts">
import { computed, ref } from 'vue'
import { NBadge, NInput, NButton, NDropdown, useDialog } from 'naive-ui'
import { Search, Plus } from 'lucide-vue-next'
import Avatar from '@/components/Avatar.vue'
import { state, deleteConversation } from '@/api/messageSync'
import { formatListTime } from '@/utils/time'

const props = defineProps<{ activeKey: string | null }>()
const emit = defineEmits<{ (e: 'select', key: string): void }>()
const dialog = useDialog()

const search = ref('')
const filtered = computed(() => {
  const q = search.value.trim().toLowerCase()
  const list = state.conversations
  if (!q) return list
  return list.filter((c) => c.name.toLowerCase().includes(q) || c.lastPreview.toLowerCase().includes(q))
})

const ctxKey = ref<string | null>(null)

function onContextMenu(e: MouseEvent, key: string) {
  e.preventDefault()
  ctxKey.value = key
  showDropdown.value = false
  dropdownX.value = e.clientX
  dropdownY.value = e.clientY
  // open on next tick
  requestAnimationFrame(() => {
    showDropdown.value = true
  })
}

const showDropdown = ref(false)
const dropdownX = ref(0)
const dropdownY = ref(0)
const dropdownOptions = [{ key: 'delete', label: '删除会话' }]
function onDropdownSelect(key: string) {
  showDropdown.value = false
  if (key === 'delete' && ctxKey.value) {
    const k = ctxKey.value
    dialog.warning({
      title: '删除会话',
      content: '会话内容仅在本地删除',
      positiveText: '删除',
      negativeText: '取消',
      onPositiveClick: () => deleteConversation(k),
    })
  }
}
</script>

<template>
  <div class="conv-pane">
    <div class="search-bar">
      <n-input v-model:value="search" placeholder="搜索" size="small" round>
        <template #prefix>
          <Search :size="14" />
        </template>
      </n-input>
      <n-button quaternary circle size="small" disabled>
        <Plus :size="16" />
      </n-button>
    </div>
    <div class="list">
      <div
        v-for="c in filtered"
        :key="c.key"
        class="row"
        :class="{ active: props.activeKey === c.key }"
        @click="emit('select', c.key)"
        @contextmenu="onContextMenu($event, c.key)"
      >
        <Avatar :src="c.avatarUrl" :name="c.name" :size="42" :rounded="6" />
        <div class="row-mid">
          <div class="row-top">
            <div class="row-name">{{ c.name }}</div>
            <div class="row-time">{{ formatListTime(c.lastTimestamp) }}</div>
          </div>
          <div class="row-bot">
            <div class="row-preview">{{ c.lastPreview || '（暂无消息）' }}</div>
            <n-badge v-if="c.unread > 0" :value="c.unread" :max="99" />
          </div>
        </div>
      </div>
      <div v-if="!filtered.length" class="empty-list">暂无会话<br/>添加好友或创建群组以开始聊天</div>
    </div>

    <n-dropdown
      :show="showDropdown"
      :options="dropdownOptions"
      :x="dropdownX"
      :y="dropdownY"
      placement="bottom-start"
      trigger="manual"
      :on-clickoutside="() => (showDropdown = false)"
      @select="onDropdownSelect"
    />
  </div>
</template>

<style scoped>
.conv-pane {
  display: flex;
  flex-direction: column;
  height: 100%;
}
.search-bar {
  display: flex;
  gap: 8px;
  padding: 12px 12px 8px;
  align-items: center;
}
.list {
  flex: 1;
  overflow-y: auto;
  padding: 4px 0;
}
.row {
  display: flex;
  gap: 10px;
  padding: 10px 12px;
  cursor: pointer;
  border-left: 2px solid transparent;
}
.row:hover {
  background: rgba(127, 127, 127, 0.10);
}
.row.active {
  background: rgba(127, 127, 127, 0.18);
  border-left-color: var(--n-primary-color, #07c160);
}
.row-mid {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  justify-content: center;
}
.row-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}
.row-name {
  font-size: 14px;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
  min-width: 0;
}
.row-time {
  font-size: 11px;
  opacity: 0.55;
  white-space: nowrap;
}
.row-bot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 2px;
  gap: 8px;
}
.row-preview {
  font-size: 12px;
  opacity: 0.6;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
  min-width: 0;
}
.empty-list {
  text-align: center;
  margin-top: 80px;
  color: rgba(127, 127, 127, 0.7);
  font-size: 13px;
  line-height: 1.7;
}
</style>
