<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { NButton, NCheckbox, NInput, useMessage } from 'naive-ui'
import Avatar from '@/components/Avatar.vue'
import { friendClient, groupClient } from '@/api/transport'
import { auth } from '@/store/auth'
import { ensureGroupConversation } from '@/api/messageSync'
import { GroupMemberType, GroupMemberStatus } from '@/gen/proto/logic/group.int_pb'

interface FriendVM {
  userId: string
  nickname: string
  avatarUrl: string
}

const emit = defineEmits<{ (e: 'created', convKey: string): void }>()
const message = useMessage()

const friends = ref<FriendVM[]>([])
const selected = ref<Set<string>>(new Set())
const groupName = ref('')
const submitting = ref(false)

async function load() {
  const r = await friendClient.getFriends({})
  friends.value = (r.friends || []).map((f) => ({
    userId: f.userId.toString(),
    nickname: f.remarks || f.nickname,
    avatarUrl: f.avatarUrl,
  }))
}

onMounted(load)

function toggle(uid: string) {
  if (selected.value.has(uid)) selected.value.delete(uid)
  else selected.value.add(uid)
  selected.value = new Set(selected.value)
}

const canSubmit = computed(() => groupName.value.trim() && selected.value.size > 0)

async function submit() {
  if (!canSubmit.value || submitting.value) return
  submitting.value = true
  try {
    const memberIds = [...selected.value]
    const members = memberIds.map((uid) => {
      const f = friends.value.find((x) => x.userId === uid)
      return {
        userId: BigInt(uid),
        nickname: f?.nickname || '',
        type: GroupMemberType.GMT_DEFAULT,
        status: GroupMemberStatus.GMS_DEFAULT,
        extra: '',
      }
    })
    members.unshift({
      userId: BigInt(auth.userId),
      nickname: auth.nickname,
      type: GroupMemberType.GMT_DEFAULT,
      status: GroupMemberStatus.GMS_DEFAULT,
      extra: '',
    })
    const r = await groupClient.create({
      group: {
        id: BigInt(0),
        name: groupName.value.trim(),
        avatarUrl: '',
        introduction: '',
        extra: '',
      },
      members,
    })
    const groupId = r.groupId.toString()
    const key = await ensureGroupConversation(groupId, groupName.value.trim(), '')
    message.success('群组已创建')
    emit('created', key)
  } catch (e: any) {
    console.error(e)
    message.error('创建失败：' + (e?.message ?? e))
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="page">
    <div class="title">创建群聊</div>

    <div class="field">
      <div class="lbl">群名称</div>
      <n-input v-model:value="groupName" placeholder="请输入群名称" maxlength="32" show-count />
    </div>

    <div class="field">
      <div class="lbl">选择成员（已选 {{ selected.size }}）</div>
      <div class="members">
        <div
          v-for="f in friends"
          :key="f.userId"
          class="m-row"
          :class="{ on: selected.has(f.userId) }"
          @click="toggle(f.userId)"
        >
          <n-checkbox :checked="selected.has(f.userId)" />
          <Avatar :src="f.avatarUrl" :name="f.nickname" :size="32" :rounded="6" />
          <div class="m-name">{{ f.nickname }}</div>
        </div>
        <div v-if="!friends.length" class="empty">暂无好友，先去添加好友</div>
      </div>
    </div>

    <div class="footer">
      <n-button type="primary" :disabled="!canSubmit" :loading="submitting" @click="submit">创建</n-button>
    </div>
  </div>
</template>

<style scoped>
.page {
  padding: 18px 24px;
  height: 100%;
  overflow-y: auto;
}
.title { font-size: 17px; font-weight: 600; margin-bottom: 16px; }
.field { margin-bottom: 18px; max-width: 520px; }
.lbl { font-size: 12px; opacity: 0.7; margin-bottom: 6px; }
.members {
  border: 1px solid rgba(127,127,127,0.16);
  border-radius: 8px;
  max-height: 360px;
  overflow-y: auto;
}
.m-row {
  display: flex;
  gap: 10px;
  align-items: center;
  padding: 8px 12px;
  cursor: pointer;
}
.m-row:hover { background: rgba(127,127,127,0.10); }
.m-row.on { background: rgba(127,127,127,0.18); }
.m-name { font-size: 13px; }
.empty { padding: 24px; text-align: center; opacity: 0.6; font-size: 12px; }
.footer { max-width: 520px; display: flex; justify-content: flex-end; }
</style>
