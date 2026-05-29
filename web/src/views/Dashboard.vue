<template>
  <div class="page-container">
    <div class="page-header"><h2>仪表盘</h2></div>
    <div class="stat-cards">
      <div class="stat-card"><div class="label">科室总数</div><div class="value blue">{{ stats.office_count || 0 }}</div></div>
      <div class="stat-card"><div class="label">耗材种类</div><div class="value green">{{ stats.consumable_count || 0 }}</div></div>
      <div class="stat-card"><div class="label">本月消耗量</div><div class="value orange">{{ stats.monthly_usage || 0 }}</div></div>
      <div class="stat-card"><div class="label">待处理工单</div><div class="value red">{{ stats.pending_tickets || 0 }}</div></div>
    </div>
    <div class="grid-2">
      <div class="chart-container">
        <h3>本月科室耗材排行</h3>
        <v-chart class="chart-box" :option="officeChart" autoresize />
      </div>
      <div class="chart-container">
        <h3>耗材使用趋势 (近12月)</h3>
        <v-chart class="chart-box" :option="trendChart" autoresize />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { getStats } from '@/api'
import VChart from 'vue-echarts'
import 'echarts'

const stats = ref({})
const officeData = ref([])
const trendData = ref([])

onMounted(async () => {
  const data = await getStats()
  stats.value = data
  officeData.value = data.office_ranking || []
  trendData.value = data.monthly_trend || []
})

const officeChart = computed(() => ({
  tooltip: { trigger: 'axis' },
  xAxis: { type: 'category', data: officeData.value.map(o => o.room_number) },
  yAxis: { type: 'value' },
  series: [{ type: 'bar', data: officeData.value.map(o => o.total), itemStyle: { color: '#409eff' } }],
}))

const trendChart = computed(() => ({
  tooltip: { trigger: 'axis' },
  xAxis: { type: 'category', data: trendData.value.map(t => t.month) },
  yAxis: { type: 'value' },
  series: [{ type: 'line', data: trendData.value.map(t => t.total), smooth: true, lineStyle: { color: '#67c23a' }, areaStyle: { color: 'rgba(103,194,58,0.1)' } }],
}))
</script>
