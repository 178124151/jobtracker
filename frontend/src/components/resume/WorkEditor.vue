<script setup lang="ts">
import { useResumeStore } from '@/stores/resume'
const store = useResumeStore()
const addHL = (id: string, w: any) => store.updateWork(id, { highlights: [...w.highlights, ''] })
const updateHL = (id: string, w: any, i: number, v: string) => { const h = [...w.highlights]; h[i] = v; store.updateWork(id, { highlights: h }) }
const removeHL = (id: string, w: any, i: number) => store.updateWork(id, { highlights: w.highlights.filter((_: any, x: number) => x !== i) })
</script>

<template>
  <div class="section-editor">
    <div v-if="store.resume.work.length === 0" class="empty-hint">暂无工作经历，点击下方按钮添加</div>
    <div v-for="w in store.resume.work" :key="w.id" class="entry-card">
      <div class="entry-header"><span class="entry-title">{{ w.company || '新增工作经历' }}</span><button class="btn-remove" @click="store.removeWork(w.id)">删除</button></div>
      <div class="form-row">
        <div class="form-item"><label>公司名称</label><input :value="w.company" @input="store.updateWork(w.id, { company: ($event.target as HTMLInputElement).value })" placeholder="如：阿里巴巴" /></div>
        <div class="form-item"><label>职位</label><input :value="w.position" @input="store.updateWork(w.id, { position: ($event.target as HTMLInputElement).value })" placeholder="如：SRE 工程师" /></div>
      </div>
      <div class="form-item"><label>在职时间</label><div class="date-range"><input :value="w.startDate" @input="store.updateWork(w.id, { startDate: ($event.target as HTMLInputElement).value })" placeholder="2022-01" /><span>至</span><input :value="w.endDate" @input="store.updateWork(w.id, { endDate: ($event.target as HTMLInputElement).value })" placeholder="至今" /></div></div>
      <div class="form-item full-width"><label>工作描述</label><textarea :value="w.description" @input="store.updateWork(w.id, { description: ($event.target as HTMLTextAreaElement).value })" placeholder="描述你的主要工作内容..." rows="2"></textarea></div>
      <div class="form-item full-width">
        <label>工作成果</label>
        <div v-for="(hl, i) in w.highlights" :key="i" class="hl-row"><span class="bullet">•</span><input :value="hl" @input="updateHL(w.id, w, i, ($event.target as HTMLInputElement).value)" placeholder="描述一项工作成果..." /><button class="btn-icon" @click="removeHL(w.id, w, i)">×</button></div>
        <button class="btn-add-sm" @click="addHL(w.id, w)">+ 添加成果</button>
      </div>
    </div>
    <button class="btn-add" @click="store.addWork()">+ 添加工作经历</button>
  </div>
</template>

<style scoped>
.section-editor{display:flex;flex-direction:column;gap:12px}.entry-card{background:var(--bg);border:1px solid var(--border);border-radius:8px;padding:16px}.entry-header{display:flex;justify-content:space-between;align-items:center;margin-bottom:12px}.entry-title{font-weight:600;font-size:14px}.btn-remove{background:none;border:none;color:#ef4444;cursor:pointer;font-size:13px;padding:4px 8px;border-radius:4px}.btn-remove:hover{background:#fef2f2}.form-row{display:grid;grid-template-columns:1fr 1fr;gap:12px;margin-bottom:8px}.form-item{display:flex;flex-direction:column;gap:4px}.form-item.full-width{grid-column:span 2}.date-range{display:flex;align-items:center;gap:8px}.date-range input{flex:1}.date-range span{color:var(--text-3)}label{font-size:13px;font-weight:500;color:var(--text-2)}input,textarea{padding:8px 12px;border:1px solid var(--border);border-radius:6px;font-size:14px;font-family:inherit;transition:border-color .2s}input:focus,textarea:focus{outline:none;border-color:var(--accent)}textarea{resize:vertical}.hl-row{display:flex;align-items:center;gap:8px;margin-bottom:6px}.bullet{color:var(--accent);font-weight:700}.hl-row input{flex:1}.btn-icon{background:none;border:none;color:#999;cursor:pointer;font-size:18px;padding:0 4px}.btn-icon:hover{color:#ef4444}.btn-add-sm{padding:6px 12px;background:none;border:1px dashed var(--border);border-radius:6px;cursor:pointer;font-size:13px;color:var(--text-3);margin-top:4px}.btn-add-sm:hover{border-color:var(--accent);color:var(--accent)}.btn-add{padding:10px;background:var(--accent-bg);color:var(--accent);border:1px dashed var(--accent);border-radius:8px;cursor:pointer;font-size:14px;transition:all .2s}.btn-add:hover{background:var(--accent);color:#fff}.empty-hint{text-align:center;color:var(--text-3);padding:20px;font-size:14px}
</style>