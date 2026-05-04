<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useMessage } from 'naive-ui'
import { uploadFile } from '@/api/upload'

const emit = defineEmits<{ (e: 'send', text: string): void }>()
const editorRef = ref<HTMLDivElement>()
const message = useMessage()

function getMarkdownFromEditor(): string {
  const root = editorRef.value
  if (!root) return ''
  // Walk the DOM and produce Markdown.
  const out: string[] = []
  function walk(node: ChildNode) {
    if (node.nodeType === Node.TEXT_NODE) {
      out.push(node.textContent || '')
      return
    }
    if (node.nodeType !== Node.ELEMENT_NODE) return
    const el = node as HTMLElement
    const tag = el.tagName
    if (tag === 'BR') {
      out.push('\n')
      return
    }
    if (tag === 'IMG') {
      const src = el.getAttribute('src') || ''
      const alt = el.getAttribute('alt') || ''
      if (src.startsWith('blob:') || src.startsWith('data:')) {
        // Skip placeholders that aren't yet uploaded; they'll be replaced.
        out.push(`![${alt}](${src})`)
      } else {
        out.push(`![${alt}](${src})`)
      }
      return
    }
    if (tag === 'DIV' || tag === 'P') {
      // Treat block as a newline break before content
      if (out.length && !out[out.length - 1].endsWith('\n')) out.push('\n')
      el.childNodes.forEach(walk)
      return
    }
    el.childNodes.forEach(walk)
  }
  root.childNodes.forEach(walk)
  return out.join('').replace(/ /g, ' ').trim()
}

function clearEditor() {
  if (editorRef.value) editorRef.value.innerHTML = ''
}

async function handleSend() {
  const text = getMarkdownFromEditor()
  if (!text) return
  // Wait briefly if any blob: image still uploading.
  if (/!\[[^\]]*\]\(blob:/.test(text)) {
    message.warning('图片上传中，请稍候')
    return
  }
  emit('send', text)
  clearEditor()
}

function onKeyDown(e: KeyboardEvent) {
  if (e.key === 'Enter' && !e.shiftKey && !e.isComposing) {
    e.preventDefault()
    handleSend()
  }
}

async function insertImageFile(file: File) {
  const placeholder = document.createElement('img')
  const blobUrl = URL.createObjectURL(file)
  placeholder.src = blobUrl
  placeholder.dataset.uploading = '1'
  placeholder.alt = ''
  insertNodeAtCursor(placeholder)
  try {
    const url = await uploadFile(file)
    placeholder.src = url
    placeholder.dataset.uploading = ''
    URL.revokeObjectURL(blobUrl)
  } catch (e: any) {
    message.error('上传失败：' + (e?.message ?? e))
    placeholder.remove()
  }
}

function insertNodeAtCursor(node: Node) {
  const root = editorRef.value
  if (!root) return
  root.focus()
  const sel = window.getSelection()
  if (!sel) return
  if (!sel.rangeCount || !root.contains(sel.anchorNode)) {
    const range = document.createRange()
    range.selectNodeContents(root)
    range.collapse(false)
    sel.removeAllRanges()
    sel.addRange(range)
  }
  const range = sel.getRangeAt(0)
  range.deleteContents()
  range.insertNode(node)
  range.setStartAfter(node)
  range.setEndAfter(node)
  sel.removeAllRanges()
  sel.addRange(range)
}

async function onPaste(e: ClipboardEvent) {
  if (!e.clipboardData) return
  const items = e.clipboardData.items
  for (let i = 0; i < items.length; i++) {
    const it = items[i]
    if (it.kind === 'file' && it.type.startsWith('image/')) {
      e.preventDefault()
      const file = it.getAsFile()
      if (file) await insertImageFile(file)
      return
    }
  }
  // Fall back to default plain-text paste.
  const text = e.clipboardData.getData('text/plain')
  if (text) {
    e.preventDefault()
    document.execCommand('insertText', false, text)
  }
}

