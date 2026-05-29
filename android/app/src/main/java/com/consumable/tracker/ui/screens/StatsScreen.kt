package com.consumable.tracker.ui.screens

import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.consumable.tracker.data.api.*
import com.patrykandpatrick.vico.compose.chart.Chart
import com.patrykandpatrick.vico.compose.chart.line.lineChart
import com.patrykandpatrick.vico.core.entry.entryModelOf

@Composable
fun StatsScreen() {
    var stats by remember { mutableStateOf<StatsResponse?>(null) }

    LaunchedEffect(Unit) {
        try { stats = ApiService.create().getStats() } catch (_: Exception) {}
    }

    Column(modifier = Modifier.fillMaxSize().padding(16.dp)) {
        Text("统计报表", style = MaterialTheme.typography.titleLarge)
        Spacer(Modifier.height(16.dp))

        stats?.let { s ->
            Row(modifier = Modifier.fillMaxWidth()) {
                Card(modifier = Modifier.weight(1f).padding(4.dp)) { Column(Modifier.padding(12.dp)) { Text("总用量"); Text("${s.total_usage}", style = MaterialTheme.typography.headlineSmall) } }
                Card(modifier = Modifier.weight(1f).padding(4.dp)) { Column(Modifier.padding(12.dp)) { Text("本月"); Text("${s.current_month}", style = MaterialTheme.typography.headlineSmall) } }
            }
            Spacer(Modifier.height(16.dp))
            Text("月度趋势", style = MaterialTheme.typography.titleMedium)
            Spacer(Modifier.height(8.dp))
            if (s.monthly_trend.isNotEmpty()) {
                val entries = s.monthly_trend.reversed().mapIndexed { i, item -> i to item.value.toFloat() }
                Chart(chart = lineChart(), model = entryModelOf(*entries.toTypedArray()))
            }
        }
    }
}
