<script setup lang="ts">
import { useResumeStore } from '@/stores/resume'
const store = useResumeStore()
const addHL = (id: string, p: any) => store.updateProject(id, { highlights: [...p.highlights, ''] })
const updateHL = (id: string, p: any, i: number, v: string) => { const h = [...p.highlights]; h[i] = v; store.updateProject(id, { highlights: h }) }
const removeHL = (id: string, p: any, i: number) => store.updateProject(id, { highlights: p.highlights.filter((_: any, x: number) => x !== i) })
</script>

<template>
  <div class="section-editor">
    <div v-if="store.resume.projects.length === 0" class="empty-hint">暂无项目经历，点击下方按钮添加</div>
    <div v-for="p in store.resume.projects" :key="p.id" class="entry-card">
      <div class="entry-header"><span class="entry-title">{{ p.name || '新增项目经历' }}</span><button class="btn-remove" @click="store.removeProject(p.id)">删除</button></div>
      <div class="form-row">
        <div class="form-item"><label>项目名称</label><input :value="p.name" @input="store.updateProject(p.id, { name: ($event.target as HTMLInputElement).value })" placeholder="如：JobTracker 求职助手" /></div>
        <div class="form-item"><label>担任角色</label><input :value="p.role" @input="store.updateProject(p.id, { role: ($event.target as HTMLInputElement).value })" placeholder="如：全栈开发" /></div>
      </div>
      <div class="form-row">
        <div class="form-item"><label>项目时间</label><div class="date-range"><input :value="p.startDate" @input="store.updateProject(p.id, { startDate: ($event.target as HTMLInputElement).value })" placeholder="2023-01" /><span>至</span><input :value="p.endDate" @input="store.updateProject(p.id, { endDate: ($event.target as HTMLInputElement).value })" placeholder="2023-06" /></div></div>
        <div class="form-item"><label>项目链接</label><input :value="p.link" @input="store.updateProject(p.id, { link: ($event.target as HTMLInputElement).value })" placeholder="GitHub 或在线地址" /></div>
      </div>
      <div class="form-item full-width"><label>项目描述</label><textarea :value="p.description" @input="store.updateProject(p.id, { description: ($event.target as HTMLTextAreaElement).value })" placeholder="简要描述项目内容..." rows="2"></textarea></div>
      <div class="form-item full-width">
        <label>项目亮点</label>
        <div v-for="(hl, i) in p.highlights" :key="i" class="hl-row"><span class="bullet">•</span><input :value="hl" @input="updateHL(p.id, p, i, ($event.target as HTMLInputElement).value)" placeholder="描述一项技术亮点或成果..." /><button class="btn-icon" @click="removeHL(p.id, p, i)">×</button></div>
        <button class="btn-add-sm" @click="addHL(p.id, p)">+ 添加亮点</button>
      </div>
    </div>
    <button class="btn-add" @click="store.addProject()">+ 添加项目经历</button>
  </div>
</template>

<style scoped>
.section-editor{display:flex;flex-direction:column;gap:12px}.entry-card{background:var(--bg);border:1px solid var(--border);border-radius:8px;padding:16px}.entry-header{display:flex;justify-content:space-between;align-items:center;margin-bottom:12px}.entry-title{font-weight:600;font-size:14px}.btn-remove{background:none;border:none;color:#ef4444;cursor:pointer;font-size:13px;padding:4px 8px;border-radius:4px}.btn-remove:hover{background:#fef2f2}.form-row{display:grid;grid-template-columns:1fr 1fr;gap:12px;margin-bottom:8px}.form-item{display:flex;flex-direction:column;gap:4px}.form-item.full-width{grid-column:span 2}.date-range{display:flex;align-items:center;gap:8px}.date-range input{flex:1}.date-range span{color:var(--text-3)}label{font-size:13px;font-weight:500;color:var(--text-2)}input,textarea{padding:8px 12px;border:1px solid var(--border);border-radius:6px;font-size:14px;font-family:inherit;transition:border-color .2s}input:focus,textarea:focus{outline:none;border-color:var(--accent)}textarea{resize:vertical}.hl-row{display:flex;align-items:center;gap:8px;margin-bottom:6px}.bullet{color:var(--accent);font-weight:700}.hl-row input{flex:1}.btn-icon{background:none;border:none;color:#999;cursor:pointer;font-size:18px;padding:0 4px}.btn-icon:hover{color:#ef4444}.btn-add-sm{padding:6px 12px;background:none;border:1px dashed var(--border);border-radius:6px;cursor:pointer;font-size:13px;color:var(--text-3);margin-top:4px}.btn-add-sm:hover{border-color:var(--accent);color:var(--accent)}.btn-add{padding:10px;background:var(--accent-bg);color:var(--accent);border:1px dashed var(--accent);border-radius:8px;cursor:pointer;font-size:14px;transition:all .2s}.btn-add:hover{background:var(--accent);color:#fff}.empty-hint{text-align:center;color:var(--text-3);padding:20px;font-size:14px}
</style>