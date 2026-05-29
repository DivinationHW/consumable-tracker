<template>
  <div>
    <el-row :gutter="16" style="margin-bottom:16px">
      <el-col :span="6"><el-card><h3>总使用量</h3><p style="font-size:28px;color:#409eff">{{ stats.total_usage }}</p></el-card></el-col>
      <el-col :span="6"><el-card><h3>本月用量</h3><p style="font-size:28px;color:#67c23a">{{ stats.current_month }}</p></el-card></el-col>
      <el-col :span="6"><el-card><h3>最多办公室</h3><p style="font-size:20px;color:#e6a23c">{{ stats.top_office || '-' }}</p></el-card></el-col>
      <el-col :span="6"><el-card><h3>最多耗材</h3><p style="font-size:20px;color:#f56c6c">{{ stats.top_consumable || '-' }}</p></el-card></el-col>
    </el-row>
    <el-row :gutter="16">
      <el-col :span="12">
        <el-card><template #header><span>按办公室统计</span></template>
          <div ref="officeChart" style="height:350px"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card><template #header><span>月度趋势</span></template>
          <div ref="trendChart" style="height:350px"></div>
        </el-card>
      </el-col>
    </el-row>
    <el-card style="margin-top:16px">
      <template #header><span>按办公室明细</span></template>
      <el-table :data="stats.by_office" stripe>
        <el-table-column prop="label" label="办公室" />
        <el-table-column prop="value" label="使用量" sortable />
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import * as echarts from 'echarts'
import { getStats } from '@/api/stats'

const stats = ref<any>({ by_office: [], monthly_trend: [] })
const officeChart = ref<HTMLElement>()
const trendChart = ref<HTMLElement>()

async function load() {
  const res = await getStats()
  stats.value = res.data
  nextTick(() => {
    if (officeChart.value) {
      const c = echarts.init(officeChart.value)
      c.setOption({
        tooltip: {},
        xAxis: { type: 'category', data: res.data.by_office.map((i: any) => i.label), axisLabel: { rotate: 45 } },
        yAxis: { type: 'value' },
        series: [{ type: 'bar', data: res.data.by_office.map((i: any) => i.value), itemStyle: { color: '#409eff' } }],
      })
    }
    if (trendChart.value) {
      const data = [...res.data.monthly_trend].reverse()
      const c = echarts.init(trendChart.value)
      c.setOption({
        tooltip: {},
        xAxis: { type: 'category', data: data.map((i: any) => i.label) },
        yAxis: { type: 'value' },
        series: [{ type: 'line', data: data.map((i: any) => i.value), smooth: true, areaStyle: {} }],
      })
    }
  })
}

onMounted(load)
</script>
