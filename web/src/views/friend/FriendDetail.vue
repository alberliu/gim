<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { NButton, NInput, useMessage } from 'naive-ui'
import Avatar from '@/components/Avatar.vue'
import { friendClient, userClient } from '@/api/transport'
import { ensurePeerConversation } from '@/api/messageSync'

const props = defineProps<{ userId: string }>()
const emit = defineEmits<{ (e: 'open-chat', convKey: string): void }>()
const message = useMessage()

const profile = ref<{ nickname: string; avatarUrl: string; sex: number } | null>(null)
const friendInfo = ref<{ remarks: string; phone?: string } | null>(null)

async function load() {
  try {
    const r = await userClient.getUser({ userId: BigInt(props.userId) })
    if (r.user) {
      profile.value = { nickname: r.user.nickname, avatarUrl: r.user.avatarUrl, sex: r.user.sex }
    }
    const fr = await friendClient.getFriends({})
    const f = fr.friends.find((x) => x.userId.toString() === props.userId)
    if (f) friendInfo.value = { remarks: f.remarks, phone: f.phoneNumber }
  } catch (e) {
    message.error('加载失败')
  }
}

watch(() => props.userId, load)
onMounted(load)

async function onSendMessage() {
  if (!profile.value) return
  const key = await ensurePeerConversation(props.userId, friendInfo.value?.remarks || profile.value.nickname, profile.value.avatarUrl)
  emit('open-chat', key)
}
</script>

<template>
  <div class="detail">
    <div class="card">
      <div class="head">
        <Avatar :src="profile?.avatarUrl" :name="profile?.nickname" :size="64" :rounded="10" />
        <div class="head-info">
          <div class="name">{{ friendInfo?.remarks || profile?.nickname || '加载中…' }}</div>
          <div class="sub" v-if="friendInfo?.remarks">昵称：{{ profile?.nickname }}</div>
          <div class="sub">UID: {{ userId }}</div>
        </div>
      </div>
      <div class="actions">
        <n-button type="primary" @click="onSendMessage">发消息</n-button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.detail {
  padding: 24px;
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 100%;
}
.card {
  width: 460px;
  padding: 24px;
  background: rgba(127,127,127,0.06);
  border: 1px solid rgba(127,127,127,0.14);
  border-radius: 10px;
}
.head {
  display: flex;
  gap: 16px;
  align-items: center;
  margin-bottom: 24px;
}
.head-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.name {
  font-size: 18px;
  font-weight: 600;
}
.sub {
  font-size: 12px;
  opacity: 0.65;
}
.actions {
  display: flex;
  justify-content: flex-end;
}
</style>
