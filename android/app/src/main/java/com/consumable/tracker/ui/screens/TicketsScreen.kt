package com.consumable.tracker.ui.screens

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.consumable.tracker.data.models.Ticket
import com.consumable.tracker.data.repository.TrackerRepository
import com.consumable.tracker.ui.theme.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun TicketsScreen() {
    val repo = remember { TrackerRepository() }
    var tickets by remember { mutableStateOf<List<Ticket>>(emptyList()) }
    var loading by remember { mutableStateOf(true) }

    LaunchedEffect(Unit) {
        try { tickets = repo.getTickets() } catch (_: Exception) {}
        loading = false
    }

    Scaffold(
        topBar = { TopAppBar(title = { Text("报修工单") }, colors = TopAppBarDefaults.topAppBarColors(containerColor = SidebarBg, titleContentColor = androidx.compose.ui.graphics.Color.White)) }
    ) { padding ->
        if (loading) {
            Box(Modifier.fillMaxSize().padding(padding), contentAlignment = androidx.compose.ui.Alignment.Center) {
                CircularProgressIndicator()
            }
        } else {
            LazyColumn(modifier = Modifier.fillMaxSize().padding(padding).padding(16.dp), verticalArrangement = Arrangement.spacedBy(8.dp)) {
                items(tickets) { ticket ->
                    Card(modifier = Modifier.fillMaxWidth()) {
                        Column(modifier = Modifier.padding(12.dp)) {
                            Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.SpaceBetween) {
                                Text(ticket.room_number, style = MaterialTheme.typography.titleSmall)
                                val statusColors = mapOf("pending" to Danger, "processing" to Warning, "completed" to Success)
                                Text(
                                    mapOf("pending" to "待处理", "processing" to "处理中", "completed" to "已完成")[ticket.status] ?: ticket.status,
                                    color = statusColors[ticket.status] ?: TextSecondary,
                                    style = MaterialTheme.typography.bodySmall
                                )
                            }
                            Spacer(Modifier.height(4.dp))
                            Text(ticket.problem_type, style = MaterialTheme.typography.bodyMedium)
                            if (ticket.description.isNotBlank()) {
                                Text(ticket.description, style = MaterialTheme.typography.bodySmall, color = TextSecondary)
                            }
                        }
                    }
                }
            }
        }
    }
}
