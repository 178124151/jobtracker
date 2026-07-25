<script setup lang="ts">
import { useResumeStore } from '@/stores/resume'
const store = useResumeStore()
const addItem = (id: string, s: any) => store.updateSkill(id, { items: [...s.items, ''] })
const updateItem = (id: string, s: any, i: number, v: string) => { const items = [...s.items]; items[i] = v; store.updateSkill(id, { items }) }
const removeItem = (id: string, s: any, i: number) => store.updateSkill(id, { items: s.items.filter((_: any, x: number) => x !== i) })
</script>

<template>
  <div class="section-editor">
    <div v-if="store.resume.skills.length === 0" class="empty-hint">暂无技能信息，点击下方按钮添加</div>
    <div v-for="s in store.resume.skills" :key="s.id" class="entry-card">
      <div class="entry-header">
        <input class="cat-input" :value="s.category" @input="store.updateSkill(s.id, { category: ($event.target as HTMLInputElement).value })" placeholder="技能分类（如：编程语言）" />
        <button class="btn-remove" @click="store.removeSkill(s.id)">删除</button>
      </div>
      <div class="tags-wrap">
        <span v-for="(item, i) in s.items" :key="i" class="tag"><input :value="item" @input="updateItem(s.id, s, i, ($event.target as HTMLInputElement).value)" placeholder="技能" /><button @click="removeItem(s.id, s, i)">×</button></span>
        <button class="btn-add-tag" @click="addItem(s.id, s)">+ 添加</button>
      </div>
    </div>
    <button class="btn-add" @click="store.addSkill()">+ 添加技能分类</button>
  </div>
</template>

<style scoped>
.section-editor{display:flex;flex-direction:column;gap:12px}.entry-card{background:var(--bg);border:1px solid var(--border);border-radius:8px;padding:16px}.entry-header{display:flex;justify-content:space-between;align-items:center;margin-bottom:12px;gap:12px}.cat-input{flex:1;padding:8px 12px;border:1px solid var(--border);border-radius:6px;font-size:14px;font-weight:600;font-family:inherit}.cat-input:focus{outline:none;border-color:var(--accent)}.btn-remove{background:none;border:none;color:#ef4444;cursor:pointer;font-size:13px;padding:4px 8px;border-radius:4px;white-space:nowrap}.btn-remove:hover{background:#fef2f2}.tags-wrap{display:flex;flex-wrap:wrap;gap:8px;align-items:center}.tag{display:flex;align-items:center;gap:4px;background:var(--surface);border:1px solid var(--border);border-radius:20px;padding:4px 4px 4px 12px}.tag input{border:none;background:none;width:80px;font-size:13px;padding:2px}.tag input:focus{outline:none}.tag button{background:none;border:none;color:#999;cursor:pointer;font-size:14px;padding:0 4px}.tag button:hover{color:#ef4444}.btn-add-tag{padding:4px 12px;background:none;border:1px dashed var(--border);border-radius:20px;cursor:pointer;font-size:13px;color:var(--text-3)}.btn-add-tag:hover{border-color:var(--accent);color:var(--accent)}.btn-add{padding:10px;background:var(--accent-bg);color:var(--accent);border:1px dashed var(--accent);border-radius:8px;cursor:pointer;font-size:14px;transition:all .2s}.btn-add:hover{background:var(--accent);color:#fff}.empty-hint{text-align:center;color:var(--text-3);padding:20px;font-size:14px}
</style>