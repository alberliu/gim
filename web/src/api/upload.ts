import { UPLOAD_URL } from './transport'

export async function uploadFile(file: File | Blob, filename?: string): Promise<string> {
  const fd = new FormData()
  const fname = filename || (file instanceof File ? file.name : 'image.png')
  fd.append('file', file, fname)
  const res = await fetch(UPLOAD_URL, { method: 'POST', body: fd })
  if (!res.ok) throw new Error(`upload failed: ${res.status}`)
  const json = (await res.json()) as { code: number; message: string; data: { url: string } }
  if (json.code !== 0) throw new Error(json.message || 'upload failed')
  return json.data.url
}
