<script setup lang="ts">
import { useResumeStore } from '@/stores/resume'
const store = useResumeStore()

const update = (field: string, value: string) => store.updateBasic({ [field]: value })

const handleAvatarUpload = (e: Event) => {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (!file) return
  const reader = new FileReader()
  reader.onload = (ev) => { store.updateBasic({ avatar: ev.target?.result as string }) }
  reader.readAsDataURL(file)
}
</script>

<template>
  <div class="section-editor">
    <div class="avatar-section">
      <div class="avatar-preview" @click="($refs.avatarInput as HTMLInputElement).click()">
        <img v-if="store.resume.basic.avatar" :src="store.resume.basic.avatar" alt="证件照" />
        <div v-else class="avatar-placeholder">
          <span class="icon">📷</span>
          <span>上传证件照</span>
        </div>
      </div>
      <input ref="avatarInput" type="file" accept="image/*" style="display:none" @change="handleAvatarUpload" />
      <p class="avatar-hint">点击上传证件照（建议 1寸/2寸）</p>
    </div>

    <div class="form-row">
      <div class="form-item">
        <label>姓名</label>
        <input :value="store.resume.basic.name" @input="update('name', ($event.target as HTMLInputElement).value)" placeholder="请输入姓名" />
      </div>
      <div class="form-item">
        <label>手机号</label>
        <input :value="store.resume.basic.phone" @input="update('phone', ($event.target as HTMLInputElement).value)" placeholder="请输入手机号" />
      </div>
    </div>
    <div class="form-row">
      <div class="form-item">
        <label>电子邮箱</label>
        <input :value="store.resume.basic.email" @input="update('email', ($event.target as HTMLInputElement).value)" placeholder="请输入邮箱" />
      </div>
      <div class="form-item">
        <label>所在城市</label>
        <input :value="store.resume.basic.city" @input="update('city', ($event.target as HTMLInputElement).value)" placeholder="如：北京、上海" />
      </div>
    </div>
    <div class="form-item full-width">
      <label>个人主页</label>
      <input :value="store.resume.basic.website" @input="update('website', ($event.target as HTMLInputElement).value)" placeholder="如：个人博客、LinkedIn 等" />
    </div>
    <div class="form-item full-width">
      <label>个人简介</label>
      <textarea :value="store.resume.basic.summary" @input="update('summary', ($event.target as HTMLTextAreaElement).value)" placeholder="简要介绍你的职业背景和求职意向..." rows="4"></textarea>
    </div>
  </div>
</template>

<style scoped>
.section-editor{display:flex;flex-direction:column;gap:14px}
.avatar-section{display:flex;flex-direction:column;align-items:center;gap:8px;padding:16px;background:var(--bg);border-radius:8px;border:1px dashed var(--border)}
.avatar-preview{width:120px;height:160px;border-radius:8px;overflow:hidden;cursor:pointer;border:2px dashed var(--border);display:flex;align-items:center;justify-content:center;transition:border-color .2s;background:var(--surface)}
.avatar-preview:hover{border-color:var(--accent)}
.avatar-preview img{width:100%;height:100%;object-fit:cover}
.avatar-placeholder{display:flex;flex-direction:column;align-items:center;gap:6px;color:var(--text-3);font-size:13px}
.avatar-placeholder .icon{font-size:28px}
.avatar-hint{font-size:12px;color:var(--text-3)}
.form-row{display:grid;grid-template-columns:1fr 1fr;gap:12px}
.form-item{display:flex;flex-direction:column;gap:4px}
.form-item.full-width{grid-column:span 2}
label{font-size:13px;font-weight:500;color:var(--text-2)}
input,textarea{padding:8px 12px;border:1px solid var(--border);border-radius:6px;font-size:14px;font-family:inherit;transition:border-color .2s}
input:focus,textarea:focus{outline:none;border-color:var(--accent)}
textarea{resize:vertical}
</style>