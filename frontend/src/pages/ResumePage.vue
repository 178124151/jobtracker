<script setup lang="ts">
import { ref } from 'vue'
import { useResumeStore } from '@/stores/resume'
import BasicEditor from '@/components/resume/BasicEditor.vue'
import EducationEditor from '@/components/resume/EducationEditor.vue'
import WorkEditor from '@/components/resume/WorkEditor.vue'
import ProjectEditor from '@/components/resume/ProjectEditor.vue'
import SkillsEditor from '@/components/resume/SkillsEditor.vue'
import ResumePreview from '@/components/resume/ResumePreview.vue'

const store = useResumeStore()
const activeTab = ref('basic')
const showPreview = ref(true)
const isExporting = ref(false)

const tabs = [
  { key: 'basic', label: '基本信息', icon: '👤' },
  { key: 'education', label: '教育背景', icon: '🎓' },
  { key: 'work', label: '工作经历', icon: '💼' },
  { key: 'projects', label: '项目经历', icon: '🚀' },
  { key: 'skills', label: '专业技能', icon: '⚡' }
]

const handleSave = async () => { await store.saveResume() }

const handleExportJSON = () => {
  const blob = new Blob([JSON.stringify(store.resume, null, 2)], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a'); a.href = url; a.download = `简历-${store.resume.basic.name || '未命名'}.json`; a.click(); URL.revokeObjectURL(url)
}

const handleImportJSON = () => {
  const input = document.createElement('input'); input.type = 'file'; input.accept = '.json'
  input.onchange = (e) => {
    const file = (e.target as HTMLInputElement).files?.[0]
    if (!file) return
    const reader = new FileReader()
    reader.onload = (ev) => { try { store.$patch({ resume: JSON.parse(ev.target?.result as string) }) } catch { alert('JSON 文件格式错误') } }
    reader.readAsText(file)
  }
  input.click()
}

const handleExportPDF = async () => {
  isExporting.value = true
  try {
    const html2canvas = (await import('html2canvas')).default
    const jsPDF = (await import('jspdf')).default

    const el = document.querySelector('.preview-paper .resume-preview') as HTMLElement
    if (!el) { alert('未找到简历预览内容'); return }

    const canvas = await html2canvas(el, { scale: 2, useCORS: true, backgroundColor: '#ffffff' })
    const imgData = canvas.toDataURL('image/png')

    const pdf = new jsPDF('p', 'mm', 'a4')
    const pdfWidth = pdf.internal.pageSize.getWidth()
    const pdfHeight = pdf.internal.pageSize.getHeight()

    const imgWidth = canvas.width
    const imgHeight = canvas.height
    const ratio = pdfWidth / imgWidth
    const scaledHeight = imgHeight * ratio

    let position = 0
    let remainingHeight = scaledHeight

    while (remainingHeight > 0) {
      if (position > 0) pdf.addPage()
      pdf.addImage(imgData, 'PNG', 0, -(position * pdfWidth / imgWidth), pdfWidth, scaledHeight)
      position += pdfHeight
      remainingHeight -= pdfHeight
    }

    pdf.save(`简历-${store.resume.basic.name || '未命名'}.pdf`)
  } catch (err) {
    console.error('PDF 导出失败:', err)
    alert('PDF 导出失败，请确保已安装 html2canvas 和 jspdf')
  } finally {
    isExporting.value = false
  }
}

const togglePreview = () => { showPreview.value = !showPreview.value }
</script>

<template>
  <div class="resume-page">
    <div class="toolbar">
      <div class="toolbar-left">
        <h2>📝 简历制作</h2>
        <span v-if="store.lastSaved" class="save-ok">已保存 {{ store.lastSaved }}</span>
      </div>
      <div class="toolbar-right">
        <button class="btn-tool" @click="togglePreview">{{ showPreview ? '隐藏' : '显示' }}预览</button>
        <button class="btn-tool" @click="handleImportJSON">导入 JSON</button>
        <button class="btn-tool" @click="handleExportJSON">导出 JSON</button>
        <button class="btn-export" @click="handleExportPDF" :disabled="isExporting">{{ isExporting ? '导出中...' : '📄 导出 PDF' }}</button>
        <button class="btn-save" @click="handleSave" :disabled="store.isSaving">{{ store.isSaving ? '保存中...' : '💾 保存' }}</button>
      </div>
    </div>

    <div class="content" :class="{ 'no-preview': !showPreview }">
      <div class="editor-panel">
        <div class="tab-bar">
          <button v-for="tab in tabs" :key="tab.key" class="tab-btn" :class="{ active: activeTab === tab.key }" @click="activeTab = tab.key">
            <span class="tab-icon">{{ tab.icon }}</span>
            <span class="tab-label">{{ tab.label }}</span>
          </button>
        </div>
        <div class="editor-content">
          <KeepAlive>
            <BasicEditor v-if="activeTab === 'basic'" />
            <EducationEditor v-else-if="activeTab === 'education'" />
            <WorkEditor v-else-if="activeTab === 'work'" />
            <ProjectEditor v-else-if="activeTab === 'projects'" />
            <SkillsEditor v-else-if="activeTab === 'skills'" />
          </KeepAlive>
        </div>
      </div>

      <div v-if="showPreview" class="preview-panel">
        <div class="preview-container">
          <div class="preview-paper"><ResumePreview /></div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.resume-page{display:flex;flex-direction:column;height:100%}
.toolbar{display:flex;justify-content:space-between;align-items:center;padding:12px 20px;background:var(--surface);border-bottom:1px solid var(--border);flex-shrink:0}
.toolbar-left{display:flex;align-items:center;gap:16px}
.toolbar-left h2{font-size:18px;font-weight:600}
.save-ok{font-size:12px;color:#16a34a;background:#f0fdf4;padding:2px 8px;border-radius:10px}
.toolbar-right{display:flex;gap:8px}
.btn-tool{padding:6px 14px;background:var(--bg);border:1px solid var(--border);border-radius:6px;cursor:pointer;font-size:13px;transition:all .2s}
.btn-tool:hover{border-color:var(--accent);color:var(--accent)}
.btn-export{padding:6px 16px;background:#16a34a;color:#fff;border:none;border-radius:6px;cursor:pointer;font-size:13px;font-weight:500;transition:all .2s}
.btn-export:hover{opacity:.9}.btn-export:disabled{opacity:.6;cursor:not-allowed}
.btn-save{padding:6px 20px;background:var(--accent);color:#fff;border:none;border-radius:6px;cursor:pointer;font-size:13px;font-weight:500;transition:all .2s}
.btn-save:hover{opacity:.9}.btn-save:disabled{opacity:.6;cursor:not-allowed}
.content{display:flex;flex:1;overflow:hidden}
.content.no-preview .editor-panel{width:100%}
.editor-panel{width:50%;min-width:380px;display:flex;flex-direction:column;border-right:1px solid var(--border)}
.tab-bar{display:flex;background:var(--bg);border-bottom:1px solid var(--border);overflow-x:auto;flex-shrink:0}
.tab-btn{display:flex;align-items:center;gap:5px;padding:10px 14px;background:none;border:none;border-bottom:2px solid transparent;cursor:pointer;font-size:13px;color:var(--text-2);white-space:nowrap;transition:all .2s}
.tab-btn:hover{background:var(--surface);color:var(--text-1)}
.tab-btn.active{background:var(--surface);color:var(--accent);border-bottom-color:var(--accent);font-weight:500}
.tab-icon{font-size:15px}
.editor-content{flex:1;overflow-y:auto;padding:20px}
.preview-panel{flex:1;overflow:hidden;background:#d1d5db}
.preview-container{height:100%;overflow-y:auto;padding:24px;display:flex;justify-content:center}
.preview-paper{width:210mm;min-height:297mm;background:#fff;box-shadow:0 4px 20px rgba(0,0,0,.2);border-radius:2px;overflow:hidden}
@media(max-width:1100px){.editor-panel{width:100%;min-width:unset;border-right:none}.content{flex-direction:column}.preview-panel{border-top:1px solid var(--border)}}
</style>