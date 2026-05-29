package com.consumable.tracker.ui.screens

import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.consumable.tracker.data.api.ApiClient
import com.consumable.tracker.data.repository.TrackerRepository
import com.consumable.tracker.ui.theme.*
import kotlinx.coroutines.launch

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun TicketFormScreen(code: String, onSubmitted: () -> Unit) {
    val repo = remember { TrackerRepository() }
    val scope = rememberCoroutineScope()
    var problemType by remember { mutableStateOf("") }
    var description by remember { mutableStateOf("") }
    var contact by remember { mutableStateOf("") }
    var loading by remember { mutableStateOf(false) }
    var submitted by remember { mutableStateOf(false) }
    var ticketId by remember { mutableStateOf("") }
    var error by remember { mutableStateOf<String?>(null) }

    Scaffold(
        topBar = { TopAppBar(title = { Text("提交报修") }, colors = TopAppBarDefaults.topAppBarColors(containerColor = SidebarBg, titleContentColor = androidx.compose.ui.graphics.Color.White)) }
    ) { padding ->
        Column(modifier = Modifier.fillMaxSize().padding(padding).padding(16.dp), verticalArrangement = Arrangement.spacedBy(12.dp)) {
            if (!submitted) {
                Text("二维码编码: $code", style = MaterialTheme.typography.bodySmall, color = TextSecondary)
                OutlinedTextField(value = problemType, onValueChange = { problemType = it }, label = { Text("问题类型") }, modifier = Modifier.fillMaxWidth())
                OutlinedTextField(value = description, onValueChange = { description = it }, label = { Text("问题描述（可选）") }, modifier = Modifier.fillMaxWidth(), maxLines = 3)
                OutlinedTextField(value = contact, onValueChange = { contact = it }, label = { Text("联系方式（可选）") }, modifier = Modifier.fillMaxWidth())

                error?.let { Text(it, color = MaterialTheme.colorScheme.error) }

                Button(onClick = {
                    scope.launch {
                        loading = true
                        error = null
                        try {
                            val res = repo.submitTicket(com.consumable.tracker.data.models.TicketCreate(
                                office_id = 0, problem_type = problemType, description = description, contact = contact
                            ))
                            ticketId = res.id
                            submitted = true
                        } catch (e: Exception) {
                            error = e.message ?: "提交失败"
                        } finally { loading = false }
                    }
                }, enabled = !loading, modifier = Modifier.fillMaxWidth()) {
                    if (loading) CircularProgressIndicator(modifier = Modifier.size(20.dp), strokeWidth = 2.dp)
                    else Text("提交报修")
                }
            } else {
                Text("报修提交成功", style = MaterialTheme.typography.headlineSmall)
                Text("工单编号: $ticketId")
                Button(onClick = onSubmitted) { Text("返回") }
            }
        }
    }
}
