<script setup lang="ts">
import { ref } from 'vue'
import { NButton, NCard, NDivider, NInput, NSwitch, useMessage } from 'naive-ui'
import { Sun, Moon } from 'lucide-vue-next'
import Avatar from '@/components/Avatar.vue'
import { auth } from '@/store/auth'
import { theme as themeStore, PRESET_COLORS } from '@/store/theme'
import { uploadFile } from '@/api/upload'
import { userClient } from '@/api/transport'

const message = useMessage()
const nickname = ref(auth.nickname)
const saving = ref(false)
const fileRef = ref<HTMLInputElement>()

async function saveProfile() {
  saving.value = true
  try {
    await userClient.updateUser({
      nickname: nickname.value,
      sex: 0,
      avatarUrl: auth.avatarUrl,
      extra: '',
    })
    auth.nickname = nickname.value
    message.success('已保存')
  } catch (e: any) {
    message.error('保存失败：' + (e?.message ?? e))
  } finally {
    saving.value = false
  }
}

function pickAvatar() {
  fileRef.value?.click()
}

async function onAvatarPicked(e: Event) {
  const t = e.target as HTMLInputElement
  if (!t.files || !t.files.length) return
  const file = t.files[0]
  try {
    const url = await uploadFile(file)
    auth.avatarUrl = url
    await userClient.updateUser({
      nickname: auth.nickname,
      sex: 0,
      avatarUrl: url,
      extra: '',
    })
    message.success('头像已更新')
  } catch (e: any) {
    message.error('上传失败：' + (e?.message ?? e))
  } finally {
    t.value = ''
  }
}
</script>

<template>
  <div class="settings">
    <div class="title">设置</div>

    <n-card title="个人资料" :bordered="false" class="card">
      <div class="profile">
        <div class="avatar-area">
          <Avatar :src="auth.avatarUrl" :name="auth.nickname" :size="72" :rounded="12" />
          <n-button size="tiny" tertiary @click="pickAvatar">更换头像</n-button>
          <input
            ref="fileRef"
            type="file"
            accept="image/*"
            style="display:none"
            @change="onAvatarPicked"
          />
        </div>
        <div class="form">
          <div class="field">
            <div class="lbl">手机号</div>
            <n-input :value="auth.phoneNumber" disabled />
          </div>
          <div class="field">
            <div class="lbl">UID</div>
            <n-input :value="auth.userId" disabled />
          </div>
          <div class="field">
            <div class="lbl">昵称</div>
            <n-input v-model:value="nickname" />
          </div>
          <div class="field-actions">
            <n-button type="primary" :loading="saving" @click="saveProfile">保存</n-button>
          </div>
        </div>
      </div>
    </n-card>

    <n-card title="主题外观" :bordered="false" class="card">
      <div class="theme-row">
        <div class="lbl">显示模式</div>
        <div class="theme-actions">
          <n-button
            :type="themeStore.mode === 'light' ? 'primary' : 'default'"
            tertiary
            size="small"
            @click="themeStore.mode = 'light'"
          >
            <template #icon><Sun :size="14" /></template>
            浅色
          </n-button>
          <n-button
            :type="themeStore.mode === 'dark' ? 'primary' : 'default'"
            tertiary
            size="small"
            @click="themeStore.mode = 'dark'"
          >
            <template #icon><Moon :size="14" /></template>
            深色
          </n-button>
        </div>
      </div>
      <n-divider />
      <div class="theme-row">
        <div class="lbl">主题色</div>
        <div class="colors">
          <button
            v-for="c in PRESET_COLORS"
            :key="c.value"
            class="color-cell"
            :title="c.label"
            :style="{ background: c.value }"
            :class="{ chosen: themeStore.primary === c.value }"
            @click="themeStore.primary = c.value"
          />
        </div>
      </div>
    </n-card>
  </div>
</template>

<style scoped>
.settings {
  padding: 24px 32px;
  overflow-y: auto;
  height: 100%;
}
.title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 16px;
}
.card {
  margin-bottom: 18px;
  background: rgba(127,127,127,0.06);
}
.profile {
  display: flex;
  gap: 24px;
}
.avatar-area {
  display: flex;
  flex-direction: column;
  gap: 10px;
  align-items: center;
}
.form {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.field { display: flex; flex-direction: column; gap: 4px; max-width: 360px; }
.lbl { font-size: 12px; opacity: 0.65; }
.field-actions { margin-top: 4px; }
.theme-row {
  display: flex;
  align-items: center;
  gap: 16px;
}
.theme-actions { display: flex; gap: 6px; }
.colors {
  display: flex;
  gap: 10px;
}
.color-cell {
  width: 26px;
  height: 26px;
  border-radius: 50%;
  border: 2px solid transparent;
  cursor: pointer;
}
.color-cell.chosen {
  outline: 2px solid var(--n-text-color, #ccc);
  outline-offset: 2px;
}
</style>
