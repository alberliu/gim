<script setup lang="ts">
import { ref } from 'vue'
import { NCard, NInput, NButton, NSpace, useMessage } from 'naive-ui'
import { userClient } from '@/api/transport'
import { auth } from '@/store/auth'
import { DeviceType } from '@/gen/proto/logic/device.int_pb'

const message = useMessage()

const phone = ref('')
const code = ref('')
const cdLeft = ref(0)
const submitting = ref(false)
let cdTimer: number | null = null

function startCountdown() {
  cdLeft.value = 60
  if (cdTimer) clearInterval(cdTimer)
  cdTimer = window.setInterval(() => {
    cdLeft.value--
    if (cdLeft.value <= 0 && cdTimer) {
      clearInterval(cdTimer)
      cdTimer = null
    }
  }, 1000)
}

function sendCode() {
  if (!/^1\d{10}$/.test(phone.value)) {
    message.warning('请输入正确的手机号')
    return
  }
  message.success('验证码已发送（测试码：000000）')
  startCountdown()
}

async function submit() {
  if (!/^1\d{10}$/.test(phone.value)) {
    message.warning('请输入正确的手机号')
    return
  }
  if (!code.value) {
    message.warning('请输入验证码')
    return
  }
  submitting.value = true
  try {
    const r = await userClient.signIn({
      phoneNumber: phone.value,
      code: code.value,
      device: {
        type: DeviceType.DT_WEB,
        brand: 'Web',
        model: navigator.userAgent.slice(0, 64),
        systemVersion: '1',
        sdkVersion: '1',
        pushToken: '',
      },
    })
    auth.userId = r.userId.toString()
    auth.deviceId = r.deviceId.toString()
    auth.token = r.token
    auth.phoneNumber = phone.value
    // Pull profile
    try {
      const u = await userClient.getUser({ userId: r.userId })
      if (u.user) {
        auth.nickname = u.user.nickname || `用户${r.userId}`
        auth.avatarUrl = u.user.avatarUrl
      }
    } catch (e) {
      console.warn('getUser after sign in failed', e)
    }
    if (!auth.nickname) auth.nickname = `用户${r.userId}`
    message.success(r.isNew ? '注册并登录成功' : '登录成功')
    // App.vue watches isLoggedIn() and bootstraps automatically.
  } catch (e: any) {
    console.error(e)
    message.error('登录失败：' + (e?.message ?? e))
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="login-page">
    <n-card class="login-card" :bordered="false">
      <div class="login-brand">
        <div class="login-logo">G</div>
        <div class="login-title">GIM</div>
        <div class="login-subtitle">Markdown 即时通讯</div>
      </div>
      <n-space vertical size="large">
        <n-input v-model:value="phone" placeholder="手机号" maxlength="11" />
        <div class="code-row">
          <n-input v-model:value="code" placeholder="验证码（测试：000000）" maxlength="6" />
          <n-button :disabled="cdLeft > 0" @click="sendCode" tertiary>
            {{ cdLeft > 0 ? `${cdLeft}s` : '获取验证码' }}
          </n-button>
        </div>
        <n-button type="primary" block :loading="submitting" @click="submit">登录 / 注册</n-button>
        <div class="login-hint">未注册手机号将自动创建新账号</div>
      </n-space>
    </n-card>
  </div>
</template>

<style scoped>
.login-page {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}
.login-card {
  width: 360px;
  background: var(--n-card-color, transparent);
  border-radius: 12px;
  padding: 12px 8px;
}
.login-brand {
  text-align: center;
  margin-bottom: 28px;
}
.login-logo {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  margin: 0 auto 14px;
  background: var(--n-primary-color, #07c160);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  font-weight: bold;
}
.login-title {
  font-size: 22px;
  font-weight: 600;
}
.login-subtitle {
  font-size: 12px;
  opacity: 0.6;
  margin-top: 4px;
}
.code-row {
  display: flex;
  gap: 8px;
}
.code-row .n-button {
  white-space: nowrap;
}
.login-hint {
  text-align: center;
  font-size: 12px;
  opacity: 0.5;
}
</style>
