package com.consumable.tracker.ui.screens

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.unit.dp
import com.consumable.tracker.data.repository.TrackerRepository
import com.consumable.tracker.ui.navigation.Screen
import com.consumable.tracker.ui.theme.*
import kotlinx.coroutines.launch

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun DashboardScreen(onNavigate: (Screen) -> Unit) {
    val repo = remember { TrackerRepository() }
    val scope = rememberCoroutineScope()
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
        topBar = { TopAppBar(title = { Text("耗材管理系统") }, colors = TopAppBarDefaults.topAppBarColors(containerColor = SidebarBg, titleContentColor = androidx.compose.ui.graphics.Color.White)) }
    ) { padding ->
        LazyColumn(modifier = Modifier.fillMaxSize().padding(padding).padding(16.dp), verticalArrangement = Arrangement.spacedBy(12.dp)) {
            item {
                Row(horizontalArrangement = Arrangement.spacedBy(12.dp), modifier = Modifier.fillMaxWidth()) {
                    StatCard("科室数", "$officeCount", Primary, Modifier.weight(1f))
                    StatCard("耗材种类", "$consumableCount", Success, Modifier.weight(1f))
                }
            }
            item {
                StatCard("待处理工单", "$pendingTickets", Danger, Modifier.fillMaxWidth())
            }
            item { MenuButton(Icons.Default.Description, "使用记录") { onNavigate(Screen.Records) } }
            item { MenuButton(Icons.Default.BarChart, "数据统计") { onNavigate(Screen.Stats) } }
            item { MenuButton(Icons.Default.Build, "报修工单") { onNavigate(Screen.Tickets) } }
            item { MenuButton(Icons.Default.QrCodeScanner, "扫码报修") { onNavigate(Screen.QrScanner) } }
            item { MenuButton(Icons.Default.Settings, "设置") { onNavigate(Screen.Settings) } }
        }
    }
}

@Composable
fun StatCard(label: String, value: String, color: androidx.compose.ui.graphics.Color, modifier: Modifier = Modifier) {
    Card(modifier = modifier.padding(4.dp)) {
        Column(modifier = Modifier.padding(16.dp)) {
            Text(label, style = MaterialTheme.typography.bodySmall, color = TextSecondary)
            Spacer(Modifier.height(4.dp))
            Text(value, style = MaterialTheme.typography.headlineMedium, color = color)
        }
    }
}

@Composable
fun MenuButton(icon: ImageVector, label: String, onClick: () -> Unit) {
    Card(modifier = Modifier.fillMaxWidth().clickable(onClick = onClick)) {
        Row(modifier = Modifier.padding(16.dp), verticalAlignment = Alignment.CenterVertically) {
            Icon(icon, contentDescription = null, tint = Primary)
            Spacer(Modifier.width(12.dp))
            Text(label, style = MaterialTheme.typography.bodyLarge)
        }
    }
}
