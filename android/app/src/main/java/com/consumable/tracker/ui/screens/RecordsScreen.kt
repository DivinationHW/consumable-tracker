package com.consumable.tracker.ui.screens

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.navigation.NavHostController
import com.consumable.tracker.data.api.*
import kotlinx.coroutines.launch

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun RecordsScreen(navController: NavHostController) {
    var records by remember { mutableStateOf<List<RecordResponse>>(emptyList()) }
    var loading by remember { mutableStateOf(true) }
    val scope = rememberCoroutineScope()

    LaunchedEffect(Unit) {
        try {
            val api = ApiService.create()
            records = api.getRecords(mapOf("page" to "1", "page_size" to "200")).data
        } catch (_: Exception) {} finally { loading = false }
    }

    Column(modifier = Modifier.fillMaxSize().padding(16.dp)) {
        Text("记录查询", style = MaterialTheme.typography.titleLarge)
        Spacer(Modifier.height(12.dp))
        if (loading) { Box(Modifier.fillMaxSize()) { CircularProgressIndicator(modifier = Modifier.align(androidx.compose.ui.Alignment.Center)) } }
        else {
            LazyColumn {
                items(records) { r ->
                    Card(modifier = Modifier.fillMaxWidth().padding(vertical = 4.dp)) {
                        Column(modifier = Modifier.padding(12.dp)) {
                            Row { Text("${r.usage_date} | ${r.office_name}", style = MaterialTheme.typography.bodyLarge) }
                            Row { Text("${r.consumable_name} x${r.quantity}", style = MaterialTheme.typography.bodyMedium) }
                            if (r.note.isNotEmpty()) Text("备注: ${r.note}", style = MaterialTheme.typography.bodySmall, color = MaterialTheme.colorScheme.onSurfaceVariant)
                        }
                    }
                }
            }
        }
    }
}
