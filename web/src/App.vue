<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { NConfigProvider, NDialogProvider, NMessageProvider, NNotificationProvider, NGlobalStyle, zhCN, dateZhCN } from 'naive-ui'
import { naiveTheme, naiveThemeOverrides, theme as themeStore } from '@/store/theme'
import { auth, isLoggedIn } from '@/store/auth'
import LoginView from '@/views/LoginView.vue'
import MainLayout from '@/views/MainLayout.vue'
import { initDb } from '@/db'
import { bootstrap, shutdown } from '@/api/messageSync'

const ready = ref(false)

async function bootForUser() {
  await initDb(auth.userId)
  await bootstrap()
  ready.value = true
}

onMounted(async () => {
  if (isLoggedIn()) {
    await bootForUser()
  } else {
    ready.value = true
  }
})

// Reactively bootstrap when auth becomes valid. We watch `isLoggedIn()`
// instead of relying on a `@logged-in` event from LoginView because the
// login flow performs an `await` between setting auth and emitting; Vue
// unmounts LoginView during that microtask, so the parent event handler
// never fires.
watch(
  () => isLoggedIn(),
  async (now, was) => {
    if (now && !was) {
      ready.value = false
      shutdown()
      await bootForUser()
    }
  },
)

function onLogout() {
  shutdown()
  auth.userId = ''
  auth.deviceId = ''
  auth.token = ''
  auth.nickname = ''
  auth.avatarUrl = ''
  auth.phoneNumber = ''
  ready.value = true
}

</script>

<template>
  <n-config-provider
    :theme="naiveTheme"
    :theme-overrides="naiveThemeOverrides"
    :locale="zhCN"
    :date-locale="dateZhCN"
  >
    <n-global-style />
    <n-message-provider>
      <n-notification-provider>
        <n-dialog-provider>
          <div class="gim-app">
            <template v-if="ready">
              <LoginView v-if="!isLoggedIn()" />
              <MainLayout v-else @logout="onLogout" />
            </template>
          </div>
        </n-dialog-provider>
      </n-notification-provider>
    </n-message-provider>
  </n-config-provider>
</template>
