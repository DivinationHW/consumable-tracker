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

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun QRCodeScreen(navController: NavHostController) {
    var codes by remember { mutableStateOf<List<QRCodeResponse>>(emptyList()) }
    var loading by remember { mutableStateOf(true) }

    LaunchedEffect(Unit) {
        try { codes = ApiService.create().getQRCodes() } catch (_: Exception) {} finally { loading = false }
    }

    Column(modifier = Modifier.fillMaxSize().padding(16.dp)) {
        Text("二维码管理", style = MaterialTheme.typography.titleLarge)
        Spacer(Modifier.height(12.dp))
        if (loading) { CircularProgressIndicator() }
        else {
            LazyColumn {
                items(codes) { q ->
                    Card(modifier = Modifier.fillMaxWidth().padding(vertical = 4.dp)) {
                        Row(modifier = Modifier.padding(12.dp)) {
                            Column {
                                Text("码值: ${q.code}", style = MaterialTheme.typography.bodyLarge)
                                Text("办公室: ${q.office_name}", style = MaterialTheme.typography.bodySmall)
                            }
                        }
                    }
                }
            }
        }
    }
}
