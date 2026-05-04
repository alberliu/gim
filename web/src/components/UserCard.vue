<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { NButton, useMessage } from 'naive-ui'
import Avatar from '@/components/Avatar.vue'
import { friendClient, userClient } from '@/api/transport'
import { ensurePeerConversation, state } from '@/api/messageSync'
import { auth } from '@/store/auth'

const props = defineProps<{ userId: string }>()
const emit = defineEmits<{ (e: 'open-chat', key: string): void }>()
const message = useMessage()

const loading = ref(true)
const profile = ref<{ nickname: string; avatarUrl: string; sex: number; phone?: string } | null>(null)
const isFriend = ref(false)

async function load() {
  loading.value = true
  try {
    const r = await userClient.getUser({ userId: BigInt(props.userId) })
    if (r.user) {
      profile.value = { nickname: r.user.nickname, avatarUrl: r.user.avatarUrl, sex: r.user.sex }
      state.profiles[props.userId] = { nickname: r.user.nickname, avatarUrl: r.user.avatarUrl }
    }
    try {
      const fr = await friendClient.getFriends({})
      isFriend.value = fr.friends.some((f) => f.userId.toString() === props.userId)
    } catch {}
  } catch (e: any) {
    message.error('获取用户信息失败')
  } finally {
    loading.value = false
  }
}

onMounted(load)

async function startChat() {
  if (!profile.value) return
  const key = await ensurePeerConversation(props.userId, profile.value.nickname, profile.value.avatarUrl)
  emit('open-chat', key)
}

const adding = ref(false)
async function addFriend() {
  adding.value = true
  try {
    await friendClient.add({ friendId: BigInt(props.userId), remarks: '', description: '' })
    message.success('已发送好友申请')
  } catch (e: any) {
    message.error('申请失败：' + (e?.message ?? e))
  } finally {
    adding.value = false
  }
}

const isSelf = props.userId === auth.userId
</script>

<template>
  <div class="user-card user-text">
    <div class="row">
      <Avatar :src="profile?.avatarUrl" :name="profile?.nickname" :size="48" :rounded="8" />
      <div class="info">
        <div class="name">{{ profile?.nickname || '加载中…' }}</div>
        <div class="sub">UID: {{ userId }}</div>
      </div>
    </div>
    <div class="footer">
      <n-button v-if="!isSelf" size="small" type="primary" @click="startChat">发消息</n-button>
      <n-button v-if="!isSelf && !isFriend" size="small" :loading="adding" tertiary @click="addFriend">添加好友</n-button>
    </div>
  </div>
</template>

<style scoped>
.user-card {
  width: 240px;
  padding: 10px 4px;
}
.row {
  display: flex;
  gap: 12px;
  align-items: center;
  margin-bottom: 12px;
}
.info {
  display: flex;
  flex-direction: column;
}
.name {
  font-size: 15px;
  font-weight: 600;
}
.sub {
  font-size: 11px;
  opacity: 0.6;
  margin-top: 2px;
}
.footer {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
}
</style>
