<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed, nextTick } from 'vue'

interface SmeCompany {
  id: string
  name: string
  boss_url: string
  zhilian_url: string
  qiancheng_url: string
  liepin_url: string
}

const companies = ref<SmeCompany[]>([])
const loading = ref(false)
const searchQuery = ref('')
const currentPage = ref(1)
const gridRef = ref<HTMLElement | null>(null)

// 卡片尺寸配置
const CARD_WIDTH = 200
const CARD_HEIGHT = 150
const GAP = 16

// 动态计算每页数量
const pageSize = ref(20)
const containerWidth = ref(0)
const containerHeight = ref(0)

const filteredCompanies = computed(() => {
  if (!searchQuery.value.trim()) {
    return companies.value
  }
  return companies.value.filter(c => 
    c.name.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

const totalPages = computed(() => Math.ceil(filteredCompanies.value.length / pageSize.value))

const paginatedCompanies = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  return filteredCompanies.value.slice(start, start + pageSize.value)
})

const fetchCompanies = async () => {
  loading.value = true
  try {
    const response = await fetch('/api/v1/sme-companies')
    const data = await response.json()
    companies.value = data.data || []
  } catch (error) {
    console.error('Failed to fetch SME companies:', error)
    companies.value = []
  } finally {
    loading.value = false
  }
}

const calculatePageSize = () => {
  // 获取内容区域的宽度和高度
  const contentArea = document.querySelector('.app-content')
  if (!contentArea) return
  
  containerWidth.value = contentArea.clientWidth - 48 // 减去padding
  containerHeight.value = contentArea.clientHeight - 120 // 减去header和padding
  
  // 计算每行卡片数
  const cols = Math.floor((containerWidth.value + GAP) / (CARD_WIDTH + GAP))
  
  // 计算能放多少行
  const rows = Math.floor((containerHeight.value + GAP) / (CARD_HEIGHT + GAP))
  
  // 每页数量 = 列数 × 行数
  pageSize.value = Math.max(cols * rows, 4) // 最少4个
}

const handleResize = () => {
  calculatePageSize()
  // 确保当前页不超过总页数
  if (currentPage.value > totalPages.value) {
    currentPage.value = Math.max(totalPages.value, 1)
  }
}

const handleSearch = () => {
  currentPage.value = 1
}

const openUrl = (url: string) => {
  if (url) window.open(url, '_blank')
}

onMounted(() => {
  fetchCompanies()
  calculatePageSize()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<template>
  <div class="sme-page">
    <div class="page-header">
      <h2>专精特新企业</h2>
      <div class="search-box">
        <input 
          v-model="searchQuery" 
          placeholder="搜索企业名称..." 
          @input="handleSearch"
        />
        <span class="search-count">共 {{ filteredCompanies.length }} 家企业</span>
      </div>
    </div>

    <div v-if="loading" class="loading">加载中...</div>

    <template v-else>
      <div ref="gridRef" class="company-grid">
        <div 
          v-for="company in paginatedCompanies" 
          :key="company.id" 
          class="company-card"
        >
          <div class="card-header">
            <div class="card-logo">{{ company.name[0] }}</div>
            <div class="card-name">{{ company.name }}</div>
          </div>
          <div class="card-buttons">
            <button 
              v-if="company.boss_url" 
              class="btn-boss" 
              @click="openUrl(company.boss_url)"
              title="BOSS直聘"
            >
              BOSS
            </button>
            <button 
              v-if="company.zhilian_url" 
              class="btn-zhilian" 
              @click="openUrl(company.zhilian_url)"
              title="智联招聘"
            >
              智联
            </button>
            <button 
              v-if="company.qiancheng_url" 
              class="btn-qiancheng" 
              @click="openUrl(company.qiancheng_url)"
              title="前程无忧"
            >
              前程
            </button>
            <button 
              v-if="company.liepin_url" 
              class="btn-liepin" 
              @click="openUrl(company.liepin_url)"
              title="猎聘"
            >
              猎聘
            </button>
          </div>
        </div>
      </div>

      <div v-if="totalPages > 1" class="pagination">
        <button 
          :disabled="currentPage === 1" 
          @click="currentPage--"
        >
          上一页
        </button>
        <span class="page-info">{{ currentPage }} / {{ totalPages }}</span>
        <button 
          :disabled="currentPage === totalPages" 
          @click="currentPage++"
        >
          下一页
        </button>
      </div>
    </template>
  </div>
</template>

<style scoped>
.sme-page {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;
  flex-shrink: 0;
}

.page-header h2 {
  font-size: 20px;
  font-weight: 600;
}

.search-box {
  display: flex;
  align-items: center;
  gap: 12px;
}

.search-box input {
  padding: 8px 16px;
  border: 1px solid var(--border);
  border-radius: 8px;
  font-size: 14px;
  width: 250px;
}

.search-box input:focus {
  outline: none;
  border-color: var(--accent);
}

.search-count {
  font-size: 13px;
  color: var(--text-3);
}

.loading {
  text-align: center;
  padding: 40px;
  color: var(--text-3);
}

.company-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
  flex: 1;
  align-content: start;
}

.company-card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 10px;
  padding: 16px;
  transition: all 0.2s;
  height: fit-content;
}

.company-card:hover {
  border-color: var(--accent);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 14px;
}

.card-logo {
  width: 36px;
  height: 36px;
  min-width: 36px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 700;
  font-size: 14px;
}

.card-name {
  font-size: 14px;
  font-weight: 500;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-buttons {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}

.card-buttons button {
  padding: 6px 8px;
  border: none;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  color: white;
}

.btn-boss {
  background: #00b38a;
}

.btn-boss:hover {
  background: #009973;
}

.btn-zhilian {
  background: #1890ff;
}

.btn-zhilian:hover {
  background: #096dd9;
}

.btn-qiancheng {
  background: #ff6600;
}

.btn-qiancheng:hover {
  background: #e65c00;
}

.btn-liepin {
  background: #ff4d4f;
}

.btn-liepin:hover {
  background: #cf1322;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 16px;
  padding: 16px 0;
  flex-shrink: 0;
}

.pagination button {
  padding: 8px 16px;
  border: 1px solid var(--border);
  background: var(--surface);
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
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