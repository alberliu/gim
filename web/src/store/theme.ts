import { reactive, computed, watch } from 'vue'
import { darkTheme, type GlobalTheme, type GlobalThemeOverrides } from 'naive-ui'

export interface ThemeState {
  mode: 'light' | 'dark'
  primary: string
}

export const PRESET_COLORS = [
  { label: '微信绿', value: '#07c160' },
  { label: '蔚蓝', value: '#1989fa' },
  { label: '番茄', value: '#ff6b35' },
  { label: '紫罗兰', value: '#8a4fff' },
  { label: '玫瑰红', value: '#e85d75' },
  { label: '深空灰', value: '#5a6c7d' },
]

const STORAGE_KEY = 'gim.theme'
const defaults: ThemeState = { mode: 'dark', primary: '#07c160' }

function load(): ThemeState {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (raw) return { ...defaults, ...JSON.parse(raw) }
  } catch {}
  return { ...defaults }
}

export const theme = reactive<ThemeState>(load())

watch(
  theme,
  (v) => localStorage.setItem(STORAGE_KEY, JSON.stringify(v)),
  { deep: true },
)

function darken(hex: string, amount = 0.15): string {
  const m = hex.replace('#', '')
  const r = parseInt(m.slice(0, 2), 16)
  const g = parseInt(m.slice(2, 4), 16)
  const b = parseInt(m.slice(4, 6), 16)
  const f = (x: number) => Math.max(0, Math.min(255, Math.round(x * (1 - amount))))
  const toHex = (x: number) => x.toString(16).padStart(2, '0')
  return `#${toHex(f(r))}${toHex(f(g))}${toHex(f(b))}`
}

function lighten(hex: string, amount = 0.15): string {
  const m = hex.replace('#', '')
  const r = parseInt(m.slice(0, 2), 16)
  const g = parseInt(m.slice(2, 4), 16)
  const b = parseInt(m.slice(4, 6), 16)
  const f = (x: number) => Math.max(0, Math.min(255, Math.round(x + (255 - x) * amount)))
  const toHex = (x: number) => x.toString(16).padStart(2, '0')
  return `#${toHex(f(r))}${toHex(f(g))}${toHex(f(b))}`
}

export const naiveTheme = computed<GlobalTheme | null>(() => (theme.mode === 'dark' ? darkTheme : null))

export const naiveThemeOverrides = computed<GlobalThemeOverrides>(() => ({
  common: {
    primaryColor: theme.primary,
    primaryColorHover: lighten(theme.primary, 0.1),
    primaryColorPressed: darken(theme.primary, 0.15),
    primaryColorSuppl: theme.primary,
  },
}))
