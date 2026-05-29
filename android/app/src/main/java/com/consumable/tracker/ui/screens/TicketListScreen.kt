package com.consumable.tracker.ui.screens

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.consumable.tracker.data.api.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun TicketListScreen() {
    var tickets by remember { mutableStateOf<List<TicketResponse>>(emptyList()) }
    var loading by remember { mutableStateOf(true) }

    LaunchedEffect(Unit) {
        try { tickets = ApiService.create().getTickets() } catch (_: Exception) {} finally { loading = false }
    }

    Column(modifier = Modifier.fillMaxSize().padding(16.dp)) {
        Text("工单管理", style = MaterialTheme.typography.titleLarge)
        Spacer(Modifier.height(12.dp))
        if (loading) { Box(Modifier.fillMaxSize()) { CircularProgressIndicator(modifier = Modifier.align(Alignment.Center)) } }
        else {
            LazyColumn {
                items(tickets) { t ->
                    Card(modifier = Modifier.fillMaxWidth().padding(vertical = 4.dp)) {
                        Column(modifier = Modifier.padding(12.dp)) {
                            Row { Text("${t.office_name} - ${t.problem_type}", style = MaterialTheme.typography.bodyLarge, modifier = Modifier.weight(1f))
                                val label = when(t.status) { "pending" -> "待处理"; "processing" -> "处理中"; else -> "已完成" }
                                val color = when(t.status) { "pending" -> MaterialTheme.colorScheme.error; "processing" -> MaterialTheme.colorScheme.tertiary; else -> MaterialTheme.colorScheme.primary }
                                Text(label, color = color)
                            }
                            Text(t.created_at, style = MaterialTheme.typography.bodySmall, color = MaterialTheme.colorScheme.onSurfaceVariant)
                        }
                    }
                }
            }
        }
    }
}
