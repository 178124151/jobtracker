<script setup lang="ts">
import { onMounted, onUnmounted, ref, computed } from 'vue'
import { useCompanyStore } from '@/stores/company'

const store = useCompanyStore()

// 分页相关
const currentPage = ref(1)
const pageSize = ref(20)

// 卡片尺寸配置
const CARD_WIDTH = 250
const CARD_HEIGHT = 200
const GAP = 16

const groups = [
  { key: 'all', label: 'All' },
  { key: 'bigtech', label: 'Big Tech' },
  { key: 'cloud', label: 'Cloud' },
  { key: 'foreign', label: 'Foreign' },
  { key: 'startup', label: 'Startup' }
]

// 计算每页数量
const calculatePageSize = () => {
  const contentArea = document.querySelector('.app-content')
  if (!contentArea) return
  
  const containerWidth = contentArea.clientWidth - 48
  const containerHeight = contentArea.clientHeight - 150
  
  const cols = Math.floor((containerWidth + GAP) / (CARD_WIDTH + GAP))
  const rows = Math.floor((containerHeight + GAP) / (CARD_HEIGHT + GAP))
  
  pageSize.value = Math.max(cols * rows, 4)
}

const handleResize = () => {
  calculatePageSize()
  if (currentPage.value > totalPages.value) {
    currentPage.value = Math.max(totalPages.value, 1)
  }
}

const totalPages = computed(() => Math.ceil(store.filteredCompanies.length / pageSize.value))

const paginatedCompanies = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  return store.filteredCompanies.slice(start, start + pageSize.value)
})

const setGroup = (group: string) => {
  store.setGroup(group)
  currentPage.value = 1
}

const openWebsite = (url: string) => {
  window.open(url, '_blank')
}

const getStatusColor = (status: string) => {
  switch (status) {
    case 'GREEN': return '#16A34A'
    case 'YELLOW': return '#D97706'
    case 'RED': return '#DC2626'
    default: return '#9CA3AF'
  }
}

const getStatusText = (status: string) => {
  switch (status) {
    case 'GREEN': return 'Online'
    case 'YELLOW': return 'Checking'
    case 'RED': return 'Offline'
    default: return 'Unknown'
  }
}

onMounted(() => {
  store.fetchCompanies()
  calculatePageSize()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<template>
  <div class="company-page">
    <div class="page-header">
      <div class="filter-bar">
        <button
          v-for="group in groups"
          :key="group.key"
          class="filter-btn"
          :class="{ active: store.activeGroup === group.key }"
          @click="setGroup(group.key)"
        >
          {{ group.label }}
        </button>
      </div>
      <span class="company-count">共 {{ store.filteredCompanies.length }} 家公司</span>
    </div>

    <div class="company-grid">
      <div
        v-for="company in paginatedCompanies"
        :key="company.id"
        class="company-card"
      >
        <div class="card-header">
          <div class="card-logo" :style="{ background: getStatusColor(company.healthStatus) }">
            {{ company.name[0] }}
          </div>
          <div class="card-info">
            <h3 class="card-name">{{ company.name }}</h3>
            <div class="card-status">
              <span class="status-dot" :style="{ background: getStatusColor(company.healthStatus) }"></span>
              <span class="status-text">{{ getStatusText(company.healthStatus) }}</span>
            </div>
          </div>
        </div>
        <p class="card-desc">{{ company.description || '暂无简介' }}</p>
        <div class="card-footer">
          <button class="card-btn" @click="openWebsite(company.website)">
            Visit Website
          </button>
        </div>
      </div>
    </div>

    <div v-if="totalPages > 1" class="pagination">
      <button 
        :disabled="currentPage === 1" 
        @click="currentPage--"
      >
        Previous
      </button>
      <span class="page-info">{{ currentPage }} / {{ totalPages }}</span>
      <button 
        :disabled="currentPage === totalPages" 
        @click="currentPage++"
      >
        Next
      </button>
    </div>
  </div>
</template>

<style scoped>
.company-page {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  flex-wrap: wrap;
  gap: 16px;
  flex-shrink: 0;
}

.filter-bar {
  display: flex;
  gap: 8px;
}

.filter-btn {
  padding: 8px 16px;
  border: 1px solid var(--border);
  background: var(--surface);
  border-radius: 20px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

.filter-btn:hover {
  border-color: var(--accent);
}

.filter-btn.active {
  background: var(--accent);
  color: white;
  border-color: var(--accent);
}

.company-count {
  font-size: 13px;
  color: var(--text-3);
}

.company-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 16px;
  flex: 1;
  align-content: start;
}

.company-card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 20px;
  transition: all 0.2s;
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 180px;
}

.company-card:hover {
  border-color: var(--accent);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.card-logo {
  width: 40px;
  height: 40px;
  min-width: 40px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 700;
  font-size: 18px;
}

.card-info {
  flex: 1;
  min-width: 0;
}

.card-name {
  font-size: 15px;
  font-weight: 600;
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.card-status {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--text-3);
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.card-desc {
  font-size: 13px;
  color: var(--text-2);
  line-height: 1.5;
  flex: 1;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
  margin-bottom: 12px;
}

.card-footer {
  margin-top: auto;
}

.card-btn {
  width: 100%;
  padding: 10px;
  background: var(--accent-bg);
  color: var(--accent);
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s;
}

.card-btn:hover {
  background: var(--accent);
  color: white;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 16px;
  padding: 16px 0 0;
  flex-shrink: 0;
}

.pagination button {
  padding: 8px 16px;
  border: 1px solid var(--border);
  background: var(--surface);
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

.pagination button:hover:not(:disabled) {
  border-color: var(--accent);
  color: var(--accent);
}

.pagination button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.page-info {
  font-size: 14px;
  color: var(--text-2);
}
</style>