import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/companies' },
    { path: '/companies', name: 'companies', component: () => import('@/pages/CompanyPage.vue'), meta: { title: '公司导航' } },
    { path: '/sme', name: 'sme', component: () => import('@/pages/SmallBusinessPage.vue'), meta: { title: '小而美企业' } },
    { path: '/tracker', name: 'tracker', component: () => import('@/pages/TrackerPage.vue'), meta: { title: '投递进度' } },
    { path: '/resume', name: 'resume', component: () => import('@/pages/ResumePage.vue'), meta: { title: '简历制作' } },
    { path: '/monitor', name: 'monitor', component: () => import('@/pages/MonitorPage.vue'), meta: { title: '链路监控' } }
  ]
})

export default router