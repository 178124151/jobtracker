<script setup lang="ts">
import { ref } from 'vue'

interface Application {
  id: string
  company: string
  position: string
  appliedAt: string
  status: string
  notes: string
}

const applications = ref<Application[]>([
  { id: '1', company: '阿里巴巴', position: 'SRE工程师', appliedAt: '2026-07-01', status: 'INTERVIEW', notes: '一面通过' },
  { id: '2', company: '腾讯', position: '运维开发', appliedAt: '2026-07-05', status: 'RESUME', notes: '' },
  { id: '3', company: '字节跳动', position: 'SRE', appliedAt: '2026-07-10', status: 'WRITTEN', notes: '笔试中' }
])

const statusOptions = [
  { value: 'RESUME', label: '简历筛选', color: '#2563EB' },
  { value: 'WRITTEN', label: '笔试', color: '#9333EA' },
  { value: 'INTERVIEW', label: '面试中', color: '#D97706' },
  { value: 'OFFER', label: 'Offer', color: '#16A34A' },
  { value: 'REJECTED', label: '已拒', color: '#DC2626' },
  { value: 'WITHDRAWN', label: '已放弃', color: '#6B7280' }
]

const getStatusLabel = (status: string) => {
  return statusOptions.find(s => s.value === status)?.label || status
}

const getStatusColor = (status: string) => {
  return statusOptions.find(s => s.value === status)?.color || '#6B7280'
}
</script>

<template>
  <div class="tracker-page">
    <div class="tracker-header">
      <h2>投递进度管理</h2>
      <button class="btn-primary">+ 新增投递</button>
    </div>

    <div class="tracker-table">
      <table>
        <thead>
          <tr>
            <th>公司</th>
            <th>岗位</th>
            <th>投递日期</th>
            <th>状态</th>
            <th>备注</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="app in applications" :key="app.id">
            <td>{{ app.company }}</td>
            <td>{{ app.position }}</td>
            <td>{{ app.appliedAt }}</td>
            <td>
              <span class="status-badge" :style="{ background: getStatusColor(app.status) + '20', color: getStatusColor(app.status) }">
                {{ getStatusLabel(app.status) }}
              </span>
            </td>
            <td>{{ app.notes }}</td>
            <td>
              <button class="btn-sm">编辑</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped>
.tracker-page {
  width: 100%;
}

.tracker-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.btn-primary {
  padding: 10px 20px;
  background: var(--accent);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
}

.btn-primary:hover {
  opacity: 0.9;
}

.tracker-table {
  background: var(--surface);
  border-radius: 12px;
  border: 1px solid var(--border);
  overflow: hidden;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th, td {
  padding: 14px 16px;
  text-align: left;
  border-bottom: 1px solid var(--border);
}

th {
  background: var(--bg);
  font-weight: 600;
  font-size: 13px;
  color: var(--text-3);
  text-transform: uppercase;
}

td {
  font-size: 14px;
}

.status-badge {
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
}

.btn-sm {
  padding: 6px 12px;
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: 6px;
  cursor: pointer;
  font-size: 12px;
}

.btn-sm:hover {
  background: var(--border);
}
</style>
