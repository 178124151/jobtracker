<script setup lang="ts">
import AppSidebar from './components/layout/AppSidebar.vue'
import AppTopbar from './components/layout/AppTopbar.vue'
import { ref } from 'vue'

const sidebarCollapsed = ref(false)

const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value
}
</script>

<template>
  <div class="app-shell">
    <AppSidebar :collapsed="sidebarCollapsed" @toggle="toggleSidebar" />
    <main class="app-main">
      <AppTopbar />
      <div class="app-content">
        <router-view />
      </div>
    </main>
  </div>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

:root {
  --bg: #F8FAFC;
  --surface: #FFFFFF;
  --border: #E2E8F0;
  --text-1: #0F172A;
  --text-2: #334155;
  --text-3: #64748B;
  --accent: #2563EB;
  --accent-bg: #EFF6FF;
  --green: #16A34A;
  --amber: #D97706;
  --red: #DC2626;
  --sidebar-w: 240px;
  --sidebar-collapsed-w: 64px;
  --topbar-h: 56px;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', sans-serif;
  background: var(--bg);
  color: var(--text-1);
  height: 100vh;
  overflow: hidden;
}

.app-shell {
  display: flex;
  height: 100vh;
}

.app-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.app-content {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}
</style>