async function onDrop(e: DragEvent) {
  if (!e.dataTransfer) return
  const files = e.dataTransfer.files
  if (files && files.length) {
    e.preventDefault()
    for (let i = 0; i < files.length; i++) {
      const f = files[i]
      if (f.type.startsWith('image/')) await insertImageFile(f)
    }
  }
}

function onDragOver(e: DragEvent) {
  if (e.dataTransfer && Array.from(e.dataTransfer.items || []).some((it) => it.kind === 'file')) {
    e.preventDefault()
  }
}

const fileInputRef = ref<HTMLInputElement>()
function pickImage() {
  fileInputRef.value?.click()
}
async function onFilePicked(e: Event) {
  const t = e.target as HTMLInputElement
  if (!t.files) return
  for (let i = 0; i < t.files.length; i++) {
    const f = t.files[i]
    if (f.type.startsWith('image/')) await insertImageFile(f)
  }
  t.value = ''
}

function applyFormat(prefix: string, suffix: string = prefix) {
  editorRef.value?.focus()
  const sel = window.getSelection()
  if (!sel || !sel.rangeCount) return
  const text = sel.toString() || '文本'
  document.execCommand('insertText', false, `${prefix}${text}${suffix}`)
}

defineExpose({ pickImage })

onMounted(() => {
  // nothing
})

function onClickToolbar(action: string) {
  if (action === 'image') pickImage()
  else if (action === 'bold') applyFormat('**')
  else if (action === 'italic') applyFormat('*')
  else if (action === 'code') applyFormat('`')
  else if (action === 'link') applyFormat('[', '](https://)')
}
</script>

<template>
  <div class="editor-container">
    <div class="toolbar">
      <button class="tb" title="粗体" @click="onClickToolbar('bold')"><b>B</b></button>
      <button class="tb" title="斜体" @click="onClickToolbar('italic')"><i>I</i></button>
      <button class="tb" title="代码" @click="onClickToolbar('code')">&lt;/&gt;</button>
      <button class="tb" title="链接" @click="onClickToolbar('link')">🔗</button>
      <span class="sep" />
      <button class="tb" title="插入图片" @click="onClickToolbar('image')">📷</button>
      <input
        ref="fileInputRef"
        type="file"
        accept="image/*"
        multiple
        style="display:none"
        @change="onFilePicked"
      />
      <span class="hint">Enter 发送 / Shift+Enter 换行</span>
    </div>
    <div
      ref="editorRef"
      class="md-editor"
      contenteditable="true"
      :data-placeholder="'输入消息（支持 Markdown）'"
      @keydown="onKeyDown"
      @paste="onPaste"
      @drop="onDrop"
      @dragover="onDragOver"
    />
    <div class="actions">
      <button class="send-btn" @click="handleSend">发送</button>
    </div>
  </div>
</template>

<style scoped>
.editor-container {
  flex: 0 0 auto;
  border-top: 1px solid rgba(127, 127, 127, 0.16);
  padding: 6px 14px 10px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  user-select: text;
}
.toolbar {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  opacity: 0.85;
}
.tb {
  border: none;
  background: transparent;
  color: inherit;
  width: 28px;
  height: 26px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
}
.tb:hover {
  background: rgba(127, 127, 127, 0.18);
}
.sep {
  width: 1px;
  height: 14px;
  background: rgba(127, 127, 127, 0.4);
  margin: 0 6px;
}
.hint {
  margin-left: auto;
  font-size: 11px;
  opacity: 0.45;
}
.md-editor {
  user-select: text;
  cursor: text;
  padding: 6px 0;
}
.actions {
  display: flex;
  justify-content: flex-end;
}
.send-btn {
  background: var(--n-primary-color, #07c160);
  color: #fff;
  border: none;
  border-radius: 4px;
  padding: 4px 18px;
  font-size: 13px;
  cursor: pointer;
}
.send-btn:hover {
  filter: brightness(1.05);
}
</style>
