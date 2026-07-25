<script setup lang="ts">
import { useResumeStore } from '@/stores/resume'
const store = useResumeStore()
const hasInfo = () => { const b = store.resume.basic; return b.name || b.email || b.phone }
</script>

<template>
  <div class="resume-preview">
    <header v-if="hasInfo()" class="header">
      <div class="header-main">
        <div v-if="store.resume.basic.avatar" class="avatar"><img :src="store.resume.basic.avatar" alt="证件照" /></div>
        <div class="info">
          <h1>{{ store.resume.basic.name || '您的姓名' }}</h1>
          <div class="contact">
            <span v-if="store.resume.basic.phone">📱 {{ store.resume.basic.phone }}</span>
            <span v-if="store.resume.basic.email">✉️ {{ store.resume.basic.email }}</span>
            <span v-if="store.resume.basic.city">📍 {{ store.resume.basic.city }}</span>
            <a v-if="store.resume.basic.website" :href="store.resume.basic.website" target="_blank">🔗 个人主页</a>
          </div>
        </div>
      </div>
      <p v-if="store.resume.basic.summary" class="summary">{{ store.resume.basic.summary }}</p>
    </header>

    <section v-if="store.resume.education.length > 0" class="section">
      <h2 class="section-title">教育背景</h2>
      <div v-for="e in store.resume.education" :key="e.id" class="item">
        <div class="item-head"><div class="item-left"><strong>{{ e.school }}</strong><span class="sub">{{ e.degree }}{{ e.major ? ' · ' + e.major : '' }}</span></div><span class="date">{{ e.startDate }} - {{ e.endDate }}</span></div>
        <p v-if="e.description" class="desc">{{ e.description }}</p>
      </div>
    </section>

    <section v-if="store.resume.work.length > 0" class="section">
      <h2 class="section-title">工作经历</h2>
      <div v-for="w in store.resume.work" :key="w.id" class="item">
        <div class="item-head"><div class="item-left"><strong>{{ w.company }}</strong><span class="sub">{{ w.position }}</span></div><span class="date">{{ w.startDate }} - {{ w.endDate }}</span></div>
        <p v-if="w.description" class="desc">{{ w.description }}</p>
        <ul v-if="w.highlights.filter(h=>h).length" class="hl"><li v-for="(h,i) in w.highlights.filter(h=>h)" :key="i">{{ h }}</li></ul>
      </div>
    </section>

    <section v-if="store.resume.projects.length > 0" class="section">
      <h2 class="section-title">项目经历</h2>
      <div v-for="p in store.resume.projects" :key="p.id" class="item">
        <div class="item-head"><div class="item-left"><strong>{{ p.name }}<a v-if="p.link" :href="p.link" target="_blank" class="link">[链接]</a></strong><span class="sub">{{ p.role }}</span></div><span class="date">{{ p.startDate }} - {{ p.endDate }}</span></div>
        <p v-if="p.description" class="desc">{{ p.description }}</p>
        <ul v-if="p.highlights.filter(h=>h).length" class="hl"><li v-for="(h,i) in p.highlights.filter(h=>h)" :key="i">{{ h }}</li></ul>
      </div>
    </section>

    <section v-if="store.resume.skills.length > 0" class="section">
      <h2 class="section-title">专业技能</h2>
      <div v-for="s in store.resume.skills" :key="s.id" class="skill"><strong>{{ s.category }}：</strong><span>{{ s.items.join('、') }}</span></div>
    </section>

    <div v-if="!hasInfo() && !store.resume.education.length && !store.resume.work.length && !store.resume.projects.length && !store.resume.skills.length" class="empty">
      <div class="empty-icon">📝</div>
      <p>在左侧面板填写信息</p>
      <p>简历将在此实时预览</p>
    </div>
  </div>
</template>

<style scoped>
.resume-preview{background:#fff;color:#1a1a1a;padding:36px 40px;font-family:'Noto Sans SC','PingFang SC','Microsoft YaHei',sans-serif;line-height:1.6;min-height:100%;font-size:13px}
.header{padding-bottom:16px;border-bottom:2px solid #2563eb;margin-bottom:18px}
.header-main{display:flex;gap:20px;align-items:center}
.avatar{width:90px;height:120px;border-radius:6px;overflow:hidden;border:1px solid #e5e7eb;flex-shrink:0}
.avatar img{width:100%;height:100%;object-fit:cover}
.info{flex:1}
.info h1{font-size:24px;font-weight:700;margin-bottom:6px}
.contact{display:flex;flex-wrap:wrap;gap:12px;font-size:12px;color:#555}
.contact a{color:#2563eb;text-decoration:none}
.summary{margin-top:10px;font-size:13px;color:#444;font-style:italic;border-left:3px solid #2563eb;padding-left:10px}
.section{margin-bottom:16px}
.section-title{font-size:14px;font-weight:700;color:#2563eb;text-transform:uppercase;letter-spacing:1.5px;border-bottom:1px solid #e5e7eb;padding-bottom:4px;margin-bottom:10px}
.item{margin-bottom:12px}
.item-head{display:flex;justify-content:space-between;align-items:flex-start}
.item-left{display:flex;flex-direction:column}
.item-left strong{font-size:14px}
.sub{font-size:12px;color:#555;font-style:italic}
.date{font-size:12px;color:#666;white-space:nowrap}
.desc{font-size:12px;color:#444;margin-top:3px}
.hl{margin:4px 0 0 18px;padding:0}
.hl li{font-size:12px;color:#444;margin-bottom:2px}
.skill{margin-bottom:4px;font-size:13px}
.skill strong{color:#1a1a1a}
.link{color:#2563eb;text-decoration:none;font-size:11px;font-weight:normal;margin-left:4px}
.link:hover{text-decoration:underline}
.empty{display:flex;flex-direction:column;align-items:center;justify-content:center;min-height:400px;color:#999;text-align:center}
.empty-icon{font-size:56px;margin-bottom:12px}
.empty p{margin-bottom:6px;font-size:14px}
</style>