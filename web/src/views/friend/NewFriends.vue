<script setup lang="ts">
import { ref } from 'vue'
import { NButton, NInput, NModal, useMessage } from 'naive-ui'
import Avatar from '@/components/Avatar.vue'
import { state, updateFriendRequest, loadFriendRequests } from '@/api/messageSync'
import { friendClient } from '@/api/transport'

const message = useMessage()
const remarksModal = ref(false)
const remarksInput = ref('')
const targetFriendId = ref<string | null>(null)

function startAgree(friendId: string) {
  targetFriendId.value = friendId
  remarksInput.value = ''
  remarksModal.value = true
}

async function confirmAgree() {
  if (!targetFriendId.value) return
  const fid = targetFriendId.value
  try {
    await friendClient.agree({ userId: BigInt(fid), remarks: remarksInput.value })
    await updateFriendRequest(fid, 'agreed')
    message.success('已同意')
    remarksModal.value = false
  } catch (e: any) {
    message.error('同意失败：' + (e?.message ?? e))
  }
}

async function ignore(friendId: string) {
  await updateFriendRequest(friendId, 'ignored')
  message.success('已忽略')
}

loadFriendRequests()
</script>

<template>
  <div class="page">
    <div class="title">新朋友</div>
    <div class="list">
      <div v-for="r in state.friendRequests" :key="r.friendId" class="row">
        <Avatar :src="r.avatarUrl" :name="r.nickname" :size="42" :rounded="6" />
        <div class="row-mid">
          <div class="row-top">
            <div class="row-name">{{ r.nickname }}</div>
            <div class="row-time">{{ new Date(r.receivedAt).toLocaleDateString() }}</div>
          </div>
          <div class="row-desc">{{ r.description || '请求加为好友' }}</div>
        </div>
        <div class="row-actions">
          <template v-if="r.status === 'pending'">
            <n-button size="small" type="primary" @click="startAgree(r.friendId)">同意</n-button>
            <n-button size="small" tertiary @click="ignore(r.friendId)">忽略</n-button>
          </template>
          <template v-else-if="r.status === 'agreed'">
            <span class="status-text">已同意</span>
          </template>
          <template v-else>
            <span class="status-text">已忽略</span>
          </template>
        </div>
      </div>
      <div v-if="!state.friendRequests.length" class="empty">暂无好友申请</div>
    </div>

    <n-modal
      v-model:show="remarksModal"
      preset="dialog"
      title="同意好友申请"
      positive-text="确认"
      negative-text="取消"
      @positive-click="confirmAgree"
    >
      <div style="margin-bottom: 6px">备注（可选）</div>
      <n-input v-model:value="remarksInput" placeholder="为该好友设置备注" />
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
.list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.row {
  display: flex;
  gap: 10px;
  align-items: center;
  padding: 10px 12px;
  border: 1px solid rgba(127,127,127,0.16);
  border-radius: 8px;
}
.row-mid { flex: 1; min-width: 0; }
.row-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}
.row-name { font-size: 14px; font-weight: 600; }
.row-time { font-size: 11px; opacity: 0.55; }
.row-desc { font-size: 12px; opacity: 0.7; margin-top: 2px; }
.row-actions { display: flex; gap: 6px; }
.status-text { color: rgba(127,127,127,0.7); font-size: 12px; }
.empty { color: rgba(127,127,127,0.7); text-align: center; margin-top: 80px; }
</style>
