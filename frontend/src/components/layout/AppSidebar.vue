<script setup lang="ts">
import { useRouter, useRoute } from 'vue-router'

defineProps<{ collapsed: boolean }>()
const emit = defineEmits(['toggle'])
const router = useRouter()
const route = useRoute()

const navItems = [
  { path: '/companies', label: '公司导航', icon: '司' },
  { path: '/sme', label: '小而美企业', icon: '美' },
  { path: '/resume', label: '简历制作', icon: '简' },
  { path: '/monitor', label: '链路监控', icon: '监' }
]
</script>

<template>
  <aside class="sidebar" :class="{ collapsed }">
    <div class="sidebar-header" @click="emit('toggle')">
      <div class="logo">
        <span class="logo-mark">JT</span>
        <span v-if="!collapsed" class="logo-text">JobTracker</span>
      </div>
    </div>
    <nav class="sidebar-nav">
      <div v-for="item in navItems" :key="item.path" class="nav-item" :class="{ active: route.path === item.path }" @click="router.push(item.path)" :title="collapsed ? item.label : ''">
        <span class="nav-icon">{{ item.icon }}</span>
        <span v-if="!collapsed" class="nav-label">{{ item.label }}</span>
      </div>
    </nav>
    <div class="sidebar-footer">
      <div class="user-info">
        <div class="avatar">L</div>
        <span v-if="!collapsed" class="username">Luci</span>
      </div>
    </div>
  </aside>
</template>

<style scoped>
.sidebar{width:var(--sidebar-w);background:var(--surface);border-right:1px solid var(--border);display:flex;flex-direction:column;transition:width .3s;overflow:hidden}
.sidebar.collapsed{width:var(--sidebar-collapsed-w)}
.sidebar-header{height:var(--topbar-h);display:flex;align-items:center;justify-content:center;padding:0 12px;border-bottom:1px solid var(--border);cursor:pointer;transition:background .2s}
.sidebar-header:hover{background:var(--bg)}
.logo{display:flex;align-items:center;gap:10px}
.logo-mark{width:32px;height:32px;min-width:32px;background:var(--accent);border-radius:6px;display:flex;align-items:center;justify-content:center;color:#fff;font-weight:700;font-size:14px}
.logo-text{font-weight:700;font-size:16px;white-space:nowrap}
.sidebar-nav{flex:1;padding:12px 8px;overflow-y:auto}
.nav-item{display:flex;align-items:center;gap:12px;padding:10px 12px;border-radius:8px;cursor:pointer;margin-bottom:4px;transition:all .2s;white-space:nowrap}
.nav-item:hover{background:var(--bg)}
.nav-item.active{background:var(--accent-bg)}
.nav-item.active .nav-icon{background:var(--accent);color:#fff}
.nav-item.active .nav-label{color:var(--accent);font-weight:500}
.nav-icon{width:24px;height:24px;min-width:24px;display:flex;align-items:center;justify-content:center;font-size:12px;font-weight:600;background:var(--accent);color:#fff;border-radius:4px;flex-shrink:0}
.nav-label{font-size:14px;overflow:hidden}
.sidebar-footer{padding:16px;border-top:1px solid var(--border)}
.user-info{display:flex;align-items:center;gap:10px;justify-content:center}
.avatar{width:32px;height:32px;min-width:32px;background:var(--accent);border-radius:50%;display:flex;align-items:center;justify-content:center;color:#fff;font-weight:600}
.username{font-size:14px;font-weight:500;white-space:nowrap}
</style>