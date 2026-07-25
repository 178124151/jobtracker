import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface Company {
  id: string
  name: string
  website: string
  industry: string
  group: string
  description: string
  healthStatus: 'GREEN' | 'YELLOW' | 'RED' | 'UNKNOWN'
}

export const useCompanyStore = defineStore('company', () => {
  const companies = ref<Company[]>([])
  const loading = ref(false)
  const activeGroup = ref('all')

  const filteredCompanies = computed(() => {
    if (activeGroup.value === 'all') {
      return companies.value
    }
    return companies.value.filter(c => c.group === activeGroup.value)
  })

  const fetchCompanies = async () => {
    loading.value = true
    try {
      const response = await fetch('/api/v1/companies')
      const data = await response.json()
      // 转换snake_case为camelCase
      companies.value = (data.data || []).map((c: any) => ({
        id: c.id,
        name: c.name,
        website: c.website,
        industry: c.industry,
        group: c.group,
        description: c.description,
        healthStatus: c.health_status || 'UNKNOWN'
      }))
    } catch (error) {
      console.error('Failed to fetch companies:', error)
      // 使用默认数据
      companies.value = getDefaultCompanies()
    } finally {
      loading.value = false
    }
  }

  const setGroup = (group: string) => {
    activeGroup.value = group
  }

  return {
    companies,
    loading,
    activeGroup,
    filteredCompanies,
    fetchCompanies,
    setGroup
  }
})

function getDefaultCompanies(): Company[] {
  return [
    {
      id: 'alibaba',
      name: '阿里巴巴',
      website: 'https://talent.alibaba.com',
      industry: '互联网',
      group: 'bigtech',
      description: '电商、云计算、数字媒体与金融科技综合互联网集团',
      healthStatus: 'GREEN'
    },
    {
      id: 'tencent',
      name: '腾讯',
      website: 'https://careers.tencent.com',
      industry: '互联网',
      group: 'bigtech',
      description: '社交、游戏、金融科技与企业服务科技集团',
      healthStatus: 'GREEN'
    },
    {
      id: 'bytedance',
      name: '字节跳动',
      website: 'https://jobs.bytedance.com',
      industry: '互联网',
      group: 'bigtech',
      description: '短视频、信息流内容与全球化企业协作服务',
      healthStatus: 'GREEN'
    }
  ]
}
