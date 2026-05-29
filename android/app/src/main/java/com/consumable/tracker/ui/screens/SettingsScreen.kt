package com.consumable.tracker.ui.screens

import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.consumable.tracker.data.api.ApiClient
import com.consumable.tracker.ui.theme.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun SettingsScreen(onLogout: () -> Unit) {
    var serverUrl by remember { mutableStateOf("http://10.0.2.2:8443") }

    Scaffold(
        topBar = { TopAppBar(title = { Text("设置") }, colors = TopAppBarDefaults.topAppBarColors(containerColor = SidebarBg, titleContentColor = androidx.compose.ui.graphics.Color.White)) }
    ) { padding ->
        Column(modifier = Modifier.fillMaxSize().padding(padding).padding(16.dp), verticalArrangement = Arrangement.spacedBy(12.dp)) {
            Card {
                Column(Modifier.padding(16.dp)) {
                    Text("服务器地址", style = MaterialTheme.typography.titleSmall)
                    Spacer(Modifier.height(8.dp))
                    OutlinedTextField(value = serverUrl, onValueChange = {
                        serverUrl = it
                        ApiClient.setBaseUrl(it)
                    }, modifier = Modifier.fillMaxWidth(), singleLine = true)
                }
            }
            Spacer(Modifier.weight(1f))
            Button(onClick = {
                ApiClient.setToken(null)
                onLogout()
            }, modifier = Modifier.fillMaxWidth(), colors = ButtonDefaults.buttonColors(containerColor = Danger)) {
                Text("退出登录")
            }
        }
    }
}
