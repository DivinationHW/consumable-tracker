<template>
  <div class="page-container">
    <div class="page-header"><h2>数据统计</h2></div>
    <div class="filters">
      <el-date-picker v-model="year" type="year" placeholder="选择年份" @change="loadData" />
      <el-button @click="exportExcel">导出Excel</el-button>
    </div>
    <div class="grid-2">
      <div class="chart-container">
        <h3>耗材消耗排行</h3>
        <v-chart class="chart-box" :option="consumableChart" autoresize />
      </div>
      <div class="chart-container">
        <h3>科室消耗排行</h3>
        <v-chart class="chart-box" :option="officeChart" autoresize />
      </div>
    </div>
    <div class="chart-container">
      <h3>月度消耗趋势</h3>
      <v-chart class="chart-box" :option="trendChart" autoresize />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { getStats } from '@/api'
import VChart from 'vue-echarts'
import 'echarts'
import * as XLSX from 'xlsx'

const year = ref(new Date())
const consumableRank = ref([])
const officeRank = ref([])
const trendData = ref([])

onMounted(loadData)

async function loadData() {
  const params = { year: year.value.getFullYear() }
  const data = await getStats(params)
  consumableRank.value = data.consumable_ranking || []
  officeRank.value = data.office_ranking || []
  trendData.value = data.monthly_trend || []
}

const consumableChart = computed(() => ({
  tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
  xAxis: { type: 'value' },
  yAxis: { type: 'category', data: consumableRank.value.map(c => c.name).reverse(), axisLabel: { fontSize: 11 } },
  series: [{ type: 'bar', data: consumableRank.value.map(c => c.total).reverse(), itemStyle: { color: '#409eff' } }],
}))

const officeChart = computed(() => ({
  tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
  xAxis: { type: 'value' },
  yAxis: { type: 'category', data: officeRank.value.map(o => o.room_number).reverse() },
  series: [{ type: 'bar', data: officeRank.value.map(o => o.total).reverse(), itemStyle: { color: '#67c23a' } }],
}))

const trendChart = computed(() => ({
  tooltip: { trigger: 'axis' },
  xAxis: { type: 'category', data: trendData.value.map(t => t.month) },
  yAxis: { type: 'value' },
  series: [{ type: 'line', data: trendData.value.map(t => t.total), smooth: true, lineStyle: { color: '#e6a23c' }, areaStyle: { color: 'rgba(230,162,60,0.1)' } }],
}))

function exportExcel() {
  const wb = XLSX.utils.book_new()
  if (consumableRank.value.length) {
    const ws = XLSX.utils.json_to_sheet(consumableRank.value)
    XLSX.utils.book_append_sheet(wb, ws, '耗材排行')
  }
  if (officeRank.value.length) {
    const ws = XLSX.utils.json_to_sheet(officeRank.value)
    XLSX.utils.book_append_sheet(wb, ws, '科室排行')
  }
  if (trendData.value.length) {
    const ws = XLSX.utils.json_to_sheet(trendData.value)
    XLSX.utils.book_append_sheet(wb, ws, '月度趋势')
  }
  XLSX.writeFile(wb, `统计数据_${year.value.getFullYear()}.xlsx`)
}
</script>
