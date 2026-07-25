<script setup lang="ts">
import { useRoute } from 'vue-router'
const route = useRoute()
const title = route.meta?.title as string || 'JobTracker'
</script>

<template>
  <header class="topbar">
    <h1>{{ title }}</h1>
    <div class="actions">
      <button class="btn-icon" title="搜索">🔍</button>
      <button class="btn-icon" title="通知">🔔</button>
    </div>
  </header>
</template>

<style scoped>
.topbar{height:var(--topbar-h);background:var(--surface);border-bottom:1px solid var(--border);display:flex;align-items:center;justify-content:space-between;padding:0 24px}
h1{font-size:18px;font-weight:600}
.actions{display:flex;gap:8px}
.btn-icon{width:36px;height:36px;border:none;background:none;border-radius:8px;cursor:pointer;display:flex;align-items:center;justify-content:center;font-size:16px}
.btn-icon:hover{background:var(--bg)}
</style>