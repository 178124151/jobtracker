import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface BasicInfo {
  name: string
  phone: string
  email: string
  city: string
  website: string
  summary: string
  avatar: string
}

export interface Education {
  id: string
  school: string
  major: string
  degree: string
  startDate: string
  endDate: string
  description: string
}

export interface WorkExperience {
  id: string
  company: string
  position: string
  startDate: string
  endDate: string
  description: string
  highlights: string[]
}

export interface Project {
  id: string
  name: string
  role: string
  startDate: string
  endDate: string
  description: string
  highlights: string[]
  link: string
}

export interface Skill {
  id: string
  category: string
  items: string[]
}

export interface ResumeData {
  basic: BasicInfo
  education: Education[]
  work: WorkExperience[]
  projects: Project[]
  skills: Skill[]
}

const defaultResume: ResumeData = {
  basic: { name: '', phone: '', email: '', city: '', website: '', summary: '', avatar: '' },
  education: [],
  work: [],
  projects: [],
  skills: []
}

export const useResumeStore = defineStore('resume', () => {
  const resume = ref<ResumeData>(JSON.parse(JSON.stringify(defaultResume)))
  const currentTemplate = ref('classic')
  const isSaving = ref(false)
  const lastSaved = ref<string | null>(null)
  const resumeId = ref<string | null>(null)

  const genId = () => Math.random().toString(36).substring(2, 11)

  const updateBasic = (data: Partial<BasicInfo>) => Object.assign(resume.value.basic, data)

  const addEducation = () => resume.value.education.push({ id: genId(), school: '', major: '', degree: '', startDate: '', endDate: '', description: '' })
  const removeEducation = (id: string) => { resume.value.education = resume.value.education.filter(e => e.id !== id) }
  const updateEducation = (id: string, data: Partial<Education>) => { const i = resume.value.education.find(e => e.id === id); if (i) Object.assign(i, data) }

  const addWork = () => resume.value.work.push({ id: genId(), company: '', position: '', startDate: '', endDate: '', description: '', highlights: [''] })
  const removeWork = (id: string) => { resume.value.work = resume.value.work.filter(w => w.id !== id) }
  const updateWork = (id: string, data: Partial<WorkExperience>) => { const i = resume.value.work.find(w => w.id === id); if (i) Object.assign(i, data) }

  const addProject = () => resume.value.projects.push({ id: genId(), name: '', role: '', startDate: '', endDate: '', description: '', highlights: [''], link: '' })
  const removeProject = (id: string) => { resume.value.projects = resume.value.projects.filter(p => p.id !== id) }
  const updateProject = (id: string, data: Partial<Project>) => { const i = resume.value.projects.find(p => p.id === id); if (i) Object.assign(i, data) }

  const addSkill = () => resume.value.skills.push({ id: genId(), category: '', items: [] })
  const removeSkill = (id: string) => { resume.value.skills = resume.value.skills.filter(s => s.id !== id) }
  const updateSkill = (id: string, data: Partial<Skill>) => { const i = resume.value.skills.find(s => s.id === id); if (i) Object.assign(i, data) }

  const setTemplate = (t: string) => { currentTemplate.value = t }

  const saveResume = async () => {
    isSaving.value = true
    try {
      const res = await fetch('/api/v1/resumes', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id: resumeId.value, title: `${resume.value.basic.name || '未命名'}的简历`, content: JSON.stringify(resume.value), template: currentTemplate.value })
      })
      const data = await res.json()
      if (data.data?.id) resumeId.value = data.data.id
      lastSaved.value = new Date().toLocaleTimeString()
    } catch (e) { console.error('保存失败:', e) }
    finally { isSaving.value = false }
  }

  const loadResume = async (id: string) => {
    try {
      const res = await fetch(`/api/v1/resumes/${id}`)
      const data = await res.json()
      if (data.data?.content) {
        resume.value = JSON.parse(data.data.content)
        resumeId.value = id
        if (data.data.template) currentTemplate.value = data.data.template
      }
    } catch (e) { console.error('加载失败:', e) }
  }

  const resetResume = () => { resume.value = JSON.parse(JSON.stringify(defaultResume)); resumeId.value = null }

  return {
    resume, currentTemplate, isSaving, lastSaved, resumeId,
    updateBasic, addEducation, removeEducation, updateEducation,
    addWork, removeWork, updateWork, addProject, removeProject, updateProject,
    addSkill, removeSkill, updateSkill, setTemplate, saveResume, loadResume, resetResume
  }
})