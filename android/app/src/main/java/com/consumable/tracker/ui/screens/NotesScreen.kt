package com.consumable.tracker.ui.screens

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.consumable.tracker.data.api.*

@Composable
fun NotesScreen() {
    var notes by remember { mutableStateOf<List<NoteResponse>>(emptyList()) }
    var loading by remember { mutableStateOf(true) }

    LaunchedEffect(Unit) {
        try { notes = ApiService.create().getNotes() } catch (_: Exception) {} finally { loading = false }
    }

    LazyColumn(modifier = Modifier.fillMaxSize().padding(16.dp)) {
        item { Text("便签", style = MaterialTheme.typography.titleLarge); Spacer(Modifier.height(12.dp)) }
        items(notes) { n ->
            Card(modifier = Modifier.fillMaxWidth().padding(vertical = 4.dp)) {
                Column(modifier = Modifier.padding(12.dp)) {
                    Text(n.title, style = MaterialTheme.typography.titleSmall)
                    Spacer(Modifier.height(4.dp))
                    Text(n.content, style = MaterialTheme.typography.bodyMedium)
                }
            }
        }
    }
}
