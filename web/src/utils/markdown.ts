// Tiny Markdown renderer — supports headings, lists, code, bold/italic, links, images,
// blockquotes, paragraphs. Output is sanitized HTML.

function escapeHtml(s: string): string {
  return s
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')
}

function inline(s: string): string {
  let out = escapeHtml(s)
  // images ![alt](url)
  out = out.replace(/!\[([^\]]*)\]\(([^)\s]+)(?:\s+"([^"]+)")?\)/g, (_, alt, url, _title) => {
    return `<img src="${url}" alt="${alt}" />`
  })
  // links [text](url)
  out = out.replace(/\[([^\]]+)\]\(([^)\s]+)\)/g, (_, text, url) => {
    return `<a href="${url}" target="_blank" rel="noopener noreferrer">${text}</a>`
  })
  // code `x`
  out = out.replace(/`([^`]+)`/g, (_, c) => `<code>${c}</code>`)
  // bold **x**
  out = out.replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>')
  // italic *x*
  out = out.replace(/(^|[^*])\*([^*]+)\*/g, '$1<em>$2</em>')
  // line breaks
  out = out.replace(/\n/g, '<br/>')
  return out
}

export function renderMarkdown(src: string): string {
  if (!src) return ''
  const lines = src.split(/\r?\n/)
  const out: string[] = []
  let i = 0
  while (i < lines.length) {
    const line = lines[i]
    // fenced code
    if (/^```/.test(line)) {
      const lang = line.slice(3).trim()
      const buf: string[] = []
      i++
      while (i < lines.length && !/^```/.test(lines[i])) {
        buf.push(lines[i])
        i++
      }
      i++ // skip closing fence
      out.push(`<pre><code class="lang-${escapeHtml(lang)}">${escapeHtml(buf.join('\n'))}</code></pre>`)
      continue
    }
    // heading
    const h = /^(#{1,6})\s+(.*)$/.exec(line)
    if (h) {
      const level = h[1].length
      out.push(`<h${level}>${inline(h[2])}</h${level}>`)
      i++
      continue
    }
    // blockquote
    if (/^>\s?/.test(line)) {
      const buf: string[] = []
      while (i < lines.length && /^>\s?/.test(lines[i])) {
        buf.push(lines[i].replace(/^>\s?/, ''))
        i++
      }
      out.push(`<blockquote>${inline(buf.join('\n'))}</blockquote>`)
      continue
    }
    // unordered list
    if (/^\s*[-*]\s+/.test(line)) {
      const buf: string[] = []
      while (i < lines.length && /^\s*[-*]\s+/.test(lines[i])) {
        buf.push(`<li>${inline(lines[i].replace(/^\s*[-*]\s+/, ''))}</li>`)
        i++
      }
      out.push(`<ul>${buf.join('')}</ul>`)
      continue
    }
    // ordered list
    if (/^\s*\d+\.\s+/.test(line)) {
      const buf: string[] = []
      while (i < lines.length && /^\s*\d+\.\s+/.test(lines[i])) {
        buf.push(`<li>${inline(lines[i].replace(/^\s*\d+\.\s+/, ''))}</li>`)
        i++
      }
      out.push(`<ol>${buf.join('')}</ol>`)
      continue
    }
    // empty line
    if (/^\s*$/.test(line)) {
      i++
      continue
    }
    // paragraph
    const buf: string[] = []
    while (i < lines.length && !/^\s*$/.test(lines[i]) && !/^(#{1,6})\s+/.test(lines[i]) && !/^```/.test(lines[i]) && !/^>\s?/.test(lines[i]) && !/^\s*[-*]\s+/.test(lines[i]) && !/^\s*\d+\.\s+/.test(lines[i])) {
      buf.push(lines[i])
      i++
    }
    out.push(`<p>${inline(buf.join('\n'))}</p>`)
  }
  return out.join('')
}

export function summarizeMarkdown(src: string, max = 60): string {
  if (!src) return ''
  // Strip image markdown to a placeholder, then strip the rest of markdown for preview.
  let s = src.replace(/!\[[^\]]*\]\([^)]+\)/g, '[图片]')
  s = s.replace(/```[\s\S]*?```/g, '[代码]')
  s = s.replace(/`([^`]+)`/g, '$1')
  s = s.replace(/\[([^\]]+)\]\([^)]+\)/g, '$1')
  s = s.replace(/[*_#>\-]+/g, '')
  s = s.replace(/\s+/g, ' ').trim()
  if (s.length > max) s = s.slice(0, max) + '…'
  return s
}
