<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{ src?: string; name?: string; size?: number; rounded?: number; nineGrid?: string[] }>(),
  { size: 36, rounded: 6, name: '', src: '' },
)

const initial = computed(() => (props.name?.[0] || '?').toUpperCase())

function colorFromName(s: string): string {
  let h = 0
  for (let i = 0; i < s.length; i++) h = (h * 31 + s.charCodeAt(i)) | 0
  const palette = ['#5b8def', '#3aa675', '#e07a5f', '#9b59b6', '#f4a261', '#2a9d8f', '#7d80da', '#d36b6b']
  return palette[Math.abs(h) % palette.length]
}
const bg = computed(() => colorFromName(props.name || 'x'))
const styleObj = computed(() => ({
  width: props.size + 'px',
  height: props.size + 'px',
  borderRadius: props.rounded + 'px',
  fontSize: Math.round(props.size * 0.45) + 'px',
}))

const showNineGrid = computed(() => Array.isArray(props.nineGrid) && props.nineGrid.length > 0 && !props.src)
</script>

<template>
  <div class="avatar-wrap" :style="styleObj">
    <img v-if="src" :src="src" :style="styleObj" class="avatar-img" />
    <div v-else-if="showNineGrid" class="avatar-nine" :style="styleObj">
      <img v-for="(u, i) in (nineGrid || []).slice(0, 9)" :key="i" :src="u" />
      <div v-if="(nineGrid || []).length === 0" class="placeholder">群</div>
    </div>
    <div v-else class="avatar-fallback" :style="{ ...styleObj, background: bg }">
      {{ initial }}
    </div>
  </div>
</template>

<style scoped>
.avatar-wrap {
  display: inline-block;
  vertical-align: middle;
  flex-shrink: 0;
}
.avatar-img {
  display: block;
  object-fit: cover;
}
.avatar-fallback {
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-family: -apple-system, sans-serif;
}
.avatar-nine {
  background: rgba(127, 127, 127, 0.15);
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  grid-template-rows: repeat(3, 1fr);
  gap: 1px;
  overflow: hidden;
}
.avatar-nine img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.placeholder {
  grid-column: 1 / span 3;
  grid-row: 1 / span 3;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #888;
}
</style>
