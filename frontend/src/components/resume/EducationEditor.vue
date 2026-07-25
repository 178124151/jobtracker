<script setup lang="ts">
import { useResumeStore } from '@/stores/resume'
const store = useResumeStore()
</script>

<template>
  <div class="section-editor">
    <div v-if="store.resume.education.length === 0" class="empty-hint">暂无教育经历，点击下方按钮添加</div>
    <div v-for="edu in store.resume.education" :key="edu.id" class="entry-card">
      <div class="entry-header">
        <span class="entry-title">{{ edu.school || '新增教育经历' }}</span>
        <button class="btn-remove" @click="store.removeEducation(edu.id)">删除</button>
      </div>
      <div class="form-row">
        <div class="form-item"><label>学校名称</label><input :value="edu.school" @input="store.updateEducation(edu.id, { school: ($event.target as HTMLInputElement).value })" placeholder="如：清华大学" /></div>
        <div class="form-item"><label>学历</label><input :value="edu.degree" @input="store.updateEducation(edu.id, { degree: ($event.target as HTMLInputElement).value })" placeholder="如：本科 / 硕士 / 博士" /></div>
      </div>
      <div class="form-row">
        <div class="form-item"><label>专业</label><input :value="edu.major" @input="store.updateEducation(edu.id, { major: ($event.target as HTMLInputElement).value })" placeholder="如：计算机科学与技术" /></div>
        <div class="form-item"><label>在读时间</label><div class="date-range"><input :value="edu.startDate" @input="store.updateEducation(edu.id, { startDate: ($event.target as HTMLInputElement).value })" placeholder="2020-09" /><span>至</span><input :value="edu.endDate" @input="store.updateEducation(edu.id, { endDate: ($event.target as HTMLInputElement).value })" placeholder="2024-06" /></div></div>
      </div>
      <div class="form-item full-width"><label>补充说明</label><textarea :value="edu.description" @input="store.updateEducation(edu.id, { description: ($event.target as HTMLTextAreaElement).value })" placeholder="如：GPA、获奖经历、主修课程等" rows="2"></textarea></div>
    </div>
    <button class="btn-add" @click="store.addEducation()">+ 添加教育经历</button>
  </div>
</template>

<style scoped>
.section-editor{display:flex;flex-direction:column;gap:12px}.entry-card{background:var(--bg);border:1px solid var(--border);border-radius:8px;padding:16px}.entry-header{display:flex;justify-content:space-between;align-items:center;margin-bottom:12px}.entry-title{font-weight:600;font-size:14px}.btn-remove{background:none;border:none;color:#ef4444;cursor:pointer;font-size:13px;padding:4px 8px;border-radius:4px}.btn-remove:hover{background:#fef2f2}.form-row{display:grid;grid-template-columns:1fr 1fr;gap:12px}.form-item{display:flex;flex-direction:column;gap:4px}.form-item.full-width{grid-column:span 2}.date-range{display:flex;align-items:center;gap:8px}.date-range input{flex:1}.date-range span{color:var(--text-3)}.label{font-size:13px;font-weight:500;color:var(--text-2)}label{font-size:13px;font-weight:500;color:var(--text-2)}input,textarea{padding:8px 12px;border:1px solid var(--border);border-radius:6px;font-size:14px;font-family:inherit;transition:border-color .2s}input:focus,textarea:focus{outline:none;border-color:var(--accent)}textarea{resize:vertical}.btn-add{padding:10px;background:var(--accent-bg);color:var(--accent);border:1px dashed var(--accent);border-radius:8px;cursor:pointer;font-size:14px;transition:all .2s}.btn-add:hover{background:var(--accent);color:#fff}.empty-hint{text-align:center;color:var(--text-3);padding:20px;font-size:14px}
</style>