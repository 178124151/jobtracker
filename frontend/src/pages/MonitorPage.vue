<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'

// Grafana 通过 NodePort 访问
const grafanaBaseUrl = computed(() => {
  const host = window.location.hostname
  return `http://${host}:30300`
})

const dashboardUid = 'server-monitoring'

const panels = ref([
  { id: 1, title: 'CPU 使用率', type: 'gauge', height: '300px' },
  { id: 2, title: '内存使用率', type: 'gauge', height: '300px' },
  { id: 3, title: 'CPU 使用趋势', type: 'timeseries', height: '400px' }
])

const getGrafanaUrl = (panelId: number) => {
  return `${grafanaBaseUrl.value}/d-solo/${dashboardUid}/server-monitoring?orgId=1&panelId=${panelId}&theme=light&kiosk`
}

const metrics = ref({
  uptime: '0天',
  totalRequests: '0',
  avgResponseTime: '0ms',
  errorRate: '0%'
})

const recentAlerts = ref([
  { id: 1, time: '暂无告警', level: 'info', message: '系统运行正常' }
])

const fetchMetrics = async () => {
  try {
    const response = await fetch('/api/v1/sre/health')
    const data = await response.json()
    if (data.status === 'ok') {
      metrics.value.uptime = '运行中'
    }
  } catch (error) {
    console.error('Failed to fetch metrics:', error)
  }
}

onMounted(() => {
  fetchMetrics()
})
</script>

<template>
  <div class="monitor-page">
    <div class="page-header">
      <h2>链路监控</h2>
      <div class="header-links">
        <a :href="grafanaBaseUrl" target="_blank" class="grafana-link">打开 Grafana 面板</a>
      </div>
    </div>

    <div class="metrics-grid">
      <div class="metric-card">
        <div class="metric-label">服务状态</div>
        <div class="metric-value green">{{ metrics.uptime }}</div>
      </div>
      <div class="metric-card">
        <div class="metric-label">总请求数</div>
        <div class="metric-value">{{ metrics.totalRequests }}</div>
      </div>
      <div class="metric-card">
        <div class="metric-label">平均响应</div>
        <div class="metric-value">{{ metrics.avgResponseTime }}</div>
      </div>
      <div class="metric-card">
        <div class="metric-label">错误率</div>
        <div class="metric-value">{{ metrics.errorRate }}</div>
      </div>
    </div>

    <div class="grafana-section">
      <h3>服务器监控</h3>
      <div class="grafana-panels">
        <div
          v-for="panel in panels"
          :key="panel.id"
          class="grafana-panel"
          :style="{ height: panel.height }"
        >
          <div class="panel-header">
            <span class="panel-title">{{ panel.title }}</span>
            <a
              :href="`${grafanaBaseUrl}/d/${dashboardUid}?orgId=1&viewPanel=${panel.id}`"
              target="_blank"
              class="panel-expand"
            >↗</a>
          </div>
          <iframe
            :src="getGrafanaUrl(panel.id)"
            frameborder="0"
            allowfullscreen
            class="grafana-iframe"
          ></iframe>
        </div>
      </div>
    </div>

    <div class="alerts-section">
      <h3>最近告警</h3>
      <div class="alerts-list">
        <div
          v-for="alert in recentAlerts"
          :key="alert.id"
          class="alert-item"
          :class="alert.level"
        >
          <span class="alert-time">{{ alert.time }}</span>
          <span class="alert-level">{{ alert.level }}</span>
          <span class="alert-message">{{ alert.message }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.monitor-page { width: 100%; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.page-header h2 { font-size: 20px; font-weight: 600; }
.header-links { display: flex; gap: 8px; }
.grafana-link { padding: 8px 16px; background: var(--accent); color: white; border-radius: 8px; text-decoration: none; font-size: 14px; }
.grafana-link:hover { opacity: 0.9; }
.metrics-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 16px; margin-bottom: 24px; }
.metric-card { background: var(--surface); border: 1px solid var(--border); border-radius: 10px; padding: 20px; text-align: center; }
.metric-label { font-size: 13px; color: var(--text-3); margin-bottom: 8px; }
.metric-value { font-size: 24px; font-weight: 700; color: var(--text-1); }
.metric-value.green { color: #16a34a; }
.grafana-section { margin-bottom: 24px; }
.grafana-section h3 { font-size: 16px; font-weight: 600; margin-bottom: 16px; }
.grafana-panels { display: grid; grid-template-columns: repeat(2, 1fr); gap: 16px; }
.grafana-panel { background: var(--surface); border: 1px solid var(--border); border-radius: 10px; overflow: hidden; display: flex; flex-direction: column; }
.grafana-panel:last-child { grid-column: span 2; }
.panel-header { display: flex; justify-content: space-between; align-items: center; padding: 12px 16px; border-bottom: 1px solid var(--border); background: var(--bg); }
.panel-title { font-size: 14px; font-weight: 500; }
.panel-expand { color: var(--text-3); text-decoration: none; font-size: 16px; }
.panel-expand:hover { color: var(--accent); }
.grafana-iframe { flex: 1; width: 100%; min-height: 200px; }
.alerts-section { background: var(--surface); border: 1px solid var(--border); border-radius: 10px; padding: 20px; }
.alerts-section h3 { font-size: 16px; font-weight: 600; margin-bottom: 16px; }
.alerts-list { display: flex; flex-direction: column; gap: 8px; }
.alert-item { display: flex; align-items: center; gap: 12px; padding: 12px; border-radius: 8px; font-size: 14px; }
.alert-item.info { background: #eff6ff; color: #2563eb; }
.alert-item.warning { background: #fef3c7; color: #d97706; }
.alert-item.error { background: #fee2e2; color: #dc2626; }
.alert-time { font-size: 12px; opacity: 0.7; min-width: 100px; }
.alert-level { font-weight: 600; text-transform: uppercase; font-size: 12px; min-width: 60px; }
.alert-message { flex: 1; }
</style>