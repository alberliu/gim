<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { NBadge, NInput } from 'naive-ui'
import { Search, UserPlus, Users, UsersRound } from 'lucide-vue-next'
import Avatar from '@/components/Avatar.vue'
import { friendClient } from '@/api/transport'
import { state } from '@/api/messageSync'

interface FriendVM {
  userId: string
  nickname: string
  avatarUrl: string
  remarks: string
}

const props = defineProps<{ active: string | null }>()
const emit = defineEmits<{
  (e: 'select', userId: string): void
  (e: 'open-new'): void
  (e: 'open-add'): void
  (e: 'open-create-group'): void
}>()

const friends = ref<FriendVM[]>([])
const search = ref('')
const loading = ref(false)

async function refresh() {
  loading.value = true
  try {
    const r = await friendClient.getFriends({})
    friends.value = (r.friends || []).map((f) => ({
      userId: f.userId.toString(),
      nickname: f.remarks || f.nickname,
      avatarUrl: f.avatarUrl,
      remarks: f.remarks,
    }))
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

defineExpose({ refresh })

onMounted(refresh)

const filtered = computed(() => {
  const q = search.value.trim().toLowerCase()
  if (!q) return friends.value
  return friends.value.filter((f) => f.nickname.toLowerCase().includes(q))
})

const groupedByInitial = computed(() => {
  const map = new Map<string, FriendVM[]>()
  for (const f of filtered.value) {
    const c = (f.nickname[0] || '#').toUpperCase()
    if (!map.has(c)) map.set(c, [])
    map.get(c)!.push(f)
  }
  return Array.from(map.entries()).sort(([a], [b]) => a.localeCompare(b))
})
</script>

<template>
  <div class="fl-pane">
    <div class="search-bar">
      <n-input v-model:value="search" placeholder="搜索好友" size="small" round>
        <template #prefix>
          <Search :size="14" />
        </template>
      </n-input>
    </div>
    <div class="entry-list">
      <div class="entry" @click="emit('open-new')">
        <div class="entry-icon" style="background: #ff7e7e"><UserPlus :size="16" color="#fff" /></div>
        <div class="entry-name">新朋友</div>
        <n-badge v-if="state.pendingFriendRequests > 0" :value="state.pendingFriendRequests" />
      </div>
      <div class="entry" @click="emit('open-create-group')">
        <div class="entry-icon" style="background: #5b8def"><UsersRound :size="16" color="#fff" /></div>
        <div class="entry-name">创建群聊</div>
      </div>
      <div class="entry" @click="emit('open-add')">
        <div class="entry-icon" style="background: #3aa675"><Users :size="16" color="#fff" /></div>
        <div class="entry-name">添加好友</div>
      </div>
    </div>
    <div class="divider"><span>我的好友</span></div>
    <div class="list">
      <div v-for="[letter, group] in groupedByInitial" :key="letter">
        <div class="letter">{{ letter }}</div>
        <div
          v-for="f in group"
          :key="f.userId"
          class="row"
          :class="{ active: props.active === f.userId }"
          @click="emit('select', f.userId)"
        >
          <Avatar :src="f.avatarUrl" :name="f.nickname" :size="34" :rounded="6" />
          <div class="row-name">{{ f.nickname }}</div>
        </div>
      </div>
      <div v-if="!loading && !friends.length" class="empty">暂无好友<br/>点击"添加好友"开始</div>
    </div>
  </div>
</template>

<style scoped>
.fl-pane {
  display: flex;
  flex-direction: column;
  height: 100%;
}
.search-bar {
  padding: 12px 12px 8px;
}
.entry-list {
  padding: 4px 0;
}
.entry {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 14px;
  cursor: pointer;
}
.entry:hover {
  background: rgba(127,127,127,0.10);
}
.entry-icon {
  width: 30px;
  height: 30px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.entry-name {
  flex: 1;
  font-size: 13px;
}
.divider {
  padding: 6px 14px;
  font-size: 11px;
  opacity: 0.5;
  border-top: 1px solid rgba(127,127,127,0.16);
  border-bottom: 1px solid rgba(127,127,127,0.16);
  background: rgba(127,127,127,0.05);
}
.list {
  flex: 1;
  overflow-y: auto;
  padding-bottom: 12px;
}
.letter {
  font-size: 11px;
  opacity: 0.5;
  padding: 6px 14px 2px;
}
.row {
  display: flex;
  gap: 10px;
  padding: 7px 14px;
  align-items: center;
  cursor: pointer;
}
.row:hover { background: rgba(127,127,127,0.10); }
.row.active { background: rgba(127,127,127,0.18); }
.row-name {
  font-size: 13px;
  font-weight: 500;
}
.empty {
  text-align: center;
  margin-top: 60px;
  color: rgba(127,127,127,0.7);
  font-size: 12px;
  line-height: 1.7;
}
</style>
