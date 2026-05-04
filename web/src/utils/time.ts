function pad(n: number) {
  return n < 10 ? `0${n}` : `${n}`
}

export function formatChatTime(ts: number, ref = Date.now()): string {
  if (!ts) return ''
  const d = new Date(ts)
  const now = new Date(ref)
  const sameDay = d.getFullYear() === now.getFullYear() && d.getMonth() === now.getMonth() && d.getDate() === now.getDate()
  if (sameDay) return `${pad(d.getHours())}:${pad(d.getMinutes())}`
  return `${d.getMonth() + 1}/${d.getDate()} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

const WEEK_LABELS = ['星期日', '星期一', '星期二', '星期三', '星期四', '星期五', '星期六']

export function formatTimeLabel(ts: number, ref = Date.now()): string {
  if (!ts) return ''
  const d = new Date(ts)
  const now = new Date(ref)
  const oneDay = 86400000
  const startOf = (x: Date) => new Date(x.getFullYear(), x.getMonth(), x.getDate()).getTime()
  const days = Math.round((startOf(now) - startOf(d)) / oneDay)
  const hm = `${pad(d.getHours())}:${pad(d.getMinutes())}`
  if (days === 0) return hm
  if (days === 1) return `昨天 ${hm}`
  if (days < 7) return `${WEEK_LABELS[d.getDay()]} ${hm}`
  return `${d.getFullYear()}/${d.getMonth() + 1}/${d.getDate()} ${hm}`
}

export function formatListTime(ts: number, ref = Date.now()): string {
  if (!ts) return ''
  const d = new Date(ts)
  const now = new Date(ref)
  const oneDay = 86400000
  const startOf = (x: Date) => new Date(x.getFullYear(), x.getMonth(), x.getDate()).getTime()
  const days = Math.round((startOf(now) - startOf(d)) / oneDay)
  if (days === 0) return `${pad(d.getHours())}:${pad(d.getMinutes())}`
  if (days === 1) return '昨天'
  if (days < 7) return WEEK_LABELS[d.getDay()]
  return `${d.getMonth() + 1}/${d.getDate()}`
}
