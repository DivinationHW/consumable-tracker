package com.consumable.tracker.ui.screens

import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.navigation.NavHostController
import com.consumable.tracker.data.api.ApiService
import com.consumable.tracker.data.api.ChangePasswordRequest
import com.consumable.tracker.data.local.Preferences
import kotlinx.coroutines.launch

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun SettingsScreen(navController: NavHostController) {
    var serverUrl by remember { mutableStateOf(Preferences.baseUrl) }
    var oldPw by remember { mutableStateOf("") }
    var newPw by remember { mutableStateOf("") }
    var msg by remember { mutableStateOf("") }
    val scope = rememberCoroutineScope()

    Column(modifier = Modifier.fillMaxSize().padding(16.dp)) {
        Text("设置", style = MaterialTheme.typography.titleLarge)
        Spacer(Modifier.height(16.dp))

        Card(modifier = Modifier.fillMaxWidth()) {
            Column(modifier = Modifier.padding(16.dp)) {
                Text("服务器配置", style = MaterialTheme.typography.titleMedium)
                Spacer(Modifier.height(8.dp))
                OutlinedTextField(value = serverUrl, onValueChange = { serverUrl = it }, label = { Text("服务器地址") }, singleLine = true, modifier = Modifier.fillMaxWidth())
                Spacer(Modifier.height(8.dp))
                Button(onClick = {
                    Preferences.saveBaseUrl(serverUrl)
                    msg = "已保存，请重启应用"
                }) { Text("保存") }
            }
        }

        Spacer(Modifier.height(16.dp))

        Card(modifier = Modifier.fillMaxWidth()) {
            Column(modifier = Modifier.padding(16.dp)) {
                Text("修改密码", style = MaterialTheme.typography.titleMedium)
                Spacer(Modifier.height(8.dp))
                OutlinedTextField(value = oldPw, onValueChange = { oldPw = it }, label = { Text("原密码") }, singleLine = true, modifier = Modifier.fillMaxWidth())
                Spacer(Modifier.height(8.dp))
                OutlinedTextField(value = newPw, onValueChange = { newPw = it }, label = { Text("新密码") }, singleLine = true, modifier = Modifier.fillMaxWidth())
                Spacer(Modifier.height(8.dp))
                Button(onClick = {
                    scope.launch {
                        try {
                            ApiService.create().changePassword(ChangePasswordRequest(oldPw, newPw))
                            msg = "密码修改成功"
                        } catch (e: Exception) { msg = "失败: ${e.message}" }
                    }
                }) { Text("修改") }
            }
        }

        if (msg.isNotEmpty()) { Spacer(Modifier.height(8.dp)); Text(msg, color = MaterialTheme.colorScheme.primary) }
    }
}
