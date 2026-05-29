package com.consumable.tracker.ui.screens

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.consumable.tracker.data.models.UsageRecord
import com.consumable.tracker.data.repository.TrackerRepository
import com.consumable.tracker.ui.theme.*
import java.text.SimpleDateFormat
import java.util.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun RecordsScreen() {
    val repo = remember { TrackerRepository() }
    var records by remember { mutableStateOf<List<UsageRecord>>(emptyList()) }
    var loading by remember { mutableStateOf(true) }

    LaunchedEffect(Unit) {
        try {
            val cal = Calendar.getInstance()
            val fmt = SimpleDateFormat("yyyy-MM-dd", Locale.CHINA)
            records = repo.getRecords(mapOf(
                "start_date" to fmt.format(Date(cal.timeInMillis - 30L * 24 * 3600 * 1000)),
                "end_date" to fmt.format(Date())
            ))
        } catch (_: Exception) {}
        loading = false
    }

    Scaffold(
        topBar = { TopAppBar(title = { Text("使用记录") }, colors = TopAppBarDefaults.topAppBarColors(containerColor = SidebarBg, titleContentColor = androidx.compose.ui.graphics.Color.White)) }
    ) { padding ->
        if (loading) {
            Box(Modifier.fillMaxSize().padding(padding), contentAlignment = androidx.compose.ui.Alignment.Center) {
                CircularProgressIndicator()
            }
        } else {
            LazyColumn(modifier = Modifier.fillMaxSize().padding(padding).padding(16.dp), verticalArrangement = Arrangement.spacedBy(8.dp)) {
                items(records) { record ->
                    Card(modifier = Modifier.fillMaxWidth()) {
                        Column(modifier = Modifier.padding(12.dp)) {
                            Text("${record.consumable_name} x${record.quantity}${record.unit}", style = MaterialTheme.typography.titleSmall)
                            Spacer(Modifier.height(4.dp))
                            Text("${record.room_number} | ${record.usage_date}", style = MaterialTheme.typography.bodySmall, color = TextSecondary)
                            if (record.note.isNotBlank()) {
                                Text(record.note, style = MaterialTheme.typography.bodySmall, color = TextSecondary)
                            }
                        }
                    }
                }
                if (records.isEmpty()) {
                    item { Text("暂无记录", modifier = Modifier.padding(16.dp), color = TextSecondary) }
                }
            }
        }
    }
}
