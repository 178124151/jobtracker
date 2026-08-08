<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'

const grafanaBaseUrl = computed(() => window.location.origin)

const dashboardUid = 'server-monitoring'

const panels = ref([
  { id: 1, title: 'CPU Usage', type: 'gauge', height: '300px' },
  { id: 2, title: 'Memory Usage', type: 'gauge', height: '300px' },
  { id: 11, title: 'CPU Trend', type: 'timeseries', height: '400px' }
])

const getGrafanaUrl = (panelId: number) => {
  return `${grafanaBaseUrl.value}/grafana/d-solo/${dashboardUid}/server-monitoring?orgId=1&panelId=${panelId}&theme=light&kiosk`
}

const metrics = ref({
  uptime: '0 days',
  totalRequests: '0',
  avgResponseTime: '0ms',
  errorRate: '0%'
})

const recentAlerts = ref([
  { id: 1, time: 'No alerts', level: 'info', message: 'System running normally' }
])

const fetchMetrics = async () => {
  try {
    const response = await fetch('/api/v1/sre/health')
    const data = await response.json()
    if (data.status === 'ok') {
      metrics.value.uptime = 'Running'
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
      <h2>Link Monitoring</h2>
      <div class="header-links">
        <a :href="`${grafanaBaseUrl}/grafana/`" target="_blank" class="grafana-link">Open Grafana Panel</a>
      </div>
    </div>

    <div class="metrics-grid">
      <div class="metric-card">
        <div class="metric-label">Service Status</div>
        <div class="metric-value green">{{ metrics.uptime }}</div>
      </div>
      <div class="metric-card">
        <div class="metric-label">Total Requests</div>
        <div class="metric-value">{{ metrics.totalRequests }}</div>
      </div>
      <div class="metric-card">
        <div class="metric-label">Avg Response</div>
        <div class="metric-value">{{ metrics.avgResponseTime }}</div>
      </div>
      <div class="metric-card">
        <div class="metric-label">Error Rate</div>
        <div class="metric-value">{{ metrics.errorRate }}</div>
      </div>
    </div>

    <div class="grafana-section">
      <h3>Server Monitoring</h3>
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
              :href="`${grafanaBaseUrl}/grafana/d/${dashboardUid}?orgId=1&viewPanel=${panel.id}`"
              target="_blank"
              class="panel-expand"
            >&#8599;</a>
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
      <h3>Recent Alerts</h3>
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
