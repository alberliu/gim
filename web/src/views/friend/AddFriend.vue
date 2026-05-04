<script setup lang="ts">
import { ref } from 'vue'
import { NButton, NInput, NModal, useMessage } from 'naive-ui'
import Avatar from '@/components/Avatar.vue'
import { friendClient, userClient } from '@/api/transport'

interface UserVM {
  userId: string
  nickname: string
  avatarUrl: string
}

const message = useMessage()
const keyword = ref('')
const results = ref<UserVM[]>([])
const loading = ref(false)
const showRequest = ref(false)
const target = ref<UserVM | null>(null)
const remarks = ref('')
const description = ref('')

async function search() {
  if (!keyword.value.trim()) return
  loading.value = true
  try {
    const r = await userClient.searchUser({ key: keyword.value.trim() })
    results.value = (r.users || []).map((u) => ({
      userId: u.userId.toString(),
      nickname: u.nickname,
      avatarUrl: u.avatarUrl,
    }))
    if (!results.value.length) message.info('没有找到用户')
  } catch (e: any) {
    message.error('搜索失败：' + (e?.message ?? e))
  } finally {
    loading.value = false
  }
}

function startAdd(u: UserVM) {
  target.value = u
  remarks.value = ''
  description.value = '我是' + u.nickname
  showRequest.value = true
}

async function submit() {
  if (!target.value) return
  try {
    await friendClient.add({
      friendId: BigInt(target.value.userId),
      remarks: remarks.value,
      description: description.value,
    })
    message.success('好友申请已发送')
    showRequest.value = false
  } catch (e: any) {
    message.error('发送失败：' + (e?.message ?? e))
  }
}
</script>

<template>
  <div class="page">
    <div class="title">添加好友</div>
    <div class="search-row">
      <n-input
        v-model:value="keyword"
        placeholder="输入手机号或昵称搜索"
        @keyup.enter="search"
      />
      <n-button type="primary" :loading="loading" @click="search">搜索</n-button>
    </div>
    <div class="results">
      <div v-for="u in results" :key="u.userId" class="row">
        <Avatar :src="u.avatarUrl" :name="u.nickname" :size="42" :rounded="6" />
        <div class="row-info">
          <div class="row-name">{{ u.nickname }}</div>
          <div class="row-uid">UID: {{ u.userId }}</div>
        </div>
        <n-button size="small" type="primary" @click="startAdd(u)">添加</n-button>
      </div>
    </div>

    <n-modal
      v-model:show="showRequest"
      preset="dialog"
      title="发送好友申请"
      positive-text="发送"
      negative-text="取消"
      @positive-click="submit"
    >
      <div style="margin: 6px 0 4px;">备注（可选）</div>
      <n-input v-model:value="remarks" placeholder="为对方设置备注" />
      <div style="margin: 12px 0 4px;">附加描述</div>
      <n-input
        v-model:value="description"
        type="textarea"
        :autosize="{ minRows: 2, maxRows: 4 }"
        placeholder="对方可见"
      />
    </n-modal>
  </div>
</template>

<style scoped>
.page {
  padding: 18px 24px;
  height: 100%;
  overflow-y: auto;
}
.title {
  font-size: 17px;
  font-weight: 600;
  margin-bottom: 16px;
}
.search-row {
  display: flex;
  gap: 8px;
  margin-bottom: 18px;
  max-width: 520px;
}
.results {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  border: 1px solid rgba(127,127,127,0.16);
  border-radius: 8px;
}
.row-info { flex: 1; }
.row-name { font-size: 14px; font-weight: 500; }
.row-uid { font-size: 11px; opacity: 0.6; }
</style>
