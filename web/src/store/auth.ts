import { reactive, watch } from 'vue'

interface AuthState {
  userId: string
  deviceId: string
  token: string
  nickname: string
  avatarUrl: string
  phoneNumber: string
}

const STORAGE_KEY = 'gim.auth'

function load(): AuthState {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (raw) return JSON.parse(raw)
  } catch {}
  return { userId: '', deviceId: '', token: '', nickname: '', avatarUrl: '', phoneNumber: '' }
}

export const auth = reactive<AuthState>(load())

watch(
  auth,
  (v) => {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(v))
  },
  { deep: true },
)

export function logout() {
  auth.userId = ''
  auth.deviceId = ''
  auth.token = ''
  auth.nickname = ''
  auth.avatarUrl = ''
  auth.phoneNumber = ''
  localStorage.removeItem(STORAGE_KEY)
  // Drop persisted state for the previous account too — stay scoped per-user via DB suffix.
}

export function isLoggedIn() {
  return !!auth.token && !!auth.userId
}
