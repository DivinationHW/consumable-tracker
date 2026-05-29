package com.consumable.tracker.ui.screens

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.consumable.tracker.data.api.*
import kotlinx.coroutines.launch

@Composable
fun BackupsScreen() {
    var backups by remember { mutableStateOf<List<BackupResponse>>(emptyList()) }
    var loading by remember { mutableStateOf(true) }
    val scope = rememberCoroutineScope()

    LaunchedEffect(Unit) {
        try { backups = ApiService.create().getBackups() } catch (_: Exception) {} finally { loading = false }
    }

    LazyColumn(modifier = Modifier.fillMaxSize().padding(16.dp)) {
        item {
            Row {
                Text("备份管理", style = MaterialTheme.typography.titleLarge, modifier = Modifier.weight(1f))
                Button(onClick = {
                    scope.launch {
                        try { ApiService.create().createBackup(); backups = ApiService.create().getBackups() } catch (_: Exception) {}
                    }
                }) { Text("创建备份") }
            }
            Spacer(Modifier.height(12.dp))
        }
        items(backups) { b ->
            Card(modifier = Modifier.fillMaxWidth().padding(vertical = 4.dp)) {
                Row(modifier = Modifier.padding(12.dp)) {
                    Column(modifier = Modifier.weight(1f)) {
                        Text(b.filename, style = MaterialTheme.typography.bodyMedium)
                        Text("${b.size / 1024} KB | ${b.created_at}", style = MaterialTheme.typography.bodySmall)
                    }
                }
            }
        }
    }
}
