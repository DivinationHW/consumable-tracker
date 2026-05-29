package com.consumable.tracker.ui.screens

import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.consumable.tracker.data.repository.TrackerRepository
import com.consumable.tracker.ui.theme.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun StatsScreen() {
    val repo = remember { TrackerRepository() }
    var officeCount by remember { mutableIntStateOf(0) }
    var consumableCount by remember { mutableIntStateOf(0) }
    var pendingTickets by remember { mutableIntStateOf(0) }

    LaunchedEffect(Unit) {
        try {
            val stats = repo.getStats()
            officeCount = stats.office_count
            consumableCount = stats.consumable_count
            pendingTickets = stats.pending_tickets
        } catch (_: Exception) {}
    }

    Scaffold(
        topBar = { TopAppBar(title = { Text("数据统计") }, colors = TopAppBarDefaults.topAppBarColors(containerColor = SidebarBg, titleContentColor = androidx.compose.ui.graphics.Color.White)) }
    ) { padding ->
        Column(modifier = Modifier.fillMaxSize().padding(padding).padding(16.dp), verticalArrangement = Arrangement.spacedBy(12.dp)) {
            Row(horizontalArrangement = Arrangement.spacedBy(12.dp), modifier = Modifier.fillMaxWidth()) {
                Card(modifier = Modifier.weight(1f)) { Column(Modifier.padding(16.dp)) { Text("科室总数", style = MaterialTheme.typography.bodySmall, color = TextSecondary); Text("$officeCount", style = MaterialTheme.typography.headlineMedium, color = Primary) } }
                Card(modifier = Modifier.weight(1f)) { Column(Modifier.padding(16.dp)) { Text("耗材种类", style = MaterialTheme.typography.bodySmall, color = TextSecondary); Text("$consumableCount", style = MaterialTheme.typography.headlineMedium, color = Success) } }
            }
            Card { Column(Modifier.padding(16.dp)) { Text("待处理工单", style = MaterialTheme.typography.bodySmall, color = TextSecondary); Text("$pendingTickets", style = MaterialTheme.typography.headlineMedium, color = Danger) } }
        }
    }
}
